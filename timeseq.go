package timeseq

import (
	"sort"
	"time"
)

// sort sequence by time
func Sort(seq Sequence) {
	if seq == nil {
		return
	}
	sort.Sort(sortableSequence{seq})
}

// return sub sequence
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

// return the index of the first item
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

// return the index of the last item
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
