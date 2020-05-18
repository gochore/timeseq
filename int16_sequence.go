// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"errors"
	"sort"
	"time"
)

// Int16Item is item of Int16Sequence
type Int16Item struct {
	Time  time.Time
	Value int16
}

// Int16Sequence is the implement of Sequence for int16
type Int16Sequence []Int16Item

// Len implements Sequence.Len
func (s Int16Sequence) Len() int {
	return len(s)
}

// Swap implements Sequence.Swap
func (s Int16Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Sequence.Time
func (s Int16Sequence) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Sequence.Slice
func (s Int16Sequence) Slice(i, j int) Sequence {
	return s[i:j]
}

// Sort will sort sequence by time
func (s Int16Sequence) Sort() {
	Sort(s)
}

// Range return sub sequence, would sort sequence if it is not sorted
func (s Int16Sequence) Range(afterOrEqual, beforeOrEqual *time.Time) Int16Sequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(Int16Sequence)
}

// First return the first item or nil if not exists, would sort sequence if it is not sorted
func (s Int16Sequence) First(afterOrEqual *time.Time) *Int16Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	i := First(s, afterOrEqual)
	if i < 0 {
		return nil
	}
	ret := s[i]
	return &ret
}

// Last return the last item or nil if not exists, would sort sequence if it is not sorted
func (s Int16Sequence) Last(beforeOrEqual *time.Time) *Int16Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	i := Last(s, beforeOrEqual)
	if i < 0 {
		return nil
	}
	ret := s[i]
	return &ret
}

// Max return the first item which has the max value, or nil if not exists
func (s Int16Sequence) Max() *Int16Item {
	var max *Int16Item
	for i, v := range s {
		if max == nil {
			max = &s[i]
		} else if v.Value > max.Value {
			max = &s[i]
		}
	}
	if max != nil {
		value := *max
		max = &value
	}
	return max
}

// Min return the first item which has the min value, or nil if not exists
func (s Int16Sequence) Min() *Int16Item {
	var min *Int16Item
	for i, v := range s {
		if min == nil {
			min = &s[i]
		} else if v.Value < min.Value {
			min = &s[i]
		}
	}
	if min != nil {
		value := *min
		min = &value
	}
	return min
}

// Sum return the value's sum
func (s Int16Sequence) Sum() int16 {
	var sum int16
	for _, v := range s {
		sum += v.Value
	}
	return sum
}

// Average return the value's average
func (s Int16Sequence) Average() int16 {
	if len(s) == 0 {
		return 0
	}

	return int16(float64(s.Sum()) / float64(len(s)))
}

// Percentile return (pct)th percentile
func (s Int16Sequence) Percentile(pct float64) int16 {
	if pct > 1 || pct < 0 {
		panic(errors.New("percentile must be [0, 1]"))
	}

	var values []int16
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
