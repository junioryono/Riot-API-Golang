package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
}

type RateLimit struct {
	shortLimiter      *limiter.Limiter
	longLimiter       *limiter.Limiter
	blockedUntil      time.Time
	blockedUntilQueue chan struct{}
}

func (rl *RateLimiter) handleRequest(req *APIRequest) {
	regionLimiter := rl.getRegionLimiter(req.Region)
	methodLimiter := rl.getMethodLimiter(req.Region + req.MethodID.String())

	isRetryRequest := req.Retries > 0
	if err := rl.waitForLimiters(req.Context, req.Priority, regionLimiter, methodLimiter, isRetryRequest); err != nil {
		req.Error <- err
		return
	}

	httpRequest, err := rl.createHTTPRequest(req)
	if err != nil {
		req.Error <- err
		rl.releaseLimiters(regionLimiter, methodLimiter)
		return
	}

	resp, err := rl.httpClient.Do(httpRequest)
	if err != nil {
		req.Error <- err
		rl.releaseLimiters(regionLimiter, methodLimiter)
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
		rl.updateRateLimits(resp, req.MethodID, regionLimiter, methodLimiter)
		return
	}

	if resp.StatusCode == http.StatusForbidden {
		req.Response <- resp
		rl.releaseLimitersAfterDelay(regionLimiter, methodLimiter, 15*time.Second)
		return
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		// Retry the request if Retries is less than maxRetries, or if maxRetries is -1. Otherwise, send the response to the channel
		if req.Retries < rl.maxRetries || rl.maxRetries == -1 {
			resp.Body.Close()
			rl.handleRateLimitedResponse(resp, regionLimiter, methodLimiter, true)

			req.Retries++
			go rl.handleRequest(req)
		} else {
			req.Response <- resp
			rl.handleRateLimitedResponse(resp, regionLimiter, methodLimiter, false)
		}

		return
	}

	if !isBadResponse(resp) && (req.Retries < rl.maxRetries || rl.maxRetries == -1) {
		resp.Body.Close()
		rl.releaseLimitersAfterDelay(regionLimiter, methodLimiter, 15*time.Second)

		req.Retries++
		go rl.handleRequest(req)
		return
	}

	// Close the response body and send an error to the channel. Do not know how to handle this response
	resp.Body.Close()
	req.Error <- fmt.Errorf("received unexpected status code %d (%s)", resp.StatusCode, req.URL)
	rl.releaseLimitersAfterDelay(regionLimiter, methodLimiter, 15*time.Second)
}

func (rl *RateLimiter) handleRateLimitedResponse(resp *http.Response, regionLimiter *RateLimit, methodLimiter *RateLimit, retryRequest bool) {
	retryAfterHeader := resp.Header.Get("Retry-After")
	rateLimitTypeHeader := resp.Header.Get("X-Rate-Limit-Type")
	retryAfter, err := strconv.Atoi(retryAfterHeader)
	if err != nil {
		retryAfter = 15
	}

	retryAfterDuration := time.Duration(retryAfter) * time.Second

	if rateLimitTypeHeader == "application" {
		regionLimiter.blockedUntil = time.Now().Add(retryAfterDuration)
	} else if rateLimitTypeHeader == "method" {
		methodLimiter.blockedUntil = time.Now().Add(retryAfterDuration)
	}

	// If we are not going to retry the request, release the limiters
	// We do this to maintain this retry request's place in the queue
	if !retryRequest {
		rl.releaseLimitersAfterDelay(regionLimiter, methodLimiter, retryAfterDuration)
	}
}

func isBadResponse(resp *http.Response) bool {
	return resp == nil || resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusUnsupportedMediaType
}

func (rl *RateLimiter) updateRateLimits(resp *http.Response, methodID MethodID, regionLimiter *RateLimit, methodLimiter *RateLimit) {
	appRateLimitHeader := resp.Header.Get("X-App-Rate-Limit")
	appRateLimitCountHeader := resp.Header.Get("X-App-Rate-Limit-Count")
	methodRateLimitHeader := resp.Header.Get("X-Method-Rate-Limit")
	methodRateLimitCountHeader := resp.Header.Get("X-Method-Rate-Limit-Count")

	if appRateLimitHeader != "" && appRateLimitCountHeader != "" {
		shortLimitInfo, longLimitInfo := getShortAndLongLimits(appRateLimitHeader)
		shortCountInfo, longCountInfo := getShortAndLongLimits(appRateLimitCountHeader)

		rl.updateRateLimit(methodID, shortLimitInfo, shortCountInfo, regionLimiter.shortLimiter, &regionLimiter.blockedUntil, rl.conserveUsage.RegionPercent, true)
		rl.updateRateLimit(methodID, longLimitInfo, longCountInfo, regionLimiter.longLimiter, &regionLimiter.blockedUntil, rl.conserveUsage.RegionPercent, true)
	} else {
		go regionLimiter.shortLimiter.ReleaseAfterDelay(15 * time.Second)
		go regionLimiter.longLimiter.ReleaseAfterDelay(15 * time.Second)
	}

	if methodRateLimitHeader != "" && methodRateLimitCountHeader != "" {
		rl.updateRateLimit(methodID, methodRateLimitHeader, methodRateLimitCountHeader, methodLimiter.shortLimiter, &methodLimiter.blockedUntil, rl.conserveUsage.MethodPercent, false)
	} else {
		go methodLimiter.shortLimiter.ReleaseAfterDelay(15 * time.Second)
	}
}

func getShortAndLongLimits(limitHeader string) (string, string) {
	limits := strings.Split(limitHeader, ",")
	return limits[0], limits[1]
}

func (rl *RateLimiter) updateRateLimit(methodID MethodID, limitInfo, countInfo string, limiterChannel *limiter.Limiter, blockedUntil *time.Time, conservePercent int, isRegionHeader bool) {
	limitSplit := strings.Split(limitInfo, ":")
	countSplit := strings.Split(countInfo, ":")

	limit, errLimit := strconv.Atoi(limitSplit[0])
	limitTimeout, errTimeout := strconv.Atoi(limitSplit[1])
	count, errCount := strconv.Atoi(countSplit[0])

	if errLimit != nil || errTimeout != nil || errCount != nil {
		if limitTimeout == 0 {
			limitTimeout = 60
		}

		go limiterChannel.ReleaseAfterDelay(time.Duration(limitTimeout) * time.Second)
		return
	}

	var limitWithConservation int = limit

	var useConservation bool = false
	if conservePercent > 0 {
		useConservation = true

		if !isRegionHeader {
			for i := 0; i < len(rl.conserveUsage.IgnoreLimits); i++ {
				if rl.conserveUsage.IgnoreLimits[i] == methodID {
					useConservation = false
					break
				}
			}
		}
	}

	if useConservation {
		limitWithConservation = limit - (limit * conservePercent / 100)
	} else {
		limitWithConservation = limit
	}

	// If the limit has been reached, block the limiter channel until the limit resets
	if count > limitWithConservation && time.Now().After(*blockedUntil) {
		*blockedUntil = time.Now().Add(time.Duration(limitTimeout) * time.Second)
	}

	// Resize the limiter channel (if necessary) to the new limit
	limiterChannel.SetCapacity(limitWithConservation)

	go limiterChannel.ReleaseAfterDelay(time.Duration(limitTimeout) * time.Second)
}
