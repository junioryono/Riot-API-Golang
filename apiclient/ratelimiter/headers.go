package ratelimiter

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	missingHeaderDelay   = 15 * time.Second
	malformedHeaderDelay = 60 * time.Second
)

var errMalformedLimitHeader = errors.New("malformed rate limit header")

// limitWindow is one comma-separated bucket from a Riot rate-limit header.
// Headers look like "1000:60,20000:10,1200000:600" (limit:windowSeconds).
type limitWindow struct {
	limit  int
	count  int
	window time.Duration
}

func parseLimitBuckets(limitHeader, countHeader string) ([]limitWindow, error) {
	limits := strings.Split(limitHeader, ",")
	counts := strings.Split(countHeader, ",")
	if len(limits) != len(counts) || len(limits) == 0 {
		return nil, errMalformedLimitHeader
	}

	out := make([]limitWindow, 0, len(limits))
	for i := range limits {
		limitPart := strings.TrimSpace(limits[i])
		countPart := strings.TrimSpace(counts[i])
		if limitPart == "" || countPart == "" {
			return nil, errMalformedLimitHeader
		}

		limit, limitSeconds, err := parseLimitPair(limitPart)
		if err != nil {
			return nil, err
		}
		count, countSeconds, err := parseLimitPair(countPart)
		if err != nil {
			return nil, err
		}
		if limit < 1 || limitSeconds < 1 || count < 0 || countSeconds != limitSeconds {
			return nil, errMalformedLimitHeader
		}

		out = append(out, limitWindow{
			limit:  limit,
			count:  count,
			window: time.Duration(limitSeconds) * time.Second,
		})
	}
	return out, nil
}

func parseLimitPair(part string) (int, int, error) {
	value, window, ok := strings.Cut(part, ":")
	if !ok || value == "" || window == "" || strings.Contains(window, ":") {
		return 0, 0, errMalformedLimitHeader
	}
	parsedValue, errValue := strconv.Atoi(value)
	parsedWindow, errWindow := strconv.Atoi(window)
	if errValue != nil || errWindow != nil {
		return 0, 0, errMalformedLimitHeader
	}
	return parsedValue, parsedWindow, nil
}

func conservedLimit(limit, percent int) int {
	if percent <= 0 {
		return limit
	}
	reduced := limit - (limit * percent / 100)
	if reduced < 1 {
		return 1
	}
	return reduced
}

func reportLimitHeader(once *sync.Once, kind, reason, header string) {
	once.Do(func() {
		log.Printf("riot rate limit: %s header %s (%q); holding tokens with the fallback delay until a valid header arrives", kind, reason, header)
	})
}
