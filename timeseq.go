package timeseq

import (
	"sort"
	"time"
)

// Sort will sort sequence by time
func Sort(seq Sequence) {
	if seq == nil {
		return
	}
	sort.Sort(sortableSequence{seq})
}

func TrimBefore(seq Sequence, before time.Time) Sequence {
	i := sort.Search(seq.Len(), func(i int) bool {
		return !seq.Time(i).Before(before)
	})
	return seq.Slice(i, seq.Len())
}

func TrimAfter(seq Sequence, after time.Time) Sequence {
	j := sort.Search(seq.Len(), func(j int) bool {
		return !seq.Time(j).Before(after)
	})
	if j < seq.Len() && seq.Time(j).Equal(after) {
		j++
	}
	return seq.Slice(0, j)
}

// Get return the index of the first item with specified time
func Get(seq Sequence, t time.Time) int {
	i := sort.Search(seq.Len(), func(i int) bool {
		return !seq.Time(i).Before(t)
	})
	if i >= seq.Len() {
		i = -1
	}
	return i
}

// First return the index of the first item
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

// Last return the index of the last item
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
