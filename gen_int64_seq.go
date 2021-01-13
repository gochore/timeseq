// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// Int64 is a time point with int64 value inside
type Int64 struct {
	Time  time.Time
	Value int64
}

// IsZero return if time and value are both zero
func (v Int64) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

// IsZero return if time and value are both equal
func (v Int64) Equal(n Int64) bool {
	return v.Value == n.Value && v.Time.Equal(n.Time)
}

// Int64s is a alias of Int64 slice
type Int64s []Int64

// Len implements Interface.Len()
func (s Int64s) Len() int {
	return len(s)
}

// Swap implements Interface.Swap()
func (s Int64s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Interface.Time()
func (s Int64s) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Interface.Slice()
func (s Int64s) Slice(i, j int) Interface {
	return s[i:j]
}

// Int64Seq is a wrapper with useful methods of Int64 slice
type Int64Seq struct {
	slice Int64s

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[int64][]int
	valueSlice []int
}

// NewInt64Seq return *Int64Seq with copied slice inside
func NewInt64Seq(slice Int64s) *Int64Seq {
	temp := make(Int64s, len(slice))
	copy(temp, slice)
	return WrapInt64Seq(temp)
}

// WrapInt64Seq return *Int64Seq with origin slice inside
func WrapInt64Seq(slice Int64s) *Int64Seq {
	if !IsSorted(slice) {
		Sort(slice)
	}
	return newInt64Seq(slice)
}

func newInt64Seq(slice Int64s) *Int64Seq {
	ret := &Int64Seq{
		slice: slice,
	}
	return ret
}

func (s *Int64Seq) buildIndex() {
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[int64][]int, len(s.slice))
		valueSlice := s.valueSlice[:0]
		for i, v := range s.slice {
			k := newTimeKey(v.Time)
			timeIndex[k] = append(timeIndex[k], i)
			valueIndex[v.Value] = append(valueIndex[v.Value], i)
			valueSlice = append(valueSlice, i)
		}
		sort.SliceStable(valueSlice, func(i, j int) bool {
			return s.slice[valueSlice[i]].Value < s.slice[valueSlice[j]].Value
		})
		s.timeIndex = timeIndex
		s.valueIndex = valueIndex
		s.valueSlice = valueSlice
	})
}

func (s *Int64Seq) resetIndex() {
	s.indexOnce = sync.Once{}
}

// Int64s return a replica of inside slice
func (s *Int64Seq) Int64s() Int64s {
	slice := make(Int64s, len(s.slice))
	copy(slice, s.slice)
	return slice
}

// Len return length of inside slice
func (s *Int64Seq) Len() int {
	return len(s.slice)
}

// Index return element of inside slice, return zero if index is out of range
func (s *Int64Seq) Index(i int) Int64 {
	if i < 0 || i >= len(s.slice) {
		return Int64{}
	}
	return s.slice[i]
}

// Time return the first element with time t, return zero if not found
func (s *Int64Seq) Time(t time.Time) Int64 {
	got := s.MTime(t)
	if len(got) == 0 {
		return Int64{}
	}
	return got[0]
}

