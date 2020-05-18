// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"errors"
	"sort"
	"time"
)

// Int8Item is item of Int8Sequence
type Int8Item struct {
	Time  time.Time
	Value int8
}

// Int8Sequence is the implement of Sequence for int8
type Int8Sequence []Int8Item

// Len implements Sequence.Len
func (s Int8Sequence) Len() int {
	return len(s)
}

// Swap implements Sequence.Swap
func (s Int8Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Sequence.Time
func (s Int8Sequence) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Sequence.Slice
func (s Int8Sequence) Slice(i, j int) Sequence {
	return s[i:j]
}

// Sort will sort sequence by time
func (s Int8Sequence) Sort() {
	Sort(s)
}

// Range return sub sequence, would sort sequence if it is not sorted
func (s Int8Sequence) Range(afterOrEqual, beforeOrEqual *time.Time) Int8Sequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(Int8Sequence)
}

// First return the first item or nil if not exists, would sort sequence if it is not sorted
func (s Int8Sequence) First(afterOrEqual *time.Time) *Int8Item {
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
func (s Int8Sequence) Last(beforeOrEqual *time.Time) *Int8Item {
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
func (s Int8Sequence) Max() *Int8Item {
	var max *Int8Item
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
func (s Int8Sequence) Min() *Int8Item {
	var min *Int8Item
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
func (s Int8Sequence) Sum() int8 {
	var sum int8
	for _, v := range s {
		sum += v.Value
	}
	return sum
}

// Average return the value's average
func (s Int8Sequence) Average() int8 {
	if len(s) == 0 {
		return 0
	}

	return int8(float64(s.Sum()) / float64(len(s)))
}

// Percentile return (pct)th percentile
func (s Int8Sequence) Percentile(pct float64) int8 {
	if pct > 1 || pct < 0 {
		panic(errors.New("percentile must be [0, 1]"))
	}

	var values []int8
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

// MergeInt8 merge two int8} seuquence into one
func MergeInt8(seq1, seq2 Int8Sequence, fn func(item1, item2 *Int8Item) *Int8Item) Int8Sequence {
	if fn == nil {
		return nil
	}

	var ret Int8Sequence
	for i1, i2 := 0, 0; i1 < seq1.Len() || i2 < seq2.Len(); {
		var item *Int8Item
		switch {
		case i1 == seq1.Len():
			item = fn(nil, &seq2[i2])
			i2++
		case i2 == seq2.Len():
			item = fn(&seq1[i1], nil)
			i1++
		case seq1[i1].Time.Equal(seq2[i2].Time):
			item = fn(&seq1[i1], &seq2[i2])
			i1++
			i2++
		case seq1[i1].Time.Before(seq2[i2].Time):
			item = fn(&seq1[i1], nil)
			i1++
		case seq1[i1].Time.After(seq2[i2].Time):
			item = fn(nil, &seq2[i2])
			i2++
		}
		if item != nil {
			ret = append(ret, *item)
		}
	}

	Sort(ret)
	return ret
}