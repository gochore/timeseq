// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"errors"
	"sort"
	"sync"
	"time"
)

type Uint64 struct {
	Time  time.Time
	Value uint64
}

func (v Uint64) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

func (v Uint64) Equal(n Uint64) bool {
	return v.Value == n.Value && v.Time.Equal(n.Time)
}

type Uint64s []Uint64

func (s Uint64s) Len() int {
	return len(s)
}

func (s Uint64s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Uint64s) Time(i int) time.Time {
	return s[i].Time
}

func (s Uint64s) Slice(i, j int) Interface {
	return s[i:j]
}

type Uint64Seq struct {
	slice Uint64s

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[uint64][]int
	valueSlice []int
}

func NewUint64Seq(slice Uint64s) *Uint64Seq {
	temp := make(Uint64s, len(slice))
	copy(temp, slice)
	return WrapUint64Seq(temp)
}

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

func (s *Uint64Seq) buildIndex() {
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[uint64][]int, len(s.slice))
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

func (s *Uint64Seq) resetIndex() {
	s.indexOnce = sync.Once{}
}

func (s *Uint64Seq) Uint64s() Uint64s {
	slice := make(Uint64s, len(s.slice))
	copy(slice, s.slice)
	return slice
}

func (s *Uint64Seq) Len() int {
	return len(s.slice)
}

func (s *Uint64Seq) Index(i int) Uint64 {
	if i < 0 || i >= len(s.slice) {
		return Uint64{}
	}
	return s.slice[i]
}

func (s *Uint64Seq) Time(t time.Time) Uint64 {
	got := s.MTime(t)
	if len(got) == 0 {
		return Uint64{}
	}
	return got[0]
}

func (s *Uint64Seq) MTime(t time.Time) Uint64s {
	s.buildIndex()
	index := s.timeIndex[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint64s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

func (s *Uint64Seq) Value(v uint64) Uint64 {
	got := s.MValue(v)
	if len(got) == 0 {
		return Uint64{}
	}
	return got[0]
}

func (s *Uint64Seq) MValue(v uint64) Uint64s {
	s.buildIndex()
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Uint64s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

func (s *Uint64Seq) Visit(fn func(i int, v Uint64) (stop bool)) {
	for i, v := range s.slice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

func (s *Uint64Seq) Sum() uint64 {
	var ret uint64
	for _, v := range s.slice {
		ret += v.Value
	}
	return ret
}

func (s *Uint64Seq) Count() int {
	return s.Len()
}

func (s *Uint64Seq) Max() Uint64 {
	var max Uint64
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

func (s *Uint64Seq) Min() Uint64 {
	var min Uint64
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

func (s *Uint64Seq) First() Uint64 {
	if len(s.slice) == 0 {
		return Uint64{}
	}
	return s.slice[0]
}

func (s *Uint64Seq) Last() Uint64 {
	if len(s.slice) == 0 {
		return Uint64{}
	}
	return s.slice[len(s.slice)-1]
}

func (s *Uint64Seq) Percentile(pct float64) Uint64 {
	s.buildIndex()
	if len(s.slice) == 0 {
		return Uint64{}
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

func (s *Uint64Seq) Range(interval Interval) *Uint64Seq {
	slice := Range(s.slice, interval).(Uint64s)
	return newUint64Seq(slice)
}

func (s *Uint64Seq) Merge(fn func(t time.Time, v1, v2 *uint64) *uint64, slices ...Uint64s) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	if len(slices) == 0 {
		return nil
	}

	slice1 := s.slice
	for _, slice2 := range slices {
		if !IsSorted(slice2) {
			temp := make(Uint64s, len(slice2))
			copy(temp, slice2)
			Sort(temp)
			slice2 = temp
		}
		var got Uint64s
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
				got = append(got, Uint64{
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

func (s *Uint64Seq) Aggregate(fn func(t time.Time, slice Uint64s) *uint64, duration time.Duration, interval Interval) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	got := Uint64s{}
	temp := Uint64s{}

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
				got = append(got, Uint64{
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
				got = append(got, Uint64{
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

func (s *Uint64Seq) Trim(fn func(i int, v Uint64) bool) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	updated := false
	slice := make(Uint64s, 0)
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

func (s *Uint64Seq) Clone() *Uint64Seq {
	if s == nil {
		return nil
	}
	return newUint64Seq(s.slice)
}
