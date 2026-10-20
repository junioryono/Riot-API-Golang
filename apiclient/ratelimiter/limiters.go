package ratelimiter

import (
	"context"
	"math"
	"time"

	"github.com/junioryono/Riot-API-Golang/apiclient/limiter"
)

const initialLimit = 20

func newRateLimit() *RateLimit {
	return &RateLimit{
		buckets:           []*limiter.Limiter{limiter.NewLimiter(initialLimit)},
		windows:           []time.Duration{0},
		active:            1,
		blockedUntil:      time.Time{},
		blockedUntilQueue: make(chan struct{}, 1),
	}
}

func (rl *RateLimiter) getRegionLimiter(region string) *RateLimit {
	rl.regionMutex.Lock()
	defer rl.regionMutex.Unlock()

	if limit, ok := rl.regionLimiters[region]; ok {
		return limit
	}

	limit := newRateLimit()
	rl.regionLimiters[region] = limit
	return limit
}

func (rl *RateLimiter) getMethodLimiter(method string) *RateLimit {
	rl.methodMutex.Lock()
	defer rl.methodMutex.Unlock()

	if limit, ok := rl.methodLimiters[method]; ok {
		return limit
	}

	limit := newRateLimit()
	rl.methodLimiters[method] = limit
	return limit
}

func (rl *RateLimiter) waitForLimiterBlocks(ctx context.Context, regionLimiter, methodLimiter *RateLimit) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := rl.waitForBlock(ctx, regionLimiter); err != nil {
		return err
	}
	return rl.waitForBlock(ctx, methodLimiter)
}

func (rl *RateLimiter) obtainBuckets(ctx context.Context, priority int, regionLimiter, methodLimiter *RateLimit) ([]*limiter.Limiter, []*limiter.Limiter, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	heldRegion := regionLimiter.snapshot()
	heldMethod := methodLimiter.snapshot()
	if err := obtainAll(ctx, priority, heldRegion); err != nil {
		return nil, nil, err
	}
	if err := obtainAll(ctx, priority, heldMethod); err != nil {
		releaseNow(heldRegion)
		return nil, nil, err
	}
	return heldRegion, heldMethod, nil
}

// waitForLimiters blocks on app and method limits, then takes one token from
// each active bucket. Retry requests skip the token take; their original tokens
// are still held.
func (rl *RateLimiter) waitForLimiters(ctx context.Context, priority int, regionLimiter, methodLimiter *RateLimit, isRetryRequest bool) error {
	if err := rl.waitForLimiterBlocks(ctx, regionLimiter, methodLimiter); err != nil {
		return err
	}
	if isRetryRequest {
		return nil
	}
	_, _, err := rl.obtainBuckets(ctx, priority, regionLimiter, methodLimiter)
	return err
}

func obtainAll(ctx context.Context, priority int, buckets []*limiter.Limiter) error {
	for i, bucket := range buckets {
		if err := bucket.Obtain(ctx, priority); err != nil {
			releaseNow(buckets[:i])
			return err
		}
	}
	return nil
}

func releaseNow(buckets []*limiter.Limiter) {
	for _, bucket := range buckets {
		bucket.Release()
	}
}

func releaseAllAfterDelay(delay time.Duration, groups ...[]*limiter.Limiter) {
	if delay > 0 {
		timer := time.NewTimer(delay)
		<-timer.C
	}
	for _, group := range groups {
		releaseNow(group)
	}
}

func (rl *RateLimiter) waitForBlock(ctx context.Context, limit *RateLimit) error {
	select {
	case limit.blockedUntilQueue <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	defer func() { <-limit.blockedUntilQueue }()

	limit.mu.Lock()
	blockedUntil := limit.blockedUntil
	limit.mu.Unlock()

	now := time.Now()
	if now.Before(blockedUntil) {
		timer := time.NewTimer(blockedUntil.Sub(now))
		defer timer.Stop()

		select {
		case <-timer.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (rl *RateLimiter) releaseLimiters(regionLimiter, methodLimiter *RateLimit) {
	releaseNow(regionLimiter.snapshot())
	releaseNow(methodLimiter.snapshot())
}

func (rl *RateLimiter) releaseLimitersAfterDelay(regionLimiter, methodLimiter *RateLimit, delay time.Duration) {
	releaseAllAfterDelay(delay, regionLimiter.snapshot(), methodLimiter.snapshot())
}

func (r *RateLimit) snapshot() []*limiter.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]*limiter.Limiter, r.active)
	copy(out, r.buckets[:r.active])
	return out
}

func (r *RateLimit) applyHeader(buckets []limitWindow, percent int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for len(r.buckets) < len(buckets) {
		r.buckets = append(r.buckets, limiter.NewLimiter(initialLimit))
		r.windows = append(r.windows, 0)
	}

	var blockWindow time.Duration
	now := time.Now()
	for i, bucket := range buckets {
		capacity := conservedLimit(bucket.limit, percent)
		r.buckets[i].SetCapacity(capacity)
		r.windows[i] = bucket.window
		if bucket.count > capacity && bucket.window > blockWindow {
			blockWindow = bucket.window
		}
	}
	for i := len(buckets); i < len(r.buckets); i++ {
		r.buckets[i].SetCapacity(math.MaxInt32)
		r.windows[i] = 0
	}
	r.active = len(buckets)

	if blockWindow > 0 && now.After(r.blockedUntil) {
		r.blockedUntil = now.Add(blockWindow)
	}
}

func (r *RateLimit) blockFor(delay time.Duration) {
	if delay <= 0 {
		return
	}
	r.mu.Lock()
	r.blockedUntil = time.Now().Add(delay)
	r.mu.Unlock()
}

func (r *RateLimit) releaseHeld(held []*limiter.Limiter, fallback time.Duration) {
	for i, bucket := range held {
		delay := r.delayFor(i, fallback)
		if delay <= 0 {
			bucket.Release()
			continue
		}
		bucket.ReleaseAfterDelay(delay)
	}
}

func (r *RateLimit) delayFor(i int, fallback time.Duration) time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	if i >= r.active {
		return 0
	}
	if i < len(r.windows) && r.windows[i] > 0 {
		return r.windows[i]
	}
	return fallback
}

func (r *RateLimit) capacities() []int {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]int, r.active)
	for i := 0; i < r.active; i++ {
		out[i] = r.buckets[i].Capacity()
	}
	return out
}

func (r *RateLimit) releaseWindows() []time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]time.Duration, r.active)
	copy(out, r.windows[:r.active])
	return out
}

func (r *RateLimit) setCapacities(caps ...int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.buckets = make([]*limiter.Limiter, len(caps))
	r.windows = make([]time.Duration, len(caps))
	for i, cap := range caps {
		r.buckets[i] = limiter.NewLimiter(cap)
		r.windows[i] = time.Second
	}
	r.active = len(caps)
}
