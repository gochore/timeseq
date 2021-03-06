// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"sort"
	"sync"
	"time"
)

// Int16 is a time point with int16 value inside
type Int16 struct {
	Time  time.Time
	Value int16
}

// IsZero returns if time and value are both zero
func (v Int16) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

// Equal returns if time and value are both equal
func (v Int16) Equal(n Int16) bool {
	return v.Value == n.Value && v.Time.Equal(n.Time)
}

// Int16s is a alias of Int16 slice
type Int16s []Int16

// Len implements Interface.Len()
func (s Int16s) Len() int {
	return len(s)
}

// Swap implements Interface.Swap()
func (s Int16s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Interface.Time()
func (s Int16s) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Interface.Slice()
func (s Int16s) Slice(i, j int) Interface {
	return s[i:j]
}

// Int16Seq is a wrapper with useful methods of Int16 slice
type Int16Seq struct {
	slice Int16s

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[int16][]int
	valueOrder []int
}

// NewInt16Seq returns *Int16Seq with copied slice inside
func NewInt16Seq(slice Int16s) *Int16Seq {
	temp := make(Int16s, len(slice))
	copy(temp, slice)
	return WrapInt16Seq(temp)
}

// WrapInt16Seq returns *Int16Seq with origin slice inside
func WrapInt16Seq(slice Int16s) *Int16Seq {
	if !IsSorted(slice) {
		Sort(slice)
	}
	return newInt16Seq(slice)
}

func newInt16Seq(slice Int16s) *Int16Seq {
	ret := &Int16Seq{
		slice: slice,
	}
	return ret
}

func (s *Int16Seq) getSlice() Int16s {
	if s == nil {
		return nil
	}
	return s.slice
}

func (s *Int16Seq) getTimeIndex() map[timeKey][]int {
	if s == nil {
		return nil
	}
	return s.timeIndex
}

func (s *Int16Seq) getValueIndex() map[int16][]int {
	if s == nil {
		return nil
	}
	return s.valueIndex
}

func (s *Int16Seq) buildIndex() {
	if s == nil {
		return
	}
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[int16][]int, len(s.slice))
		valueSlice := s.valueOrder[:0]
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
		s.valueOrder = valueSlice
	})
}

// Int16s returns a replica of inside slice
func (s *Int16Seq) Int16s() Int16s {
	sslice := s.getSlice()

	slice := make(Int16s, len(sslice))
	copy(slice, sslice)
	return slice
}

// Len returns length of inside slice
func (s *Int16Seq) Len() int {
	sslice := s.getSlice()
	return len(sslice)
}

// Index returns element of inside slice, returns zero if index is out of range
func (s *Int16Seq) Index(i int) Int16 {
	sslice := s.getSlice()
	if i < 0 || i >= len(sslice) {
		return Int16{}
	}
	return sslice[i]
}

// Time returns the first element with time t, returns zero if not found
func (s *Int16Seq) Time(t time.Time) Int16 {
	got := s.MTime(t)
	if len(got) == 0 {
		return Int16{}
	}
	return got[0]
}

// MTime returns all elements with time t, returns nil if not found
func (s *Int16Seq) MTime(t time.Time) Int16s {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getTimeIndex()[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Int16s, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Value returns the first element with value v, returns zero if not found
func (s *Int16Seq) Value(v int16) Int16 {
	got := s.MValue(v)
	if len(got) == 0 {
		return Int16{}
	}
	return got[0]
}

// MValue returns all elements with value v, returns nil if not found
func (s *Int16Seq) MValue(v int16) Int16s {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getValueIndex()[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Int16s, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Traverse call fn for every element one by one, break if fn returns true
func (s *Int16Seq) Traverse(fn func(i int, v Int16) (stop bool)) {
	sslice := s.getSlice()
	for i, v := range sslice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

// Sum returns sum of all values
func (s *Int16Seq) Sum() int16 {
	var ret int16
	sslice := s.getSlice()
	for _, v := range sslice {
		ret += v.Value
	}
	return ret
}

// Max returns the element with max value, returns zero if empty
func (s *Int16Seq) Max() Int16 {
	var max Int16
	found := false
	sslice := s.getSlice()
	for _, v := range sslice {
		if !found {
			max = v
			found = true
		} else if v.Value > max.Value {
			max = v
		}
	}
	return max
}

// Min returns the element with min value, returns zero if empty
func (s *Int16Seq) Min() Int16 {
	var min Int16
	found := false
	sslice := s.getSlice()
	for _, v := range sslice {
		if !found {
			min = v
			found = true
		} else if v.Value < min.Value {
			min = v
		}
	}
	return min
}

// First returns the first element, returns zero if empty
func (s *Int16Seq) First() Int16 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Int16{}
	}
	return sslice[0]
}

// Last returns the last element, returns zero if empty
func (s *Int16Seq) Last() Int16 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Int16{}
	}
	return sslice[len(sslice)-1]
}

// Percentile returns the element matched with percentile pct, returns zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Int16Seq) Percentile(pct float64) Int16 {
	s.buildIndex()
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Int16{}
	}
	if pct > 1 {
		pct = 1
	}
	if pct < 0 {
		pct = 0
	}
	i := int(float64(len(sslice))*pct - 1)
	if i < 0 {
		i = 0
	}
	return sslice[s.valueOrder[i]]
}

