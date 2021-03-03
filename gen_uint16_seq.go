// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"sort"
	"sync"
	"time"
)

// Uint16 is a time point with uint16 value inside
type Uint16 struct {
	Time  time.Time
	Value uint16
}

// IsZero returns if time and value are both zero
func (v Uint16) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

// Equal returns if time and value are both equal
func (v Uint16) Equal(n Uint16) bool {
	return v.Value == n.Value && v.Time.Equal(n.Time)
}

// Uint16s is a alias of Uint16 slice
type Uint16s []Uint16

// Len implements Interface.Len()
func (s Uint16s) Len() int {
	return len(s)
}

// Swap implements Interface.Swap()
func (s Uint16s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Interface.Time()
func (s Uint16s) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Interface.Slice()
func (s Uint16s) Slice(i, j int) Interface {
	return s[i:j]
}

// Uint16Seq is a wrapper with useful methods of Uint16 slice
type Uint16Seq struct {
	slice Uint16s

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[uint16][]int
	valueOrder []int
}

// NewUint16Seq returns *Uint16Seq with copied slice inside
func NewUint16Seq(slice Uint16s) *Uint16Seq {
	temp := make(Uint16s, len(slice))
	copy(temp, slice)
	return WrapUint16Seq(temp)
}

// WrapUint16Seq returns *Uint16Seq with origin slice inside
func WrapUint16Seq(slice Uint16s) *Uint16Seq {
	if !IsSorted(slice) {
		Sort(slice)
	}
	return newUint16Seq(slice)
}

func newUint16Seq(slice Uint16s) *Uint16Seq {
	ret := &Uint16Seq{
		slice: slice,
	}
	return ret
}

func (s *Uint16Seq) getSlice() Uint16s {
	if s == nil {
		return nil
	}
	return s.slice
}

func (s *Uint16Seq) getTimeIndex() map[timeKey][]int {
	if s == nil {
		return nil
	}
	return s.timeIndex
}

func (s *Uint16Seq) getValueIndex() map[uint16][]int {
	if s == nil {
		return nil
	}
	return s.valueIndex
}

func (s *Uint16Seq) buildIndex() {
	if s == nil {
		return
	}
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[uint16][]int, len(s.slice))
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

// Uint16s returns a replica of inside slice
func (s *Uint16Seq) Uint16s() Uint16s {
	sslice := s.getSlice()

	slice := make(Uint16s, len(sslice))
	copy(slice, sslice)
	return slice
}

// Len returns length of inside slice
func (s *Uint16Seq) Len() int {
	sslice := s.getSlice()
	return len(sslice)
}

// Index returns element of inside slice, returns zero if index is out of range
func (s *Uint16Seq) Index(i int) Uint16 {
	sslice := s.getSlice()
	if i < 0 || i >= len(sslice) {
		return Uint16{}
	}
	return sslice[i]
}

// Time returns the first element with time t, returns zero if not found
func (s *Uint16Seq) Time(t time.Time) Uint16 {
	got := s.MTime(t)
	if len(got) == 0 {
		return Uint16{}
	}
	return got[0]
}

// MTime returns all elements with time t, returns nil if not found
func (s *Uint16Seq) MTime(t time.Time) Uint16s {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getTimeIndex()[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint16s, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Value returns the first element with value v, returns zero if not found
func (s *Uint16Seq) Value(v uint16) Uint16 {
	got := s.MValue(v)
	if len(got) == 0 {
		return Uint16{}
	}
	return got[0]
}

// MValue returns all elements with value v, returns nil if not found
func (s *Uint16Seq) MValue(v uint16) Uint16s {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getValueIndex()[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint16s, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Traverse call fn for every element one by one, break if fn returns true
func (s *Uint16Seq) Traverse(fn func(i int, v Uint16) (stop bool)) {
	sslice := s.getSlice()
	for i, v := range sslice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

// Sum returns sum of all values
func (s *Uint16Seq) Sum() uint16 {
	var ret uint16
	sslice := s.getSlice()
	for _, v := range sslice {
		ret += v.Value
	}
	return ret
}

// Max returns the element with max value, returns zero if empty
func (s *Uint16Seq) Max() Uint16 {
	var max Uint16
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
func (s *Uint16Seq) Min() Uint16 {
	var min Uint16
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
func (s *Uint16Seq) First() Uint16 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint16{}
	}
	return sslice[0]
}

// Last returns the last element, returns zero if empty
func (s *Uint16Seq) Last() Uint16 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint16{}
	}
	return sslice[len(sslice)-1]
}

// Percentile returns the element matched with percentile pct, returns zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Uint16Seq) Percentile(pct float64) Uint16 {
	s.buildIndex()
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint16{}
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

// Range returns a sub *Uint16Seq with specified interval
func (s *Uint16Seq) Range(interval Interval) *Uint16Seq {
	sslice := s.getSlice()
	slice := Range(sslice, interval).(Uint16s)
	return newUint16Seq(slice)
}

// Slice returns a sub *Uint16Seq with specified index,
// (1, 2) means [1:2], (-1, 2) means [:2], (-1, -1) means [:]
func (s *Uint16Seq) Slice(i, j int) *Uint16Seq {
	if i < 0 && j < 0 {
		return s
	}
	sslice := s.getSlice()
	if i < 0 {
		return newUint16Seq(sslice[:j])
	}
	if j < 0 {
		return newUint16Seq(sslice[i:])
	}
	return newUint16Seq(sslice[i:j])
}

// Trim returns a *Uint16Seq without elements which make fn returns true
func (s *Uint16Seq) Trim(fn func(i int, v Uint16) bool) *Uint16Seq {
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

	slice := make(Uint16s, 0, len(sslice)-len(removeM))
	for i, v := range sslice {
		if _, ok := removeM[i]; ok {
			continue
		}
		slice = append(slice, v)
	}

	return newUint16Seq(slice)
}

// Merge returns a new *Uint16Seq with merged data according to the specified rule
func (s *Uint16Seq) Merge(fn func(t time.Time, v1, v2 *uint16) *uint16, seq *Uint16Seq) *Uint16Seq {
	if fn == nil {
		return s
	}

	var ret Uint16s

	slice1 := s.getSlice()
	var slice2 Uint16s
	if seq != nil {
		slice2 = seq.getSlice()
	}

	for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
		var (
			t time.Time
			v *uint16
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
			ret = append(ret, Uint16{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newUint16Seq(ret)
}

// Aggregate returns a aggregated *Uint16Seq according to the specified rule
func (s *Uint16Seq) Aggregate(fn func(t time.Time, slice Uint16s) *uint16, duration time.Duration, interval Interval) *Uint16Seq {
	if fn == nil {
		return s
	}

	ret := Uint16s{}
	temp := Uint16s{}

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
				ret = append(ret, Uint16{
					Time:  t,
					Value: *v,
				})
			}
		}
		return newUint16Seq(ret)
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
			ret = append(ret, Uint16{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newUint16Seq(ret)
}
