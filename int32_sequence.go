// Code generated by cmd/generate. DO NOT EDIT.
//
// Deprecated: migrate to github.com/gochore/timeseq/v2
//
// timeseq v1 is frozen, you should use v2
package timeseq

import (
	"errors"
	"sort"
	"time"
)

// Deprecated: migrate to github.com/gochore/timeseq/v2
type Int32Item struct {
	Time  time.Time
	Value int32
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
type Int32Sequence []Int32Item

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Len() int {
	return len(s)
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Time(i int) time.Time {
	return s[i].Time
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Slice(i, j int) Sequence {
	return s[i:j]
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Sort() {
	Sort(s)
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Range(afterOrEqual, beforeOrEqual *time.Time) Int32Sequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(Int32Sequence)
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) First(afterOrEqual *time.Time) *Int32Item {
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

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Last(beforeOrEqual *time.Time) *Int32Item {
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

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Max() *Int32Item {
	var max *Int32Item
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

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Min() *Int32Item {
	var min *Int32Item
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

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Sum() int32 {
	var sum int32
	for _, v := range s {
		sum += v.Value
	}
	return sum
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Average() int32 {
	if len(s) == 0 {
		return 0
	}

	return int32(float64(s.Sum()) / float64(len(s)))
}

// Deprecated: migrate to github.com/gochore/timeseq/v2
func (s Int32Sequence) Percentile(pct float64) int32 {
	if pct > 1 || pct < 0 {
		panic(errors.New("percentile must be [0, 1]"))
	}

	var values []int32
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

// Deprecated: migrate to github.com/gochore/timeseq/v2
func MergeInt32(seq1, seq2 Int32Sequence, fn func(item1, item2 *Int32Item) *Int32Item) Int32Sequence {
	if fn == nil {
		return nil
	}

	var ret Int32Sequence
	for i1, i2 := 0, 0; i1 < seq1.Len() || i2 < seq2.Len(); {
		var item *Int32Item
		switch {
		case i1 == seq1.Len():
			v2 := seq2[i2]
			item = fn(nil, &v2)
			i2++
		case i2 == seq2.Len():
			v1 := seq1[i1]
			item = fn(&v1, nil)
			i1++
		case seq1[i1].Time.Equal(seq2[i2].Time):
			v1 := seq1[i1]
			v2 := seq2[i2]
			item = fn(&v1, &v2)
			i1++
			i2++
		case seq1[i1].Time.Before(seq2[i2].Time):
			v1 := seq1[i1]
			item = fn(&v1, nil)
			i1++
		case seq1[i1].Time.After(seq2[i2].Time):
			v2 := seq2[i2]
			item = fn(nil, &v2)
			i2++
		}
		if item != nil {
			ret = append(ret, *item)
		}
	}

	Sort(ret)
	return ret
}
