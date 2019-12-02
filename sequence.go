package timeseq

import (
	"time"
)

type Sequence interface {
	Len() int
	Swap(i, j int)
	Time(i int) time.Time
	Slice(i, j int) Sequence
}

type sortableSequence struct {
	Sequence
}

func (s sortableSequence) Less(i, j int) bool {
	return s.Time(i).Before(s.Time(j))
}
