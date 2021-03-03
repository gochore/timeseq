// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"sort"
	"sync"
	"time"
)

// Uint64 is a time point with uint64 value inside
type Uint64 struct {
	Time  time.Time
	Value uint64
}

// IsZero returns if time and value are both zero
func (v Uint64) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

// Equal returns if time and value are both equal
func (v Uint64) Equal(n Uint64) bool {
	return v.Value == n.Value && v.Time.Equal(n.Time)
}

// Uint64s is a alias of Uint64 slice
type Uint64s []Uint64

// Len implements Interface.Len()
func (s Uint64s) Len() int {
	return len(s)
}

// Swap implements Interface.Swap()
func (s Uint64s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Time implements Interface.Time()
func (s Uint64s) Time(i int) time.Time {
	return s[i].Time
}

// Slice implements Interface.Slice()
func (s Uint64s) Slice(i, j int) Interface {
	return s[i:j]
}

// Uint64Seq is a wrapper with useful methods of Uint64 slice
type Uint64Seq struct {
	slice Uint64s

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[uint64][]int
	valueOrder []int
}

// NewUint64Seq returns *Uint64Seq with copied slice inside
func NewUint64Seq(slice Uint64s) *Uint64Seq {
	temp := make(Uint64s, len(slice))
	copy(temp, slice)
	return WrapUint64Seq(temp)
}

// WrapUint64Seq returns *Uint64Seq with origin slice inside
func WrapUint64Seq(slice Uint64s) *Uint64Seq {
	if !IsSorted(slice) {
		Sort(slice)
	}
	return newUint64Seq(slice)
}

func newUint64Seq(slice Uint64s) *Uint64Seq {
	ret := &Uint64Seq{
		slice: slice,
	}
	return ret
}

func (s *Uint64Seq) getSlice() Uint64s {
	if s == nil {
		return nil
	}
	return s.slice
}

func (s *Uint64Seq) getTimeIndex() map[timeKey][]int {
	if s == nil {
		return nil
	}
	return s.timeIndex
}

func (s *Uint64Seq) getValueIndex() map[uint64][]int {
	if s == nil {
		return nil
	}
	return s.valueIndex
}

func (s *Uint64Seq) buildIndex() {
	if s == nil {
		return
	}
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[uint64][]int, len(s.slice))
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

// Uint64s returns a replica of inside slice
func (s *Uint64Seq) Uint64s() Uint64s {
	sslice := s.getSlice()

	slice := make(Uint64s, len(sslice))
	copy(slice, sslice)
	return slice
}

// Len returns length of inside slice
func (s *Uint64Seq) Len() int {
	sslice := s.getSlice()
	return len(sslice)
}

// Index returns element of inside slice, returns zero if index is out of range
func (s *Uint64Seq) Index(i int) Uint64 {
	sslice := s.getSlice()
	if i < 0 || i >= len(sslice) {
		return Uint64{}
	}
	return sslice[i]
}

// Time returns the first element with time t, returns zero if not found
func (s *Uint64Seq) Time(t time.Time) Uint64 {
	got := s.MTime(t)
	if len(got) == 0 {
		return Uint64{}
	}
	return got[0]
}

// MTime returns all elements with time t, returns nil if not found
func (s *Uint64Seq) MTime(t time.Time) Uint64s {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getTimeIndex()[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint64s, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Value returns the first element with value v, returns zero if not found
func (s *Uint64Seq) Value(v uint64) Uint64 {
	got := s.MValue(v)
	if len(got) == 0 {
		return Uint64{}
	}
	return got[0]
}

// MValue returns all elements with value v, returns nil if not found
func (s *Uint64Seq) MValue(v uint64) Uint64s {
	s.buildIndex()
	sslice := s.getSlice()
	index := s.getValueIndex()[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint64s, len(index))
	for i, v := range index {
		ret[i] = sslice[v]
	}
	return ret
}

// Traverse call fn for every element one by one, break if fn returns true
func (s *Uint64Seq) Traverse(fn func(i int, v Uint64) (stop bool)) {
	sslice := s.getSlice()
	for i, v := range sslice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

// Sum returns sum of all values
func (s *Uint64Seq) Sum() uint64 {
	var ret uint64
	sslice := s.getSlice()
	for _, v := range sslice {
		ret += v.Value
	}
	return ret
}

// Max returns the element with max value, returns zero if empty
func (s *Uint64Seq) Max() Uint64 {
	var max Uint64
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
func (s *Uint64Seq) Min() Uint64 {
	var min Uint64
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
func (s *Uint64Seq) First() Uint64 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint64{}
	}
	return sslice[0]
}

// Last returns the last element, returns zero if empty
func (s *Uint64Seq) Last() Uint64 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint64{}
	}
	return sslice[len(sslice)-1]
}

// Percentile returns the element matched with percentile pct, returns zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Uint64Seq) Percentile(pct float64) Uint64 {
	s.buildIndex()
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Uint64{}
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

// Range returns a sub *Uint64Seq with specified interval
func (s *Uint64Seq) Range(interval Interval) *Uint64Seq {
	sslice := s.getSlice()
	slice := Range(sslice, interval).(Uint64s)
	return newUint64Seq(slice)
}

// Slice returns a sub *Uint64Seq with specified index
func (s *Uint64Seq) Slice(i, j int) *Uint64Seq {
	sslice := s.getSlice()
	slice := sslice[i:j]
	return newUint64Seq(slice)
}

// Trim returns a *Uint64Seq without elements which make fn returns true
func (s *Uint64Seq) Trim(fn func(i int, v Uint64) bool) *Uint64Seq {
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

	slice := make(Uint64s, 0, len(sslice)-len(removeM))
	for i, v := range sslice {
		if _, ok := removeM[i]; ok {
			continue
		}
		slice = append(slice, v)
	}

	return newUint64Seq(slice)
}

// Merge returns a new *Uint64Seq with merged data according to the specified rule
func (s *Uint64Seq) Merge(fn func(t time.Time, v1, v2 *uint64) *uint64, seq *Uint64Seq) *Uint64Seq {
	if fn == nil {
		return s
	}

	var ret Uint64s

	slice1 := s.getSlice()
	var slice2 Uint64s
	if seq != nil {
		slice2 = seq.getSlice()
	}

	for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
		var (
			t time.Time
			v *uint64
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
			ret = append(ret, Uint64{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newUint64Seq(ret)
}

// Aggregate returns a aggregated *Uint64Seq according to the specified rule
func (s *Uint64Seq) Aggregate(fn func(t time.Time, slice Uint64s) *uint64, duration time.Duration, interval Interval) *Uint64Seq {
	if fn == nil {
		return s
	}

	ret := Uint64s{}
	temp := Uint64s{}

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
				ret = append(ret, Uint64{
					Time:  t,
					Value: *v,
				})
			}
		}
		return newUint64Seq(ret)
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
			ret = append(ret, Uint64{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newUint64Seq(ret)
}
