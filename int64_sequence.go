package timeseq

import (
	"errors"
	"sort"
	"time"
)

// Int64Item is item of Int64Sequence
type Int64Item struct {
	Time  time.Time
	Value int64
}

// Int64Sequence is the implement of Sequence for int64
type Int64Sequence []*Int64Item

// implement of Sequence.Len
func (s Int64Sequence) Len() int {
	return len(s)
}

// implement of Sequence.Swap
func (s Int64Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// implement of Sequence.Time
func (s Int64Sequence) Time(i int) time.Time {
	return s[i].Time
}

// implement of Sequence.Slice
func (s Int64Sequence) Slice(i, j int) Sequence {
	return s[i:j]
}

// sort by time
func (s Int64Sequence) Sort() {
	Sort(s)
}

// return sub sequence, would sort sequence if it is not sorted
func (s Int64Sequence) Range(afterOrEqual, beforeOrEqual *time.Time) Int64Sequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(Int64Sequence)
}

// return the first item or nil if not exists, would sort sequence if it is not sorted
func (s Int64Sequence) First(afterOrEqual *time.Time) *Int64Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	i := First(s, afterOrEqual)
	if i < 0 {
		return nil
	}
	return s[i]
}

// return the last item or nil if not exists, would sort sequence if it is not sorted
func (s Int64Sequence) Last(beforeOrEqual *time.Time) *Int64Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	i := Last(s, beforeOrEqual)
	if i < 0 {
		return nil
	}
	return s[i]
}

// return the first item which has the max value, or nil if not exists
func (s Int64Sequence) Max() *Int64Item {
	var max *Int64Item
	for _, v := range s {
		if max == nil {
			max = v
		} else if v.Value > max.Value {
			max = v
		}
	}
	return max
}

// return the first item which has the min value, or nil if not exists
func (s Int64Sequence) Min() *Int64Item {
	var min *Int64Item
	for _, v := range s {
		if min == nil {
			min = v
		} else if v.Value < min.Value {
			min = v
		}
	}
	return min
}

// return the value's sum
func (s Int64Sequence) Sum() int64 {
	var sum int64
	for _, v := range s {
		sum += v.Value
	}
	return sum
}

// return the value's average
func (s Int64Sequence) Average() int64 {
	if len(s) == 0 {
		return 0
	}

	return int64(float64(s.Sum()) / float64(len(s)))
}

// return (pct)th percentile
func (s Int64Sequence) Percentile(pct float64) int64 {
	if pct > 1 || pct < 0 {
		panic(errors.New("percentile must be [0, 1]"))
	}

	var values []int64
	for _, v := range s {
		values = append(values, v.Value)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	if len(values) == 0 {
		return 0
	}

	index := int(float64(len(s))*pct - 1)
	if index < 0 {
		index = 0
	}

	return values[index]
}
