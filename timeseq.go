package timeseq

import (
	"sort"
	"time"
)

func Sort(seq Sequence) {
	sort.Sort(sortableSequence{seq})
}

func Range(seq Sequence, afterOrEqual, beforeOrEqual time.Time) Sequence {
	i := sort.Search(seq.Len(), func(i int) bool {
		return !seq.Time(i).Before(afterOrEqual)
	})
	j := sort.Search(seq.Len(), func(i int) bool {
		return !seq.Time(i).Before(beforeOrEqual)
	})
	if i < seq.Len() && !seq.Time(i).Equal(afterOrEqual) {
		i++
	}
	if j < seq.Len() && !seq.Time(i).Equal(beforeOrEqual) {
		j--
	}
	return seq.Slice(i, j+1)
}

func First(seq Sequence, afterOrEqual time.Time) int {
	i := sort.Search(seq.Len(), func(i int) bool {
		return !seq.Time(i).Before(afterOrEqual)
	})
	if i < seq.Len() && seq.Time(i).Equal(afterOrEqual) {
		return i
	}
	i++
	if i < seq.Len() {
		return i
	}
	return -1
}

func Last(seq Sequence, beforeOrEqual time.Time) int {
	i := sort.Search(seq.Len(), func(i int) bool {
		return !seq.Time(i).Before(beforeOrEqual)
	})
	if i < seq.Len() {
		return i
	}
	return -1
}
