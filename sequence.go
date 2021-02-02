// Deprecated: migrate to github.com/gochore/timeseq/v2
//
// timeseq v1 is frozen, you should use v2
package timeseq

import (
	"time"
)

// Deprecated: migrate to github.com/gochore/timeseq/v2
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
