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
	seq = append(seq, &timeseq.Int64Item{
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
