package timeseq

import (
	"errors"
	"sort"
	"time"
)

type Int struct {
	Time  time.Time
	Value int
}

type Ints []Int

func (s Ints) Len() int {
	return len(s)
}

func (s Ints) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Ints) Time(i int) time.Time {
	return s[i].Time
}

func (s Ints) Slice(i, j int) Slice {
	return s[i:j]
}

type IntSeq struct {
	slice      Ints
	timeIndex  map[timeKey][]int
	valueIndex map[int][]int
	valueSlice []int
}

func NewIntSeq(ints Ints) *IntSeq {
	slice := make(Ints, len(ints))
	copy(slice, ints)
	Sort(slice)
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].Time.Before(slice[j].Time)
	})
	return newIntSeq(slice)
}

func newIntSeq(slice Ints) *IntSeq {
	ret := &IntSeq{
		slice: slice,
	}
	ret.buildIndex()
	return ret
}

func (s *IntSeq) buildIndex() {
	timeIndex := make(map[timeKey][]int, len(s.slice))
	valueIndex := make(map[int][]int, len(s.slice))
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

func (s *IntSeq) Ints() Ints {
	slice := make(Ints, len(s.slice))
	copy(slice, s.slice)
	return slice
}

func (s *IntSeq) Index(i int) Int {
	if i < 0 || i >= len(s.slice) {
		return Int{}
	}
	return s.slice[i]
}

func (s *IntSeq) Time(t time.Time) Ints {
	index := s.timeIndex[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make(Ints, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

func (s *IntSeq) Value(v int) Ints {
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Ints, len(index))
	for i, v := range index {
		ret[i] = s.slice[v]
	}
	return ret
}

func (s *IntSeq) Visit(fn func(i int, v Int) (stop bool)) {
	for i, v := range s.slice {
		if fn != nil && fn(i, v) {
			break
		}
	}
}

func (s *IntSeq) Sum() int {
	var ret int
	for _, v := range s.slice {
		ret += v.Value
	}
	return ret
}

func (s *IntSeq) Count() int {
	return len(s.slice)
}

func (s *IntSeq) Max() Int {
	var max Int
	found := false
	for _, v := range s.slice {
		if !found {
			max = v
			found = true
		} else if v.Value < max.Value {
			max = v
		}
	}
	return max
}

func (s *IntSeq) Min() Int {
	var min Int
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

func (s *IntSeq) First() Int {
	if len(s.slice) == 0 {
		return Int{}
	}
	return s.slice[0]
}

func (s *IntSeq) Last() Int {
	if len(s.slice) == 0 {
		return Int{}
	}
	return s.slice[len(s.slice)-1]
}

func (s *IntSeq) Percentile(pct float64) Int {
	if len(s.slice) == 0 {
		return Int{}
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

func (s *IntSeq) Range(interval Interval) *IntSeq {
	slice := Range(s.slice, interval).(Ints)
	return newIntSeq(slice)
}

func (s *IntSeq) Merge(fn func(t time.Time, v1, v2 *int) *int, ints ...Ints) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	if len(ints) == 0 {
		return nil
	}

	slice1 := s.slice
	for _, slice2 := range ints {
		if !IsSorted(slice2) {
			temp := make(Ints, len(slice2))
			copy(temp, slice2)
			Sort(temp)
			slice2 = temp
		}
		var got Ints
		for i1, i2 := 0, 0; i1 < len(slice1) || i2 < len(slice2); {
			var (
				t time.Time
				v *int
			)
			switch {
			case i1 == len(slice1):
				t = slice2[i2].Time
				v2 := slice2[i2].Value
				v = fn(t, nil, &v2)
				i2++
			case i2 == len(slice2):
				t = slice2[i1].Time
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
				t = slice1[i2].Time
				v2 := slice2[i2].Value
				v = fn(t, nil, &v2)
				i2++
			}
			if v != nil {
				got = append(got, Int{
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

func (s *IntSeq) Aggregate(fn func(t time.Time, ints Ints) *int, duration time.Duration, begin, end *time.Time) error {
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

	slice := Ints{}
	ints := Ints{}
	i := 0
	for t := bg; t.Before(ed); t = t.Add(duration) {
		ints = ints[:0]
		for i < s.slice.Len() &&
			!s.slice[i].Time.After(t) &&
			s.slice[i].Time.Before(t.Add(duration)) {
			ints = append(ints, s.slice[i])
			i++
		}
		v := fn(t, ints)
		if v != nil {
			slice = append(slice, Int{
				Time:  t,
				Value: *v,
			})
		}
	}

	s.slice = slice
	s.buildIndex()
	return nil
}

func (s *IntSeq) Trim(fn func(i int, v Int) bool) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	updated := false
	slice := s.slice[:0]
	for i, v := range s.slice {
		if fn(i, v) {
			updated = true
		} else {
			slice = append(slice, v)
		}
	}

	if !updated {
		s.slice = slice
		s.buildIndex()
	}
	return nil
}