// MTime return all elements with time t, return nil if not found
func (s *Int64Seq) MTime(t time.Time) Int64s {
	s.buildIndex()
	index := s.timeIndex[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Int64s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

// Value return the first element with value v, return zero if not found
func (s *Int64Seq) Value(v int64) Int64 {
	got := s.MValue(v)
	if len(got) == 0 {
		return Int64{}
	}
	return got[0]
}

// MValue return all elements with value v, return nil if not found
func (s *Int64Seq) MValue(v int64) Int64s {
	s.buildIndex()
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Int64s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

// Traverse call fn for every element one by one, break if fn return true
func (s *Int64Seq) Traverse(fn func(i int, v Int64) (stop bool)) {
	for i, v := range s.slice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

// Sum return sum of all values
func (s *Int64Seq) Sum() int64 {
	var ret int64
	for _, v := range s.slice {
		ret += v.Value
	}
	return ret
}

// Count return count of elements, same as Len
func (s *Int64Seq) Count() int {
	return s.Len()
}

// Max return the element with max value, return zero if empty
func (s *Int64Seq) Max() Int64 {
	var max Int64
	found := false
	for _, v := range s.slice {
		if !found {
			max = v
			found = true
		} else if v.Value > max.Value {
			max = v
		}
	}
	return max
}

// Max return the element with min value, return zero if empty
func (s *Int64Seq) Min() Int64 {
	var min Int64
	found := false
	for _, v := range s.slice {
		if !found {
			min = v
			found = true
		} else if v.Value < min.Value {
			min = v
		}
	}
	return min
}

// First return the first element, return zero if empty
func (s *Int64Seq) First() Int64 {
	if len(s.slice) == 0 {
		return Int64{}
	}
	return s.slice[0]
}

// Last return the last element, return zero if empty
func (s *Int64Seq) Last() Int64 {
	if len(s.slice) == 0 {
		return Int64{}
	}
	return s.slice[len(s.slice)-1]
}

// Percentile return the element matched with percentile pct, return zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Int64Seq) Percentile(pct float64) Int64 {
	s.buildIndex()
	if len(s.slice) == 0 {
		return Int64{}
	}
	if pct > 1 {
		pct = 1
	}
	if pct < 0 {
		pct = 0
	}
	i := int(float64(len(s.slice))*pct - 1)
	if i < 0 {
		i = 0
	}
	return s.slice[s.valueSlice[i]]
}

// Range return a sub *Int64Seq with specified interval
func (s *Int64Seq) Range(interval Interval) *Int64Seq {
	slice := Range(s.slice, interval).(Int64s)
	return newInt64Seq(slice)
}

// Merge merge slices to inside slice according to the specified rule
func (s *Int64Seq) Merge(fn func(t time.Time, v1, v2 *int64) *int64, slices ...Int64s) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	if len(slices) == 0 {
		return nil
	}

	slice1 := s.slice
	for _, slice2 := range slices {
		if !IsSorted(slice2) {
			temp := make(Int64s, len(slice2))
			copy(temp, slice2)
			Sort(temp)
			slice2 = temp
		}
		var got Int64s
		for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
			var (
				t time.Time
				v *int64
			)
			switch {
			case i1 == len(slice1):
				t = slice2[i2].Time
				v2 := slice2[i2].Value
				v = fn(t, nil, &v2)
				i2++
			case i2 == len(slice2):
				t = slice1[i1].Time
				v1 := slice1[i1].Value
				v = fn(t, &v1, nil)
				i1++
			case slice1[i1].Time.Equal(slice2[i2].Time):
				t = slice1[i1].Time
				v1 := slice1[i1].Value
				v2 := slice2[i2].Value
				v = fn(t, &v1, &v2)
				i1++
				i2++
			case slice1[i1].Time.Before(slice2[i2].Time):
				t = slice1[i1].Time
				v1 := slice1[i1].Value
				v = fn(t, &v1, nil)
				i1++
			case slice1[i1].Time.After(slice2[i2].Time):
				t = slice2[i2].Time
				v2 := slice2[i2].Value
				v = fn(t, nil, &v2)
				i2++
			}
			if v != nil {
				got = append(got, Int64{
					Time:  t,
					Value: *v,
				})
			}
		}
		slice1 = got
	}

	s.slice = slice1
	s.resetIndex()
	return nil
}

// Aggregate aggregate inside slice according to the specified rule
func (s *Int64Seq) Aggregate(fn func(t time.Time, slice Int64s) *int64, duration time.Duration, interval Interval) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	got := Int64s{}
	temp := Int64s{}

	if duration <= 0 {
		for i := 0; i < s.Len(); {
			t := s.slice[i].Time
			if !interval.Contain(t) {
				i++
				continue
			}
			temp = temp[:0]
			for i < s.slice.Len() && t.Equal(s.slice[i].Time) {
				temp = append(temp, s.slice[i])
				i++
			}
			v := fn(t, temp)
			if v != nil {
				got = append(got, Int64{
					Time:  t,
					Value: *v,
				})
			}
		}
	} else {
		var begin time.Time
		if len(s.slice) > 0 {
			begin = s.slice[0].Time.Truncate(duration)
		}
		if interval.NotBefore != nil {
			begin = (*interval.NotBefore).Truncate(duration)
			if begin.Before(*interval.NotBefore) {
				begin = begin.Add(duration)
			}
		}
		for t, i := begin, 0; i < s.Len() && interval.Contain(t); t = t.Add(duration) {
			temp = temp[:0]
			itv := BeginAt(t).EndAt(t.Add(duration))
			for i < s.slice.Len() && itv.Contain(s.slice[i].Time) {
				temp = append(temp, s.slice[i])
				i++
			}
			v := fn(t, temp)
			if v != nil {
				got = append(got, Int64{
					Time:  t,
					Value: *v,
				})
			}
		}
	}

	s.slice = got
	s.resetIndex()
	return nil
}

// Trim remove the elements which make fn return true
func (s *Int64Seq) Trim(fn func(i int, v Int64) bool) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	updated := false
	slice := make(Int64s, 0)
	for i, v := range s.slice {
		if fn(i, v) {
			updated = true
		} else {
			slice = append(slice, v)
		}
	}

	if updated {
		s.slice = slice
		s.resetIndex()
	}
	return nil
}

// Clone return a new *Int64Seq with copied slice inside
func (s *Int64Seq) Clone() *Int64Seq {
	if s == nil {
		return nil
	}
	return newInt64Seq(s.slice)
}
