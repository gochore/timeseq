// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"sort"
	"sync"
	"time"
)

// Uint is a time point with uint value inside
type Uint struct {
	Time  time.Time
	Value uint
}

// IsZero returns if time and value are both zero
func (v Uint) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

// Equal returns if time and value are both equal
func (v Uint) Equal(n Uint) bool {
	return v.Value == n.Value && v.Time.Equal(n.Time)
}

// Uints is a alias of Uint slice
type Uints []Uint

// Len implements Interface.Len()
func (s Uints) Len() int {
	return len(s)
}

// Swap implements Interface.Swap()
func (s Uints) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Interface.Time()
func (s Uints) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Interface.Slice()
func (s Uints) Slice(i, j int) Interface {
	return s[i:j]
}

// UintSeq is a wrapper with useful methods of Uint slice
type UintSeq struct {
	slice Uints

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[uint][]int
	valueOrder []int
}

// NewUintSeq returns *UintSeq with copied slice inside
func NewUintSeq(slice Uints) *UintSeq {
	temp := make(Uints, len(slice))
	copy(temp, slice)
	return WrapUintSeq(temp)
}

// WrapUintSeq returns *UintSeq with origin slice inside
func WrapUintSeq(slice Uints) *UintSeq {
	if !IsSorted(slice) {
		Sort(slice)
	}
	return newUintSeq(slice)
}

func newUintSeq(slice Uints) *UintSeq {
	ret := &UintSeq{
		slice: slice,
	}
	return ret
}

func (s *UintSeq) getSlice() Uints {
	if s == nil {
		return nil
	}
	return s.slice
}

func (s *UintSeq) getTimeIndex() map[timeKey][]int {
	if s == nil {
		return nil
	}
	return s.timeIndex
}

func (s *UintSeq) getValueIndex() map[uint][]int {
	if s == nil {
		return nil
	}
	return s.valueIndex
}

func (s *UintSeq) buildIndex() {
	if s == nil {
		return
	}
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[uint][]int, len(s.slice))
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

// Uints returns a replica of inside slice
func (s *UintSeq) Uints() Uints {
	sslice := s.getSlice()

	slice := make(Uints, len(sslice))
	copy(slice, sslice)
	return slice
}

// Len returns length of inside slice
func (s *UintSeq) Len() int {
	sslice := s.getSlice()
	return len(sslice)
}

// Index returns element of inside slice, returns zero if index is out of range
func (s *UintSeq) Index(i int) Uint {
	sslice := s.getSlice()
	if i < 0 || i >= len(sslice) {
		return Uint{}
	}
	return sslice[i]
}

// Time returns the first element with time t, returns zero if not found
func (s *UintSeq) Time(t time.Time) Uint {
	got := s.MTime(t)
	if len(got) == 0 {
		return Uint{}
	}
	return got[0]
}

// MTime returns all elements with time t, returns nil if not found
func (s *UintSeq) MTime(t time.Time) Uints {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getTimeIndex()[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uints, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Value returns the first element with value v, returns zero if not found
func (s *UintSeq) Value(v uint) Uint {
	got := s.MValue(v)
	if len(got) == 0 {
		return Uint{}
	}
	return got[0]
}

// MValue returns all elements with value v, returns nil if not found
func (s *UintSeq) MValue(v uint) Uints {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getValueIndex()[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uints, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Traverse call fn for every element one by one, break if fn returns true
func (s *UintSeq) Traverse(fn func(i int, v Uint) (stop bool)) {
	sslice := s.getSlice()
	for i, v := range sslice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

// Sum returns sum of all values
func (s *UintSeq) Sum() uint {
	var ret uint
	sslice := s.getSlice()
	for _, v := range sslice {
		ret += v.Value
	}
	return ret
}

// Max returns the element with max value, returns zero if empty
func (s *UintSeq) Max() Uint {
	var max Uint
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
func (s *UintSeq) Min() Uint {
	var min Uint
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
func (s *UintSeq) First() Uint {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint{}
	}
	return sslice[0]
}

// Last returns the last element, returns zero if empty
func (s *UintSeq) Last() Uint {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint{}
	}
	return sslice[len(sslice)-1]
}

// Percentile returns the element matched with percentile pct, returns zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *UintSeq) Percentile(pct float64) Uint {
	s.buildIndex()
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint{}
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

// Truncate returns a sub *UintSeq with specified interval
func (s *UintSeq) Truncate(interval Interval) *UintSeq {
	sslice := s.getSlice()
	slice := Truncate(sslice, interval).(Uints)
	return newUintSeq(slice)
}

// Slice returns a sub *UintSeq with specified index,
// (1, 2) means [1:2], (-1, 2) means [:2], (-1, -1) means [:]
func (s *UintSeq) Slice(i, j int) *UintSeq {
	if i < 0 && j < 0 {
		return s
	}
	sslice := s.getSlice()
	if i < 0 {
		return newUintSeq(sslice[:j])
	}
	if j < 0 {
		return newUintSeq(sslice[i:])
	}
	return newUintSeq(sslice[i:j])
}

// Trim returns a *UintSeq without elements which make fn returns true
func (s *UintSeq) Trim(fn func(i int, v Uint) bool) *UintSeq {
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

	slice := make(Uints, 0, len(sslice)-len(removeM))
	for i, v := range sslice {
		if _, ok := removeM[i]; ok {
			continue
		}
		slice = append(slice, v)
	}

	return newUintSeq(slice)
}

// Merge returns a new *UintSeq with merged data according to the specified rule
func (s *UintSeq) Merge(fn func(t time.Time, v1, v2 *uint) *uint, seq *UintSeq) *UintSeq {
	if fn == nil {
		return s
	}

	var ret Uints

	slice1 := s.getSlice()
	var slice2 Uints
	if seq != nil {
		slice2 = seq.getSlice()
	}

	for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
		var (
			t time.Time
			v *uint
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
			ret = append(ret, Uint{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newUintSeq(ret)
}

// Aggregate returns a aggregated *UintSeq according to the specified rule
func (s *UintSeq) Aggregate(fn func(t time.Time, slice Uints) *uint, duration time.Duration, interval Interval) *UintSeq {
	if fn == nil {
		return s
	}

	ret := Uints{}
	temp := Uints{}

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
				ret = append(ret, Uint{
					Time:  t,
					Value: *v,
				})
			}
		}
		return newUintSeq(ret)
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
			ret = append(ret, Uint{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newUintSeq(ret)
}
