package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/junioryono/Riot-API-Golang/apiclient/limiter"
)

type APIRequest struct {
	Context  context.Context
	Priority int
	Region   string
	MethodID MethodID
	URL      string
	Response chan<- *http.Response
	Error    chan<- error
	Retries  int

	// heldRegion and heldMethod are the buckets this request took tokens from.
	// Retries reuse them instead of taking another token per attempt.
	heldRegion []*limiter.Limiter
	heldMethod []*limiter.Limiter
}

// RateLimit is one app or method limiter. Each header bucket has its own
// token bucket and release window. active is how many of those buckets
// currently match the last valid header.
type RateLimit struct {
	mu                sync.Mutex
	buckets           []*limiter.Limiter
	windows           []time.Duration
	active            int
	blockedUntil      time.Time
	blockedUntilQueue chan struct{}
}

func (rl *RateLimiter) handleRequest(req *APIRequest) {
	regionLimiter := rl.getRegionLimiter(req.Region)
	methodLimiter := rl.getMethodLimiter(req.Region + req.MethodID.String())

	if err := rl.waitForLimiterBlocks(req.Context, regionLimiter, methodLimiter); err != nil {
		req.Error <- err
		return
	}

	// A retry keeps the tokens taken on the first attempt.
	if req.Retries == 0 {
		heldRegion, heldMethod, err := rl.obtainBuckets(req.Context, req.Priority, regionLimiter, methodLimiter)
		if err != nil {
			req.Error <- err
			return
		}
		req.heldRegion = heldRegion
		req.heldMethod = heldMethod
	}

	httpRequest, err := rl.createHTTPRequest(req)
	if err != nil {
		req.Error <- err
		releaseNow(req.heldRegion)
		releaseNow(req.heldMethod)
		return
	}

	resp, err := rl.httpClient.Do(httpRequest)
	if err != nil {
		req.Error <- err
		releaseNow(req.heldRegion)
		releaseNow(req.heldMethod)
		return
	}

	rl.handleHTTPResponse(req, resp, regionLimiter, methodLimiter)
}

func (rl *RateLimiter) createHTTPRequest(req *APIRequest) (*http.Request, error) {
	var (
		httpRequest *http.Request
		err         error
	)

	if req.Context == nil {
		httpRequest, err = http.NewRequest("GET", req.URL, nil)
	} else {
		httpRequest, err = http.NewRequestWithContext(req.Context, "GET", req.URL, nil)
	}

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("X-Riot-Token", rl.apiKey)
	return httpRequest, nil
}

func (rl *RateLimiter) handleHTTPResponse(req *APIRequest, resp *http.Response, regionLimiter, methodLimiter *RateLimit) {
	if resp.StatusCode == http.StatusOK {
		req.Response <- resp
		rl.updateRateLimits(resp, req.MethodID, regionLimiter, methodLimiter, req.heldRegion, req.heldMethod)
		return
	}

	if resp.StatusCode == http.StatusForbidden {
		req.Response <- resp
		go releaseAllAfterDelay(15*time.Second, req.heldRegion, req.heldMethod)
		return
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		// Retry the request if Retries is less than maxRetries, or if maxRetries is -1. Otherwise, send the response to the channel
		if req.Retries < rl.maxRetries || rl.maxRetries == -1 {
			resp.Body.Close()
			rl.handleRateLimitedResponse(req, resp, regionLimiter, methodLimiter, true)

			req.Retries++
			go rl.handleRequest(req)
		} else {
			req.Response <- resp
			rl.handleRateLimitedResponse(req, resp, regionLimiter, methodLimiter, false)
		}

		return
	}

	if !isBadResponse(resp) && (req.Retries < rl.maxRetries || rl.maxRetries == -1) {
		resp.Body.Close()
		// Back off before the retry. Tokens stay held so the retry cannot
		// exceed the bucket, and the later success or give-up releases them once.
		timer := time.NewTimer(15 * time.Second)
		<-timer.C

		req.Retries++
		go rl.handleRequest(req)
		return
	}

	// Close the response body and send an error to the channel. Do not know how to handle this response
	resp.Body.Close()
	req.Error <- fmt.Errorf("received unexpected status code %d (%s)", resp.StatusCode, req.URL)
	go releaseAllAfterDelay(15*time.Second, req.heldRegion, req.heldMethod)
}

