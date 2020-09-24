// Copyright (c) 2020 Uber Technologies, Inc.
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

package index

import (
	"time"

	"github.com/m3db/m3/src/x/ident"
)

// SeriesLimitExceeded returns whether a given size exceeds the
// series limit the query options imposes, if it is enabled.
func (o QueryOptions) SeriesLimitExceeded(size int) bool {
	return o.SeriesLimit > 0 && size >= o.SeriesLimit
}

// DocsLimitExceeded returns whether a given size exceeds the
// docs limit the query options imposes, if it is enabled.
func (o QueryOptions) DocsLimitExceeded(size int) bool {
	return o.DocsLimit > 0 && size >= o.DocsLimit
}

// LimitsExceeded returns whether a given size exceeds the given limits.
func (o QueryOptions) LimitsExceeded(seriesCount, docsCount int) bool {
	return o.SeriesLimitExceeded(seriesCount) || o.DocsLimitExceeded(docsCount)
}

func (o QueryOptions) exhaustive(seriesCount, docsCount int) bool {
	return !o.SeriesLimitExceeded(seriesCount) && !o.DocsLimitExceeded(docsCount)
}

// NewWideQueryOptions creates a new wide query options, snapped to block start.
func NewWideQueryOptions(
	queryStart time.Time,
	batchSize int,
	collector chan ident.IDBatch,
	blockSize time.Duration,
	iterOpts IterationOptions,
) WideQueryOptions {
	start := queryStart.Truncate(blockSize)
	end := start.Add(blockSize)
	return WideQueryOptions{
		StartInclusive:      start,
		EndExclusive:        end,
		BatchSize:           batchSize,
		IndexBatchCollector: collector,
		IterationOptions:    iterOpts,
	}
}

// ToQueryOptions converts a WideQueryOptions to appropriate QueryOptions that
// will not enforce any limits.
func (q *WideQueryOptions) ToQueryOptions() QueryOptions {
	return QueryOptions{
		StartInclusive:   q.StartInclusive,
		EndExclusive:     q.EndExclusive,
		IterationOptions: q.IterationOptions,
	}
}
