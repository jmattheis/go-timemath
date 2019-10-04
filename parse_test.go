package timemath_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmattheis/go-timemath"
	"github.com/stretchr/testify/assert"
)

func TestParseRange(t *testing.T) {
	entries := []testEntry{
		now("2019-05-13T15:55:23Z").
			Range("now-1d").
			Expect("2019-05-12T15:55:23Z", "2019-05-12T15:55:23Z"),
		now("2019-05-13T15:55:23Z").
			Range("now-1d").
			Expect("2019-05-12T15:55:23Z", "2019-05-12T15:55:23Z"),
		now("2019-05-13T15:55:23Z").
			Range("now-1w").
			Expect("2019-05-06T15:55:23Z", "2019-05-06T15:55:23Z"),
		now("2019-05-13T15:55:23Z").
			Range("now-10d").
			Expect("2019-05-03T15:55:23Z", "2019-05-03T15:55:23Z"),
		now("2019-05-13T15:55:23Z").
			Range("now-1M").
			Expect("2019-04-13T15:55:23Z", "2019-04-13T15:55:23Z"),
		now("2019-05-13T15:55:23Z").
			Range("now-1y").
			Expect("2018-05-13T15:55:23Z", "2018-05-13T15:55:23Z"),
		now("2019-05-13T15:55:23Z").
			Range("now-5s-3m-5h-2d-2w-3M-1y").
			Expect("2018-01-27T10:52:18Z", "2018-01-27T10:52:18Z"),
		now("2019-05-13T15:55:23Z").
			Range("now/s").
			Expect("2019-05-13T15:55:23Z", "2019-05-13T15:55:23Z"),
		now("2019-05-13T15:55:23Z").
			Range("now/m").
			Expect("2019-05-13T15:55:00Z", "2019-05-13T15:55:59Z"),
		now("2019-05-13T15:55:23Z").
			Range("now/h").
			Expect("2019-05-13T15:00:00Z", "2019-05-13T15:59:59Z"),
		now("2019-05-13T15:55:23Z").
			Range("now/d").
			Expect("2019-05-13T00:00:00Z", "2019-05-13T23:59:59Z"),
		now("2019-05-15T15:55:23Z").
			Range("now/w").
			Expect("2019-05-13T00:00:00Z", "2019-05-19T23:59:59Z"),
		now("2019-05-15T15:55:23Z").
			Range("now/M").
			Expect("2019-05-01T00:00:00Z", "2019-05-31T23:59:59Z"),
		now("2019-05-15T15:55:23Z").
			Range("now/y").
			Expect("2019-01-01T00:00:00Z", "2019-12-31T23:59:59Z"),
		now("2019-05-15T15:55:23Z").
			Range("now-1w/w").
			Expect("2019-05-06T00:00:00Z", "2019-05-12T23:59:59Z"),
		now("2019-05-15T15:55:23Z").
			Range("now-1w/w-2d").
			Expect("2019-05-04T00:00:00Z", "2019-05-10T23:59:59Z"),
		now("2019-05-15T15:55:23Z").
			Range("now-1w/w+5d").
			Expect("2019-05-11T00:00:00Z", "2019-05-17T23:59:59Z"),
	}

	for _, entry := range entries {
		t.Run(fmt.Sprintf("now=%s;parse=%s", entry.now, entry.toParse), func(t *testing.T) {
			start, err := timemath.Parse(entry.now, entry.toParse, true, time.Monday)
			assert.NoError(t, err)
			end, err := timemath.Parse(entry.now, entry.toParse, false, time.Sunday)
			assert.NoError(t, err)
			assert.Equal(t, entry.expectStart, start.Format(time.RFC3339))
			assert.Equal(t, entry.expectEnd, end.Format(time.RFC3339))
		})
	}
}

func TestParseErrors(t *testing.T) {
	_, err := timemath.Parse(time.Time{}, "1d+now", true, time.Monday)
	assert.EqualError(t, err, "'now' must be at the start")
	_, err = timemath.Parse(time.Time{}, "1d", true, time.Monday)
	assert.EqualError(t, err, "value must be a valid rfc3339 date or start with 'now'")
	_, err = timemath.Parse(time.Time{}, "2019-123123", true, time.Monday)
	assert.EqualError(t, err, "value must be a valid rfc3339 date or start with 'now'")
	_, err = timemath.Parse(time.Time{}, "now-d", true, time.Monday)
	assert.EqualError(t, err, "expected number at index 4 but was d")
	_, err = timemath.Parse(time.Time{}, "now/1d", true, time.Monday)
	assert.EqualError(t, err, "expected unit y M w d h m s at index 4 but was '1'")
	_, err = timemath.Parse(time.Time{}, "now5", true, time.Monday)
	assert.EqualError(t, err, "expected operation / + - at index 3 but was '5'")
	_, err = timemath.Parse(time.Time{}, "now-", true, time.Monday)
	assert.EqualError(t, err, "expected number at the end but got nothing")
	_, err = timemath.Parse(time.Time{}, "now-5", true, time.Monday)
	assert.EqualError(t, err, "expected unit at the end but got nothing")
}

func TestPanics(t *testing.T) {
	assert.Panics(t, func() {
		timemath.Unit('x').EndOf(time.Time{}, time.Monday)
	})
	assert.Panics(t, func() {
		timemath.Unit('x').StartOf(time.Time{}, time.Monday)
	})
	assert.Panics(t, func() {
		timemath.Unit('x').Add(time.Time{}, 5)
	})
	assert.Panics(t, func() {
		timemath.Unit('x').Subtract(time.Time{}, 5)
	})
}

func now(now string) testEntry {
	return testEntry{now: parseTime(now)}
}

type testEntry struct {
	now         time.Time
	toParse     string
	expectStart string
	expectEnd   string
}

func (t testEntry) Expect(startOf string, endOf string) testEntry {
	t.expectStart = startOf
	t.expectEnd = endOf
	return t
}

func (t testEntry) Range(toParse string) testEntry {
	t.toParse = toParse
	return t
}
func parseTime(value string) time.Time {
	parse, err := time.ParseInLocation(time.RFC3339, value, time.UTC)
	if err != nil {
		panic(err)
	}
	return parse
}
