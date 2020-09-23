package timeseq

import (
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

type IntSeq struct {
	slice Ints
	index map[timeKey][]int
}

func NewIntSeq(slice Ints) *IntSeq {
	s := make(Ints, len(slice))
	copy(s, slice)
	Sort(slice)
	index := make(map[timeKey][]int, len(slice))
	key := timeKey{}
	for i, v := range slice {
		key.Put(v.Time)
		index[key] = append(index[key], i)
	}
	return &IntSeq{
		slice: slice,
		index: index,
	}
}

func (s *IntSeq) Ints() Ints {
	slice := make(Ints, len(s.slice))
	copy(slice, s.slice)
	return slice
}

func (s *IntSeq) Index(i int) Int {
	panic("TODO")
}

func (s *IntSeq) Time(t time.Time) Ints {
	panic("TODO")
}

func (s *IntSeq) Visit(fn func(i int, v Int)) {
	panic("TODO")
}

func (s *IntSeq) Sum() int {
	panic("TODO")
}

func (s *IntSeq) Count() int {
	panic("TODO")
}

func (s *IntSeq) Max() (int, Int) {
	panic("TODO")
}

func (s *IntSeq) Min() (int, Int) {
	panic("TODO")
}

func (s *IntSeq) First() (int, Int) {
	panic("TODO")
}

func (s *IntSeq) Last() (int, Int) {
	panic("TODO")
}

func (s IntSeq) Percentile(pct float64) (int, Int) {
	panic("TODO")
}

func (s IntSeq) Range(interval Interval) IntSeq {
	panic("TODO")
}

func (s IntSeq) Merge(fn func(e1, e2 *Int) *Int, slice ...Ints) error {
	panic("TODO")
}

func (s IntSeq) Aggregate(fn func(t time.Time, es ...Int) *Int, duration time.Duration) error {
	panic("TODO")
}

func (s IntSeq) Trim(trim func(i int, v Int) bool) error {
	panic("TODO")
}
