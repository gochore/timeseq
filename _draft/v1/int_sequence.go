package timeseq

import (
	"errors"
	"sort"
	"time"
)

// IntItem is item of IntSequence
type IntItem struct {
	Time  time.Time
	Value int
}

// IntSequence is the implement of Sequence for int
type IntSequence []IntItem

// Len implements Sequence.Len
func (s IntSequence) Len() int {
	return len(s)
}

// Swap implements Sequence.Swap
func (s IntSequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Sequence.Time
func (s IntSequence) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Sequence.Slice
func (s IntSequence) Slice(i, j int) Sequence {
	return s[i:j]
}

// Sort will sort sequence by time
func (s IntSequence) Sort() {
	Sort(s)
}

// Range return sub sequence, would sort sequence if it is not sorted
func (s IntSequence) Range(afterOrEqual, beforeOrEqual *time.Time) IntSequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(IntSequence)
}

// First return the first item or nil if not exists, would sort sequence if it is not sorted
func (s IntSequence) First(afterOrEqual *time.Time) *IntItem {
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
func (s IntSequence) Last(beforeOrEqual *time.Time) *IntItem {
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
func (s IntSequence) Max() *IntItem {
	var max *IntItem
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
func (s IntSequence) Min() *IntItem {
	var min *IntItem
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
func (s IntSequence) Sum() int {
	var sum int
	for _, v := range s {
		sum += v.Value
	}
	return sum
}

// Average return the value's average
func (s IntSequence) Average() int {
	if len(s) == 0 {
		return 0
	}

	return int(float64(s.Sum()) / float64(len(s)))
}

// Percentile return (pct)th percentile
func (s IntSequence) Percentile(pct float64) int {
	if pct > 1 || pct < 0 {
		panic(errors.New("percentile must be [0, 1]"))
	}

	var values []int
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

// MergeInt merge two int} seuquence into one
func MergeInt(seq1, seq2 IntSequence, fn func(item1, item2 *IntItem) *IntItem) IntSequence {
	if fn == nil {
		return nil
	}

	var ret IntSequence
	for i1, i2 := 0, 0; i1 < seq1.Len() || i2 < seq2.Len(); {
		var item *IntItem
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
