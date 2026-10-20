package ratelimiter

import (
	"bytes"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLimitBucketsMethodHeader(t *testing.T) {
	buckets, err := parseLimitBuckets("1000:60,20000:10,1200000:600", "1:60,2:10,3:600")
	require.NoError(t, err)
	require.Equal(t, []limitWindow{
		{limit: 1000, count: 1, window: 60 * time.Second},
		{limit: 20000, count: 2, window: 10 * time.Second},
		{limit: 1200000, count: 3, window: 600 * time.Second},
	}, buckets)
}

func TestParseLimitBucketsAppHeader(t *testing.T) {
	buckets, err := parseLimitBuckets("500:10,30000:600", "4:10,5:600")
	require.NoError(t, err)
	require.Equal(t, []limitWindow{
		{limit: 500, count: 4, window: 10 * time.Second},
		{limit: 30000, count: 5, window: 600 * time.Second},
	}, buckets)
}

func TestParseLimitBucketsSingleBucket(t *testing.T) {
	buckets, err := parseLimitBuckets("20:1", "1:1")
	require.NoError(t, err)
	require.Equal(t, []limitWindow{
		{limit: 20, count: 1, window: time.Second},
	}, buckets)
}

func TestParseLimitBucketsRejectsMalformed(t *testing.T) {
	cases := []struct {
		limit string
		count string
	}{
		{limit: "1000:60,20000", count: "1:60,2:10"},
		{limit: "1000:60,20000:10", count: "1:60"},
		{limit: "500:10,30000:600", count: "1:10,2:10"},
		{limit: "", count: ""},
		{limit: "0:10", count: "1:10"},
		{limit: "10:0", count: "1:0"},
		{limit: "10:10,", count: "1:10,"},
	}
	for _, tc := range cases {
		_, err := parseLimitBuckets(tc.limit, tc.count)
		assert.ErrorIs(t, err, errMalformedLimitHeader, "limit %q count %q", tc.limit, tc.count)
	}
}

func TestConservedLimitFloorsAtOne(t *testing.T) {
	assert.Equal(t, 450, conservedLimit(500, 10))
	assert.Equal(t, 1000, conservedLimit(1000, 0))
	assert.Equal(t, 1, conservedLimit(1, 100))
}

func TestReportLimitHeaderLogsOnce(t *testing.T) {
	var buf bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(previous) })

	var once sync.Once
	reportLimitHeader(&once, "method", "malformed", "1000:60,20000")
	reportLimitHeader(&once, "method", "malformed", "1000:60,20000")

	assert.Equal(t, 1, bytes.Count(buf.Bytes(), []byte("method header malformed")))
}
