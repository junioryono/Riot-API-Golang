package ratelimiter

import (
	"context"
	"time"

	"github.com/junioryono/Riot-API-Golang/apiclient/limiter"
)

const initialLimit = 20

func newRateLimit() *RateLimit {
	return &RateLimit{
		shortLimiter:      limiter.NewLimiter(initialLimit),
		longLimiter:       limiter.NewLimiter(initialLimit),
		blockedUntil:      time.Time{},
		blockedUntilQueue: make(chan struct{}, 1),
	}
}

func (rl *RateLimiter) getRegionLimiter(region string) *RateLimit {
	rl.regionMutex.Lock()
	defer rl.regionMutex.Unlock()

	if limiter, ok := rl.regionLimiters[region]; ok {
		return limiter
	}

	limiter := newRateLimit()
	rl.regionLimiters[region] = limiter
	return limiter
}

func (rl *RateLimiter) getMethodLimiter(method string) *RateLimit {
	rl.methodMutex.Lock()
	defer rl.methodMutex.Unlock()

	if limiter, ok := rl.methodLimiters[method]; ok {
		return limiter
	}

	limiter := newRateLimit()
	rl.methodLimiters[method] = limiter
	return limiter
}

func (rl *RateLimiter) waitForLimiters(ctx context.Context, priority int, regionLimiter, methodLimiter *RateLimit, isRetryRequest bool) error {
	var lCtx context.Context
	if ctx == nil {
		lCtx = context.Background()
	} else {
		lCtx = ctx
	}

	if err := rl.waitForBlock(lCtx, regionLimiter); err != nil {
		return err
	}

	if err := rl.waitForBlock(lCtx, methodLimiter); err != nil {
		return err
	}

	// If this is a not a retry request, obtain tokens from all limiters
	// We do this because retry requests do not release tokens
	if !isRetryRequest {
		if err := regionLimiter.shortLimiter.Obtain(lCtx, priority); err != nil {
			return err
		}

		if err := regionLimiter.longLimiter.Obtain(lCtx, priority); err != nil {
			regionLimiter.shortLimiter.Release()
			return err
		}

		if err := methodLimiter.shortLimiter.Obtain(lCtx, priority); err != nil {
			regionLimiter.shortLimiter.Release()
			regionLimiter.longLimiter.Release()
			return err
		}
	}

	return nil
}

func (rl *RateLimiter) waitForBlock(ctx context.Context, limiter *RateLimit) error {
	// Wait for our turn in the queue
	select {
	case limiter.blockedUntilQueue <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}

	defer func() { <-limiter.blockedUntilQueue }() // Leave the queue

	now := time.Now()
	if now.Before(limiter.blockedUntil) {
		blockedUntil := limiter.blockedUntil
		delay := blockedUntil.Sub(now)
		timer := time.NewTimer(delay)
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
	regionLimiter.shortLimiter.Release()
	regionLimiter.longLimiter.Release()
	methodLimiter.shortLimiter.Release()
}

func (rl *RateLimiter) releaseLimitersAfterDelay(regionLimiter, methodLimiter *RateLimit, delay time.Duration) {
	<-time.After(delay)
	rl.releaseLimiters(regionLimiter, methodLimiter)
}
