// Copyright (c) 2021 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tally

import (
	"runtime"
	"testing"
	"time"
)

func BenchmarkCounterInc(b *testing.B) {
	c := &counter{}
	for n := 0; n < b.N; n++ {
		c.Inc(1)
	}
}

func BenchmarkReportCounterNoData(b *testing.B) {
	c := &counter{}
	for n := 0; n < b.N; n++ {
		c.report("foo", nil, NullStatsReporter)
	}
}

func BenchmarkReportCounterWithData(b *testing.B) {
	c := &counter{}
	for n := 0; n < b.N; n++ {
		c.Inc(1)
		c.report("foo", nil, NullStatsReporter)
	}
}

func BenchmarkGaugeSet(b *testing.B) {
	g := &gauge{}
	for n := 0; n < b.N; n++ {
		g.Update(42)
	}
}

func BenchmarkReportGaugeNoData(b *testing.B) {
	g := &gauge{}
	for n := 0; n < b.N; n++ {
		g.report("bar", nil, NullStatsReporter)
	}
}

func BenchmarkReportGaugeWithData(b *testing.B) {
	g := &gauge{}
	for n := 0; n < b.N; n++ {
		g.Update(73)
		g.report("bar", nil, NullStatsReporter)
	}
}

func BenchmarkTimerStopwatch(b *testing.B) {
	t := &timer{
		name:     "bencher",
		tags:     nil,
		reporter: NullStatsReporter,
	}
	for n := 0; n < b.N; n++ {
		t.Start().Stop() // start and stop
	}
}

func BenchmarkTimerReport(b *testing.B) {
	t := &timer{
		name:     "bencher",
		tags:     nil,
		reporter: NullStatsReporter,
	}
	for n := 0; n < b.N; n++ {
		start := time.Now()
		t.Record(time.Since(start))
	}
}

func BenchmarkHistogramRecordDuration(b *testing.B) {
	h := newHistogram(
		durationHistogramType,
		"foo",
		map[string]string{"foo": "bar"},
		noopCachedReporter{},
		newBucketStorage(durationHistogramType, MustMakeLinearDurationBuckets(0, time.Second, 20)),
		nil,
	)
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		h.RecordDuration(time.Duration(n%20) * time.Second)
	}
	runtime.KeepAlive(h)
}

func BenchmarkHistogramRecordValue(b *testing.B) {
	h := newHistogram(
		valueHistogramType,
		"foo",
		map[string]string{"foo": "bar"},
		noopCachedReporter{},
		newBucketStorage(valueHistogramType, MustMakeLinearValueBuckets(0.0, 1.0, 20)),
		nil,
	)
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		h.RecordValue(float64(n % 20))
	}
	runtime.KeepAlive(h)
}
