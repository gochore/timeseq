// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"errors"
	"sort"
	"sync"
	"time"
)

type Float64 struct {
	Time  time.Time
	Value float64
}

func (v Float64) IsZero() bool {
	return v.Value == 0 && v.Time.IsZero()
}

type Float64s []Float64

func (s Float64s) Len() int {
	return len(s)
}

func (s Float64s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Float64s) Time(i int) time.Time {
	return s[i].Time
}

func (s Float64s) Slice(i, j int) Interface {
	return s[i:j]
}

type Float64Seq struct {
	slice Float64s

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[float64][]int
	valueSlice []int
}

func NewFloat64Seq(slice Float64s) *Float64Seq {
	temp := make(Float64s, len(slice))
	copy(temp, slice)
	return WrapFloat64Seq(temp)
}

func WrapFloat64Seq(slice Float64s) *Float64Seq {
	if !IsSorted(slice) {
		Sort(slice)
	}
	return newFloat64Seq(slice)
}

func newFloat64Seq(slice Float64s) *Float64Seq {
	ret := &Float64Seq{
		slice: slice,
	}
	return ret
}

func (s *Float64Seq) buildIndex() {
	s.indexOnce.Do(func() {
		timeIndex := make(map[timeKey][]int, len(s.slice))
		valueIndex := make(map[float64][]int, len(s.slice))
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

func (s *Float64Seq) resetIndex() {
	s.indexOnce = sync.Once{}
}

func (s *Float64Seq) Float64s() Float64s {
	slice := make(Float64s, len(s.slice))
	copy(slice, s.slice)
	return slice
}

func (s *Float64Seq) Len() int {
	return len(s.slice)
}

func (s *Float64Seq) Index(i int) Float64 {
	if i < 0 || i >= len(s.slice) {
		return Float64{}
	}
	return s.slice[i]
}

func (s *Float64Seq) Time(t time.Time) Float64s {
	s.buildIndex()
	index := s.timeIndex[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Float64s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

func (s *Float64Seq) Value(v float64) Float64s {
	s.buildIndex()
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Float64s, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

func (s *Float64Seq) Visit(fn func(i int, v Float64) (stop bool)) {
	for i, v := range s.slice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

func (s *Float64Seq) Sum() float64 {
	var ret float64
	for _, v := range s.slice {
		ret += v.Value
	}
	return ret
}

func (s *Float64Seq) Count() int {
	return s.Len()
}

func (s *Float64Seq) Max() Float64 {
	var max Float64
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

func (s *Float64Seq) Min() Float64 {
	var min Float64
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

func (s *Float64Seq) First() Float64 {
	if len(s.slice) == 0 {
		return Float64{}
	}
	return s.slice[0]
}

func (s *Float64Seq) Last() Float64 {
	if len(s.slice) == 0 {
		return Float64{}
	}
	return s.slice[len(s.slice)-1]
}

func (s *Float64Seq) Percentile(pct float64) Float64 {
	s.buildIndex()
	if len(s.slice) == 0 {
		return Float64{}
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

func (s *Float64Seq) Range(interval Interval) *Float64Seq {
	slice := Range(s.slice, interval).(Float64s)
	return newFloat64Seq(slice)
}

func (s *Float64Seq) Merge(fn func(t time.Time, v1, v2 *float64) *float64, slices ...Float64s) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	if len(slices) == 0 {
		return nil
	}

	slice1 := s.slice
	for _, slice2 := range slices {
		if !IsSorted(slice2) {
			temp := make(Float64s, len(slice2))
			copy(temp, slice2)
			Sort(temp)
			slice2 = temp
		}
		var got Float64s
		for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
			var (
				t time.Time
				v *float64
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
				got = append(got, Float64{
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

func (s *Float64Seq) Aggregate(fn func(t time.Time, slice Float64s) *float64, duration time.Duration, begin, end *time.Time) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	got := Float64s{}
	temp := Float64s{}

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
				got = append(got, Float64{
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
				got = append(got, Float64{
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

func (s *Float64Seq) Trim(fn func(i int, v Float64) bool) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	updated := false
	slice := make(Float64s, 0)
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