// Range returns a sub *Int16Seq with specified interval
func (s *Int16Seq) Range(interval Interval) *Int16Seq {
	sslice := s.getSlice()
	slice := Range(sslice, interval).(Int16s)
	return newInt16Seq(slice)
}

// Slice returns a sub *Int16Seq with specified index,
// (1, 2) means [1:2], (-1, 2) means [:2], (-1, -1) means [:]
func (s *Int16Seq) Slice(i, j int) *Int16Seq {
	if i < 0 && j < 0 {
		return s
	}
	sslice := s.getSlice()
	if i < 0 {
		return newInt16Seq(sslice[:j])
	}
	if j < 0 {
		return newInt16Seq(sslice[i:])
	}
	return newInt16Seq(sslice[i:j])
}

// Trim returns a *Int16Seq without elements which make fn returns true
func (s *Int16Seq) Trim(fn func(i int, v Int16) bool) *Int16Seq {
	sslice := s.getSlice()
	if fn == nil || len(sslice) == 0 {
		return s
	}

	removeM := map[int]struct{}{}
	for i, v := range sslice {
		if fn(i, v) {
			removeM[i] = struct{}{}
		}
	}
	if len(removeM) == 0 {
		return s
	}

	slice := make(Int16s, 0, len(sslice)-len(removeM))
	for i, v := range sslice {
		if _, ok := removeM[i]; ok {
			continue
		}
		slice = append(slice, v)
	}

	return newInt16Seq(slice)
}

// Merge returns a new *Int16Seq with merged data according to the specified rule
func (s *Int16Seq) Merge(fn func(t time.Time, v1, v2 *int16) *int16, seq *Int16Seq) *Int16Seq {
	if fn == nil {
		return s
	}

	var ret Int16s

	slice1 := s.getSlice()
	var slice2 Int16s
	if seq != nil {
		slice2 = seq.getSlice()
	}

	for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
		var (
			t time.Time
			v *int16
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
			ret = append(ret, Int16{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newInt16Seq(ret)
}

// Aggregate returns a aggregated *Int16Seq according to the specified rule
func (s *Int16Seq) Aggregate(fn func(t time.Time, slice Int16s) *int16, duration time.Duration, interval Interval) *Int16Seq {
	if fn == nil {
		return s
	}

	ret := Int16s{}
	temp := Int16s{}

	sslice := s.getSlice()
	if duration <= 0 {
		for i := 0; i < s.Len(); {
			t := sslice[i].Time
			if !interval.Contain(t) {
				i++
				continue
			}
			temp = temp[:0]
			for i < sslice.Len() && t.Equal(sslice[i].Time) {
				temp = append(temp, sslice[i])
				i++
			}
			v := fn(t, temp)
			if v != nil {
				ret = append(ret, Int16{
					Time:  t,
					Value: *v,
				})
			}
		}
		return newInt16Seq(ret)
	}

	if len(sslice) == 0 && interval.Duration() < 0 {
		return s
	}

	var begin time.Time
	if len(sslice) > 0 {
		begin = sslice[0].Time.Truncate(duration)
	}
	if interval.NotBefore != nil {
		begin = (*interval.NotBefore).Truncate(duration)
		if begin.Before(*interval.NotBefore) {
			begin = begin.Add(duration)
		}
	}

	var end time.Time
	if len(sslice) > 0 {
		end = sslice[len(sslice)-1].Time.Truncate(duration)
	}
	if interval.NotAfter != nil {
		end = (*interval.NotAfter).Truncate(duration)
	}

	for t, i := begin, 0; !t.After(end); t = t.Add(duration) {
		temp = temp[:0]
		itv := BeginAt(t).EndAt(t.Add(duration))
		for i < len(sslice) {
			if sslice[i].Time.After(*itv.NotAfter) {
				break
			}
			if itv.Contain(sslice[i].Time) {
				temp = append(temp, sslice[i])
			}
			i++
		}
		v := fn(t, temp)
		if v != nil {
			ret = append(ret, Int16{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newInt16Seq(ret)
}
