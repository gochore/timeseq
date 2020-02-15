# timeseq

[![Build Status](https://travis-ci.com/gochore/timeseq.svg?branch=master)](https://travis-ci.com/gochore/timeseq)
[![codecov](https://codecov.io/gh/gochore/timeseq/branch/master/graph/badge.svg)](https://codecov.io/gh/gochore/timeseq)
[![Go Report Card](https://goreportcard.com/badge/github.com/gochore/timeseq)](https://goreportcard.com/report/github.com/gochore/timeseq)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gochore/timeseq)](https://github.com/gochore/timeseq/blob/master/go.mod)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/gochore/timeseq)](https://github.com/gochore/timeseq/releases)

Time sequence.

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/gochore/timeseq"
)

func main() {
	now := time.Now()
	seq := timeseq.Int64Sequence{
		{
			Time:  now,
			Value: 0,
		},
		{
			Time:  now.Add(time.Second),
			Value: 1,
		},
		{
			Time:  now.Add(-time.Second),
			Value: 2,
		},
	}
	seq = append(seq, timeseq.Int64Item{
		Time:  now.Add(-2 * time.Second),
		Value: -1,
	})

	// sort by time
	seq.Sort()

	// get the first one
	fmt.Println(seq.First(nil))
	// get the last one before now
	fmt.Println(seq.Last(&now))

	// get the sub sequence after now
	subSeq := seq.Range(&now, nil)
	// get the sub sequence's length
	fmt.Println(subSeq.Len())
	// get the first one of the sub sequence
	fmt.Println(subSeq.First(nil))
}
```
