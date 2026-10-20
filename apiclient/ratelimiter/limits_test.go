package ratelimiter

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateRateLimitsResizesEveryBucket(t *testing.T) {
	rl := NewRateLimiter(nil, "key")
	region := newRateLimit()
	method := newRateLimit()

	assert.Equal(t, []int{initialLimit}, region.capacities())
	assert.Equal(t, []int{initialLimit}, method.capacities())

	resp := rateLimitResponse(
		"500:10,30000:600",
		"1:10,2:600",
		"1000:60,20000:10,1200000:600",
		"1:60,2:10,3:600",
	)
	rl.updateRateLimits(resp, GetAccountByPuuid, region, method, nil, nil)

	assert.Equal(t, []int{500, 30000}, region.capacities())
	assert.Equal(t, []time.Duration{10 * time.Second, 600 * time.Second}, region.releaseWindows())
	assert.Equal(t, []int{1000, 20000, 1200000}, method.capacities())
	assert.Equal(t, []time.Duration{60 * time.Second, 10 * time.Second, 600 * time.Second}, method.releaseWindows())
}

func TestUpdateRateLimitsAppliesConservationAndIgnoreList(t *testing.T) {
	rl := NewRateLimiter(nil, "key")
	require.NoError(t, rl.SetUsageConservation(ConserveUsage{
		RegionPercent: 10,
		MethodPercent: 50,
		IgnoreLimits:  []MethodID{GetLeagueEntriesByPuuid},
	}))

	conservedRegion := newRateLimit()
	conservedMethod := newRateLimit()
	rl.updateRateLimits(rateLimitResponse(
		"500:10,30000:600",
		"1:10,1:600",
		"1000:60,20000:10,1200000:600",
		"1:60,1:10,1:600",
	), GetAccountByPuuid, conservedRegion, conservedMethod, nil, nil)

	assert.Equal(t, []int{450, 27000}, conservedRegion.capacities())
	assert.Equal(t, []int{500, 10000, 600000}, conservedMethod.capacities())

	ignoredRegion := newRateLimit()
	ignoredMethod := newRateLimit()
	rl.updateRateLimits(rateLimitResponse(
		"500:10,30000:600",
		"1:10,1:600",
		"20000:10,1200000:600",
		"1:10,1:600",
	), GetLeagueEntriesByPuuid, ignoredRegion, ignoredMethod, nil, nil)

	assert.Equal(t, []int{450, 27000}, ignoredRegion.capacities())
	assert.Equal(t, []int{20000, 1200000}, ignoredMethod.capacities())
}

func TestUpdateRateLimitsKeepsFallbackCapacityWhenMalformed(t *testing.T) {
	rl := NewRateLimiter(nil, "key")
	region := newRateLimit()
	method := newRateLimit()

	rl.updateRateLimits(rateLimitResponse(
		"500:10,30000:600",
		"1:10,1:600",
		"1000:60,20000",
		"1:60,1:10",
	), GetAccountByPuuid, region, method, nil, nil)

	assert.Equal(t, []int{500, 30000}, region.capacities())
	assert.Equal(t, []int{initialLimit}, method.capacities())
	assert.Equal(t, []time.Duration{0}, method.releaseWindows())
}

func TestUpdateRateLimitsAcceptsExtraAppBucket(t *testing.T) {
	rl := NewRateLimiter(nil, "key")
	region := newRateLimit()
	method := newRateLimit()

	rl.updateRateLimits(rateLimitResponse(
		"100:1,200:10,300:60",
		"1:1,1:10,1:60",
		"20:1",
		"1:1",
	), GetLeagueEntriesByPuuid, region, method, nil, nil)

	assert.Equal(t, []int{100, 200, 300}, region.capacities())
	assert.Equal(t, []int{20}, method.capacities())
}

func TestUpdateRateLimitsBlocksWhenCountExceedsCapacity(t *testing.T) {
	rl := NewRateLimiter(nil, "key")
	region := newRateLimit()
	method := newRateLimit()
	before := time.Now()

	rl.updateRateLimits(rateLimitResponse(
		"500:10,30000:600",
		"501:10,1:600",
		"1000:60",
		"1:60",
	), GetAccountByPuuid, region, method, nil, nil)

	region.mu.Lock()
	blockedUntil := region.blockedUntil
	region.mu.Unlock()
	assert.True(t, blockedUntil.After(before.Add(9*time.Second)), "blockedUntil = %s", blockedUntil)
}

func rateLimitResponse(appLimit, appCount, methodLimit, methodCount string) *http.Response {
	header := make(http.Header)
	header.Set("X-App-Rate-Limit", appLimit)
	header.Set("X-App-Rate-Limit-Count", appCount)
	header.Set("X-Method-Rate-Limit", methodLimit)
	header.Set("X-Method-Rate-Limit-Count", methodCount)
	return &http.Response{Header: header, StatusCode: http.StatusOK}
}
