// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// Uint8 is a time point with uint8 value inside
type Uint8 struct {
	Time  time.Time
	Value uint8
}

// IsZero return if time and value are both zero
func (v Uint8) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

// Equal return if time and value are both equal
func (v Uint8) Equal(n Uint8) bool {
	return v.Value == n.Value && v.Time.Equal(n.Time)
}

// Uint8s is a alias of Uint8 slice
type Uint8s []Uint8

// Len implements Interface.Len()
func (s Uint8s) Len() int {
	return len(s)
}

// Swap implements Interface.Swap()
func (s Uint8s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Interface.Time()
func (s Uint8s) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Interface.Slice()
func (s Uint8s) Slice(i, j int) Interface {
	return s[i:j]
}

// Uint8Seq is a wrapper with useful methods of Uint8 slice
type Uint8Seq struct {
	slice Uint8s

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[uint8][]int
	valueSlice []int
}

// NewUint8Seq return *Uint8Seq with copied slice inside
func NewUint8Seq(slice Uint8s) *Uint8Seq {
	temp := make(Uint8s, len(slice))
	copy(temp, slice)
	return WrapUint8Seq(temp)
}

// WrapUint8Seq return *Uint8Seq with origin slice inside
func WrapUint8Seq(slice Uint8s) *Uint8Seq {
	if !IsSorted(slice) {
		Sort(slice)
	}
	return newUint8Seq(slice)
}

func newUint8Seq(slice Uint8s) *Uint8Seq {
	ret := &Uint8Seq{
		slice: slice,
	}
	return ret
}

func (s *Uint8Seq) buildIndex() {
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[uint8][]int, len(s.slice))
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

func (s *Uint8Seq) resetIndex() {
	s.indexOnce = sync.Once{}
}

// Uint8s return a replica of inside slice
func (s *Uint8Seq) Uint8s() Uint8s {
	slice := make(Uint8s, len(s.slice))
	copy(slice, s.slice)
	return slice
}

// Len return length of inside slice
func (s *Uint8Seq) Len() int {
	return len(s.slice)
}

// Index return element of inside slice, return zero if index is out of range
func (s *Uint8Seq) Index(i int) Uint8 {
	if i < 0 || i >= len(s.slice) {
		return Uint8{}
	}
	return s.slice[i]
}

// Time return the first element with time t, return zero if not found
func (s *Uint8Seq) Time(t time.Time) Uint8 {
	got := s.MTime(t)
	if len(got) == 0 {
		return Uint8{}
	}
	return got[0]
}

// MTime return all elements with time t, return nil if not found
func (s *Uint8Seq) MTime(t time.Time) Uint8s {
	s.buildIndex()
	index := s.timeIndex[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint8s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

// Value return the first element with value v, return zero if not found
func (s *Uint8Seq) Value(v uint8) Uint8 {
	got := s.MValue(v)
	if len(got) == 0 {
		return Uint8{}
	}
	return got[0]
}

// MValue return all elements with value v, return nil if not found
func (s *Uint8Seq) MValue(v uint8) Uint8s {
	s.buildIndex()
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint8s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

// Traverse call fn for every element one by one, break if fn return true
func (s *Uint8Seq) Traverse(fn func(i int, v Uint8) (stop bool)) {
	for i, v := range s.slice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

// Sum return sum of all values
func (s *Uint8Seq) Sum() uint8 {
	var ret uint8
	for _, v := range s.slice {
		ret += v.Value
	}
	return ret
}

// Max return the element with max value, return zero if empty
func (s *Uint8Seq) Max() Uint8 {
	var max Uint8
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

// Min return the element with min value, return zero if empty
func (s *Uint8Seq) Min() Uint8 {
	var min Uint8
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
func (s *Uint8Seq) First() Uint8 {
	if len(s.slice) == 0 {
		return Uint8{}
	}
	return s.slice[0]
}

// Last return the last element, return zero if empty
func (s *Uint8Seq) Last() Uint8 {
	if len(s.slice) == 0 {
		return Uint8{}
	}
	return s.slice[len(s.slice)-1]
}

// Percentile return the element matched with percentile pct, return zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Uint8Seq) Percentile(pct float64) Uint8 {
	s.buildIndex()
	if len(s.slice) == 0 {
		return Uint8{}
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

// Range return a sub *Uint8Seq with specified interval
func (s *Uint8Seq) Range(interval Interval) *Uint8Seq {
	slice := Range(s.slice, interval).(Uint8s)
	return newUint8Seq(slice)
}

// Merge merge slices to inside slice according to the specified rule
func (s *Uint8Seq) Merge(fn func(t time.Time, v1, v2 *uint8) *uint8, slices ...Uint8s) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	if len(slices) == 0 {
		return nil
	}

	slice1 := s.slice
	for _, slice2 := range slices {
		if !IsSorted(slice2) {
			temp := make(Uint8s, len(slice2))
			copy(temp, slice2)
			Sort(temp)
			slice2 = temp
		}
		var got Uint8s
		for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
			var (
				t time.Time
				v *uint8
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
				got = append(got, Uint8{
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
func (s *Uint8Seq) Aggregate(fn func(t time.Time, slice Uint8s) *uint8, duration time.Duration, interval Interval) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	got := Uint8s{}
	temp := Uint8s{}

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
				got = append(got, Uint8{
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
				got = append(got, Uint8{
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
func (s *Uint8Seq) Trim(fn func(i int, v Uint8) bool) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	updated := false
	slice := make(Uint8s, 0)
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

// Clone return a new *Uint8Seq with copied slice inside
func (s *Uint8Seq) Clone() *Uint8Seq {
	if s == nil {
		return nil
	}
	return newUint8Seq(s.slice)
}
