// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"errors"
	"sort"
	"time"
)

// Float32Item is item of Float32Sequence
type Float32Item struct {
	Time  time.Time
	Value float32
}

// Float32Sequence is the implement of Sequence for float32
type Float32Sequence []Float32Item

// Len implements Sequence.Len
func (s Float32Sequence) Len() int {
	return len(s)
}

// Swap implements Sequence.Swap
func (s Float32Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Sequence.Time
func (s Float32Sequence) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Sequence.Slice
func (s Float32Sequence) Slice(i, j int) Sequence {
	return s[i:j]
}

// Sort will sort sequence by time
func (s Float32Sequence) Sort() {
	Sort(s)
}

// Range return sub sequence, would sort sequence if it is not sorted
func (s Float32Sequence) Range(afterOrEqual, beforeOrEqual *time.Time) Float32Sequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(Float32Sequence)
}

// First return the first item or nil if not exists, would sort sequence if it is not sorted
func (s Float32Sequence) First(afterOrEqual *time.Time) *Float32Item {
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
func (s Float32Sequence) Last(beforeOrEqual *time.Time) *Float32Item {
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
func (s Float32Sequence) Max() *Float32Item {
	var max *Float32Item
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
func (s Float32Sequence) Min() *Float32Item {
	var min *Float32Item
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
func (s Float32Sequence) Sum() float32 {
	var sum float32
	for _, v := range s {
		sum += v.Value
	}
	return sum
}

// Average return the value's average
func (s Float32Sequence) Average() float32 {
	if len(s) == 0 {
		return 0
	}

	return float32(float64(s.Sum()) / float64(len(s)))
}

// Percentile return (pct)th percentile
func (s Float32Sequence) Percentile(pct float64) float32 {
	if pct > 1 || pct < 0 {
		panic(errors.New("percentile must be [0, 1]"))
	}

	var values []float32
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

// MergeFloat32 merge two float32} seuquence into one
func MergeFloat32(seq1, seq2 Float32Sequence, fn func(item1, item2 *Float32Item) *Float32Item) Float32Sequence {
	if fn == nil {
		return nil
	}

	var ret Float32Sequence
	for i1, i2 := 0, 0; i1 < seq1.Len() || i2 < seq2.Len(); {
		var item *Float32Item
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
