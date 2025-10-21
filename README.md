# timeseq

[![Go Reference](https://pkg.go.dev/badge/github.com/gochore/timeseq/v3.svg)](https://pkg.go.dev/github.com/gochore/timeseq/v3)
[![Actions](https://github.com/gochore/timeseq/actions/workflows/check.yaml/badge.svg)](https://github.com/gochore/timeseq/actions)
[![codecov](https://codecov.io/gh/gochore/timeseq/branch/master/graph/badge.svg)](https://codecov.io/gh/gochore/timeseq)
[![Go Report Card](https://goreportcard.com/badge/github.com/gochore/timeseq)](https://goreportcard.com/report/github.com/gochore/timeseq)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gochore/timeseq)](https://github.com/gochore/timeseq/blob/master/go.mod)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/gochore/timeseq)](https://github.com/gochore/timeseq/releases)

Time sequence.

This package defines `Seq[T]` for time series data,
and provides useful methods to deal with them, you can `Merge`, `Filter` or `Aggregate` your data.

## Install

```bash
go get github.com/gochore/timeseq/v3
```

## Example

See more in [doc](https://pkg.go.dev/github.com/gochore/timeseq/v3).

```go
package main

import (
	"fmt"
	"time"

	"github.com/gochore/timeseq/v3"
)

func main() {
	now := time.Now()

	// define a data point
	point := timeseq.NewPoint(now, 1)

	// define a data series
	slice := []timeseq.Point[int]{
		{
			Time:  now.Add(time.Second),
			Value: 1,
		},
		{
			Time:  now.Add(2 * time.Second),
			Value: 2,
		},
	}

	// you can append and modify slice
	slice = append(slice, point)
	slice[0].Value = 100
	slice[1].Time = now.Add(time.Hour)

	// define a seq according to the slice, it will copy and sort data
	seq := timeseq.NewSeq(slice)

	// now you can not modify or add elem to seq, seq has protected slice inside
	point = seq.Index(0)
	point.Value = 100 // not work, it does not change the data in seq

	// for better performance, you can use existing slice as inside data
	seq = timeseq.WrapSeq(slice)

	// it should be noted that you should not modify the slice any longer
	// slice[0].Value = 0 // please don't do that!

	// now you enjoy the convenience it brings

	// get time of the first one
	fmt.Println(seq.First().Time)
	// get value of the last one
	fmt.Println(seq.Last().Value)

	// get the sub sequence after now
	seq = seq.Range(now, time.Time{})
	// get the sub sequence's length
	fmt.Println(seq.Len())
	// get the first one of the sub sequence
	fmt.Println(seq.Max().Value)

	// traverse
	seq.Traverse(func(i int, v timeseq.Point[int]) (stop bool) {
		fmt.Println(i, v.Time, v.Value)
		return false
	})

	// merge data
	seq2 := timeseq.WrapSeq([]timeseq.Point[int]{
		{
			Time:  now.Add(time.Second),
			Value: 1,
		},
		{
			Time:  now.Add(-time.Second),
			Value: 2,
		},
	})
	newSeq := timeseq.Merge(seq, seq2, func(t time.Time, v1, v2 *int) *int {
		if v1 == nil {
			return v2
		}
		if v2 == nil {
			return v1
		}
		ret := *v1 + *v2
		return &ret
	})

	fmt.Println(newSeq.Points())
}
```