func (rl *RateLimiter) handleRateLimitedResponse(req *APIRequest, resp *http.Response, regionLimiter *RateLimit, methodLimiter *RateLimit, retryRequest bool) {
	retryAfterHeader := resp.Header.Get("Retry-After")
	rateLimitTypeHeader := resp.Header.Get("X-Rate-Limit-Type")
	retryAfter, err := strconv.Atoi(retryAfterHeader)
	if err != nil {
		retryAfter = 15
	}

	retryAfterDuration := time.Duration(retryAfter) * time.Second

	if rateLimitTypeHeader == "application" {
		regionLimiter.blockFor(retryAfterDuration)
	} else if rateLimitTypeHeader == "method" {
		methodLimiter.blockFor(retryAfterDuration)
	}

	// If we are not going to retry the request, release the limiters.
	// A retry keeps its tokens so it stays in line without taking another.
	if !retryRequest {
		go releaseAllAfterDelay(retryAfterDuration, req.heldRegion, req.heldMethod)
	}
}

func isBadResponse(resp *http.Response) bool {
	return resp == nil || resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusUnsupportedMediaType
}

var (
	logMissingApp      sync.Once
	logMalformedApp    sync.Once
	logMissingMethod   sync.Once
	logMalformedMethod sync.Once
)

func (rl *RateLimiter) updateRateLimits(resp *http.Response, methodID MethodID, regionLimiter *RateLimit, methodLimiter *RateLimit, heldRegion, heldMethod []*limiter.Limiter) {
	rl.syncLimitHeader(
		"app",
		resp.Header.Get("X-App-Rate-Limit"),
		resp.Header.Get("X-App-Rate-Limit-Count"),
		methodID,
		true,
		regionLimiter,
		heldRegion,
		&logMissingApp,
		&logMalformedApp,
	)
	rl.syncLimitHeader(
		"method",
		resp.Header.Get("X-Method-Rate-Limit"),
		resp.Header.Get("X-Method-Rate-Limit-Count"),
		methodID,
		false,
		methodLimiter,
		heldMethod,
		&logMissingMethod,
		&logMalformedMethod,
	)
}

func (rl *RateLimiter) syncLimitHeader(kind, limitHeader, countHeader string, methodID MethodID, isRegion bool, limit *RateLimit, held []*limiter.Limiter, missingOnce, malformedOnce *sync.Once) {
	// Missing headers release after 15s, malformed ones after 60s, until a valid
	// header has set that bucket's window. A window learned earlier is reused.
	if limitHeader == "" || countHeader == "" {
		reportLimitHeader(missingOnce, kind, "missing", limitHeader+" "+countHeader)
		limit.releaseHeld(held, missingHeaderDelay)
		return
	}

	buckets, err := parseLimitBuckets(limitHeader, countHeader)
	if err != nil {
		reportLimitHeader(malformedOnce, kind, "malformed", limitHeader)
		limit.releaseHeld(held, malformedHeaderDelay)
		return
	}

	limit.applyHeader(buckets, rl.conservationPercent(methodID, isRegion))
	limit.releaseHeld(held, malformedHeaderDelay)
}

func (rl *RateLimiter) conservationPercent(methodID MethodID, isRegion bool) int {
	if isRegion {
		return rl.conserveUsage.RegionPercent
	}

	percent := rl.conserveUsage.MethodPercent
	if percent <= 0 {
		return 0
	}
	for _, ignored := range rl.conserveUsage.IgnoreLimits {
		if ignored == methodID {
			return 0
		}
	}
	return percent
}
