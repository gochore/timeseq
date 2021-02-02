// Deprecated: migrate to github.com/gochore/timeseq/v2
//
// timeseq v1 is frozen, you should use v2
package timeseq

import (
	"sort"
	"time"
)

// Deprecated: migrate to github.com/gochore/timeseq/v2
func Sort(seq Sequence) {
	if seq == nil {
		return
	}
	sort.Sort(sortableSequence{seq})
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func Range(seq Sequence, afterOrEqual, beforeOrEqual *time.Time) Sequence {
	i := 0
	if afterOrEqual != nil {
		i = sort.Search(seq.Len(), func(i int) bool {
			return !seq.Time(i).Before(*afterOrEqual)
		})
	}
	j := seq.Len()
	if beforeOrEqual != nil {
		j = sort.Search(seq.Len(), func(j int) bool {
			return !seq.Time(j).Before(*beforeOrEqual)
		})
		if j < seq.Len() && seq.Time(j).Equal(*beforeOrEqual) {
			j++
		}
	}
	return seq.Slice(i, j)
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func First(seq Sequence, afterOrEqual *time.Time) int {
	i := 0
	if afterOrEqual != nil {
		i = sort.Search(seq.Len(), func(i int) bool {
			return !seq.Time(i).Before(*afterOrEqual)
		})
		if i >= seq.Len() {
			i = -1
		}
	}
	return i
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func Last(seq Sequence, beforeOrEqual *time.Time) int {
	j := seq.Len() - 1
	if beforeOrEqual != nil {
		j = sort.Search(seq.Len(), func(i int) bool {
			return !seq.Time(i).Before(*beforeOrEqual)
		})
		if j == seq.Len() || j < seq.Len() && !seq.Time(j).Equal(*beforeOrEqual) {
			j--
		}
	}
	return j
}
