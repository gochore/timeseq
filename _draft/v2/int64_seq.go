package timeseq

import (
	"errors"
	"sort"
	"time"
)

type Int64 struct {
	Time  time.Time
	Value int64
}

type Int64s []Int64

func (s Int64s) Len() int {
	return len(s)
}

func (s Int64s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Int64s) Time(i int) time.Time {
	return s[i].Time
}

func (s Int64s) Slice(i, j int) Interface {
	return s[i:j]
}

type Int64Seq struct {
	slice      Int64s
	timeIndex  map[timeKey][]int
	valueIndex map[int64][]int
	valueSlice []int
}

func NewInt64Seq(slice Int64s) *Int64Seq {
	temp := make(Int64s, len(slice))
	copy(temp, slice)
	slice = temp

	Sort(slice)
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].Time.Before(slice[j].Time)
	})
	return newInt64Seq(slice)
}

func newInt64Seq(slice Int64s) *Int64Seq {
	ret := &Int64Seq{
		slice: slice,
	}
	ret.buildIndex()
	return ret
}

func (s *Int64Seq) buildIndex() {
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
}

func (s *Int64Seq) Int64s() Int64s {
	slice := make(Int64s, len(s.slice))
	copy(slice, s.slice)
	return slice
}

func (s *Int64Seq) Index(i int) Int64 {
	if i < 0 || i >= len(s.slice) {
		return Int64{}
	}
	return s.slice[i]
}

func (s *Int64Seq) Time(t time.Time) Int64s {
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

func (s *Int64Seq) Value(v int64) Int64s {
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

func (s *Int64Seq) Visit(fn func(i int, v Int64) (stop bool)) {
	for i, v := range s.slice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

func (s *Int64Seq) Sum() int64 {
	var ret int64
	for _, v := range s.slice {
		ret += v.Value
	}
	return ret
}

func (s *Int64Seq) Count() int {
	return len(s.slice)
}

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

func (s *Int64Seq) First() Int64 {
	if len(s.slice) == 0 {
		return Int64{}
	}
	return s.slice[0]
}

func (s *Int64Seq) Last() Int64 {
	if len(s.slice) == 0 {
		return Int64{}
	}
	return s.slice[len(s.slice)-1]
}

func (s *Int64Seq) Percentile(pct float64) Int64 {
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

func (s *Int64Seq) Range(interval Interval) *Int64Seq {
	slice := Range(s.slice, interval).(Int64s)
	return newInt64Seq(slice)
}

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
	s.buildIndex()
	return nil
}

func (s *Int64Seq) Aggregate(fn func(t time.Time, slice Int64s) *int64, duration time.Duration, begin, end *time.Time) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	if duration <= 0 {
		return errors.New("invalid duration")
	}

	var bg, ed time.Time
	if len(s.slice) > 0 {
		bg = s.slice[0].Time.Truncate(duration)
		ed = s.slice[len(s.slice)-1].Time
	}
	if begin != nil {
		bg = (*begin).Truncate(duration)
	}
	if end != nil {
		ed = *end
	}

	got := Int64s{}
	slice := Int64s{}
	i := 0
	for t := bg; t.Before(ed); t = t.Add(duration) {
		slice = slice[:0]
		for i < s.slice.Len() &&
			!s.slice[i].Time.After(t) &&
			s.slice[i].Time.Before(t.Add(duration)) {
			slice = append(slice, s.slice[i])
			i++
		}
		v := fn(t, slice)
		if v != nil {
			got = append(got, Int64{
				Time:  t,
				Value: *v,
			})
		}
	}

	s.slice = got
	s.buildIndex()
	return nil
}

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
		s.buildIndex()
	}
	return nil
}
