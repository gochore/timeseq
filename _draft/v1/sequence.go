package timeseq

import (
	"time"
)

// Sequence can be implemented for specific data type
type Sequence interface {
	// return length
	Len() int
	// swap items
	Swap(i, j int)
	// return time of item i
	Time(i int) time.Time
	// return Sequence[i:j]
	Slice(i, j int) Sequence
}

type sortableSequence struct {
	Sequence
}

func (s sortableSequence) Less(i, j int) bool {
	return s.Time(i).Before(s.Time(j))
}
