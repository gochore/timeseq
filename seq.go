package timeseq

import (
	"slices"
	"sort"
	"sync"
	"time"
)

type Item[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64] struct {
	Time  time.Time
	Value T
}

func (i Item[T]) IsZero() bool {
	var zero T
	return i.Time.IsZero() && i.Value == zero
}

type Seq[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64] struct {
	slice []Item[T]

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[T][]int
	valueOrder []int
}

// NewSeq returns Seq with copied slice inside
func NewSeq[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](slice []Item[T]) *Seq[T] {
	temp := make([]Item[T], len(slice))
	copy(temp, slice)
	return WrapSeq(temp)
}

// WrapSeq returns Seq with origin slice inside
func WrapSeq[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](slice []Item[T]) *Seq[T] {
	if !slices.IsSortedFunc(slice, compareItems[T]) {
		slices.SortStableFunc(slice, compareItems[T])
	}
	return &Seq[T]{
		slice: slice,
	}
}

func compareItems[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](a, b Item[T]) int {
	if a.Time.Before(b.Time) {
		return -1
	}
	if a.Time.After(b.Time) {
		return 1
	}
	return 0
}

func (s *Seq[T]) buildIndex() {
	if s == nil {
		return
	}
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[T][]int, len(s.slice))
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

// Slice returns a replica of inside slice
func (s *Seq[T]) Slice() []Item[T] {
	ret := make([]Item[T], len(s.slice))
	copy(ret, s.slice)
	return ret
}

// Len returns length of inside slice
func (s *Seq[T]) Len() int {
	return len(s.slice)
}

// Index returns element of inside slice, returns zero if index is out of range
func (s *Seq[T]) Index(i int) Item[T] {
	if i < 0 || i >= len(s.slice) {
		return Item[T]{}
	}
	return s.slice[i]
}

// Time returns the first element with time t, returns zero if not found
func (s *Seq[T]) Time(t time.Time) Item[T] {
	got := s.MTime(t)
	if len(got) == 0 {
		return Item[T]{}
	}
	return got[0]
}

// MTime returns all elements with time t, returns nil if not found
func (s *Seq[T]) MTime(t time.Time) []Item[T] {
	s.buildIndex()
	index := s.timeIndex[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make([]Item[T], len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

// Value returns the first element with value v, returns zero if not found
func (s *Seq[T]) Value(v T) Item[T] {
	got := s.MValue(v)
	if len(got) == 0 {
		return Item[T]{}
	}
	return got[0]
}

// MValue returns all elements with value v, returns nil if not found
func (s *Seq[T]) MValue(v T) []Item[T] {
	s.buildIndex()
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make([]Item[T], len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

// Traverse call fn for every element one by one, break if fn returns false
func (s *Seq[T]) Traverse(fn func(i int, v Item[T]) (stop bool)) {
	for i, v := range s.slice {
		if !fn(i, v) {
			break
		}
	}
}

// Sum returns sum of all values
func (s *Seq[T]) Sum() T {
	var ret T
	for _, v := range s.slice {
		ret += v.Value
	}
	return ret
}

// Max returns the element with max value, returns zero if empty
func (s *Seq[T]) Max() Item[T] {
	var ret Item[T]
	for _, v := range s.slice {
		if ret.IsZero() {
			ret = v
		} else if v.Value > ret.Value {
			ret = v
		}
	}
	return ret
}

// Min returns the element with min value, returns zero if empty
func (s *Float32Seq) Min() Float32 {
	var min Float32
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
func (s *Float32Seq) First() Float32 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Float32{}
	}
	return sslice[0]
}

// Last returns the last element, returns zero if empty
func (s *Float32Seq) Last() Float32 {
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Float32{}
	}
	return sslice[len(sslice)-1]
}

// Percentile returns the element matched with percentile pct, returns zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Float32Seq) Percentile(pct float64) Float32 {
	s.buildIndex()
	sslice := s.getSlice()
	if len(sslice) == 0 {
		return Float32{}
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

// Range returns a sub *Float32Seq with specified interval
func (s *Float32Seq) Range(interval Interval) *Float32Seq {
	sslice := s.getSlice()
	slice := Range(sslice, interval).(Float32s)
	return newFloat32Seq(slice)
}

// Slice returns a sub *Float32Seq with specified index,
// (1, 2) means [1:2], (-1, 2) means [:2], (-1, -1) means [:]
func (s *Float32Seq) Slice(i, j int) *Float32Seq {
	if i < 0 && j < 0 {
		return s
	}
	sslice := s.getSlice()
	if i < 0 {
		return newFloat32Seq(sslice[:j])
	}
	if j < 0 {
		return newFloat32Seq(sslice[i:])
	}
	return newFloat32Seq(sslice[i:j])
}

// Trim returns a *Float32Seq without elements which make fn returns true
func (s *Float32Seq) Trim(fn func(i int, v Float32) bool) *Float32Seq {
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

	slice := make(Float32s, 0, len(sslice)-len(removeM))
	for i, v := range sslice {
		if _, ok := removeM[i]; ok {
			continue
		}
		slice = append(slice, v)
	}

	return newFloat32Seq(slice)
}

// Merge returns a new *Float32Seq with merged data according to the specified rule
func (s *Float32Seq) Merge(fn func(t time.Time, v1, v2 *float32) *float32, seq *Float32Seq) *Float32Seq {
	if fn == nil {
		return s
	}

	var ret Float32s

	slice1 := s.getSlice()
	var slice2 Float32s
	if seq != nil {
		slice2 = seq.getSlice()
	}

	for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
		var (
			t time.Time
			v *float32
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
			ret = append(ret, Float32{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newFloat32Seq(ret)
}

// Aggregate returns a aggregated *Float32Seq according to the specified rule
func (s *Float32Seq) Aggregate(fn func(t time.Time, slice Float32s) *float32, duration time.Duration, interval Interval) *Float32Seq {
	if fn == nil {
		return s
	}

	ret := Float32s{}
	temp := Float32s{}

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
				ret = append(ret, Float32{
					Time:  t,
					Value: *v,
				})
			}
		}
		return newFloat32Seq(ret)
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
			ret = append(ret, Float32{
				Time:  t,
				Value: *v,
			})
		}
	}

	return newFloat32Seq(ret)
}
