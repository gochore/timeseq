package timeseq

import (
	"fmt"
	"sort"
	"time"
)

type Int struct {
	Time  time.Time
	Value int
}

type Ints []Int

type IntSeq struct {
	slice      Ints
	timeIndex  map[timeKey][]int
	valueIndex map[int][]int
	valueSlice []int
}

func NewIntSeq(slice Ints) *IntSeq {
	s := make(Ints, len(slice))
	copy(s, slice)
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].Time.Before(slice[j].Time)
	})
	return newIntSeq(slice)
}

func newIntSeq(slice Ints) *IntSeq {
	timeIndex := make(map[timeKey][]int, len(slice))
	valueIndex := make(map[int][]int, len(slice))
	valueSlice := make([]int, len(slice))
	for i, v := range slice {
		k := newTimeKey(v.Time)
		timeIndex[k] = append(timeIndex[k], i)
		valueIndex[v.Value] = append(valueIndex[v.Value], i)
		valueSlice[i] = i
	}
	sort.SliceStable(valueSlice, func(i, j int) bool {
		return slice[valueSlice[i]].Value < slice[valueSlice[j]].Value
	})
	return &IntSeq{
		slice:      slice,
		timeIndex:  timeIndex,
		valueIndex: valueIndex,
		valueSlice: valueSlice,
	}
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
	for _, i := range index {
		ret[i] = s.slice[i]
	}
	return ret
}

func (s *IntSeq) Value(v int) Ints {
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make(Ints, len(index))
	for _, i := range index {
		ret[i] = s.slice[i]
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

func (s IntSeq) Percentile(pct float64) Int {
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

func (s IntSeq) Range(interval Interval) *IntSeq {
	i := 0
	if interval.NotBefore != nil {
		i = sort.Search(len(s.slice), func(i int) bool {
			return !s.slice[i].Time.Before(*interval.NotBefore)
		})
	}
	j := len(s.slice)
	if interval.NotAfter != nil {
		j = sort.Search(len(s.slice), func(j int) bool {
			return !s.slice[j].Time.Before(*interval.NotBefore)
		})
		if j < len(s.slice) && s.slice[j].Time.Equal(*interval.NotBefore) {
			j++
		}
	}
	return newIntSeq(s.slice[i:j])
}

func (s IntSeq) Merge(fn func(t time.Time, v1, v2 *int) *int, slices ...Ints) error {
	if len(slices) == 0 {
		return nil
	}

	if fn == nil {
		return fmt.Errorf("nil fn")
	}

	seq1 := s.slice
	for _, seq2 := range slices {
		var got Ints
		for i1, i2 := 0, 0; i1 < len(seq1) || i2 < len(seq2); {
			var (
				t time.Time
				v *int
			)
			switch {
			case i1 == len(seq1):
				t = seq2[i2].Time
				v2 := seq2[i2].Value
				v = fn(t, nil, &v2)
				i2++
			case i2 == len(seq2):
				t = seq2[i1].Time
				v1 := seq1[i1].Value
				v = fn(t, &v1, nil)
				i1++
			case seq1[i1].Time.Equal(seq2[i2].Time):
				t = seq1[i1].Time
				v1 := seq1[i1].Value
				v2 := seq2[i2].Value
				v = fn(t, &v1, &v2)
				i1++
				i2++
			case seq1[i1].Time.Before(seq2[i2].Time):
				t = seq1[i1].Time
				v1 := seq1[i1].Value
				v = fn(t, &v1, nil)
				i1++
			case seq1[i1].Time.After(seq2[i2].Time):
				t = seq1[i2].Time
				v2 := seq2[i2].Value
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
		seq1 = got
	}

	return got
}

func (s IntSeq) Aggregate(fn func(t time.Time, es ...Int) *Int, duration time.Duration) error {
	panic("TODO")
}

func (s IntSeq) Trim(trim func(i int, v Int) bool) error {
	panic("TODO")
}
