package timeseq

import (
	"slices"
	"sort"
	"sync"
	"time"
)

type Number interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Point[T Number] struct {
	Time  time.Time
	Value T
}

func (i Point[T]) IsZero() bool {
	var zero T
	return i.Time.IsZero() && i.Value == zero
}

type Seq[T Number] struct {
	points []Point[T]

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[T][]int
	valueOrder []int
}

// NewSeq returns Seq with copied points inside
func NewSeq[T Number](points []Point[T]) *Seq[T] {
	temp := make([]Point[T], len(points))
	copy(temp, points)
	return WrapSeq(temp)
}

// WrapSeq returns Seq with origin points inside
func WrapSeq[T Number](points []Point[T]) *Seq[T] {
	if !slices.IsSortedFunc(points, compareItems[T]) {
		slices.SortStableFunc(points, compareItems[T])
	}
	return &Seq[T]{
		points: points,
	}
}

// ConvertSeq converts a slice of any type to Seq using the convert function
func ConvertSeq[T Number, S []P, P any](s S, convert func(p P) Point[T]) *Seq[T] {
	points := make([]Point[T], 0, len(s))
	for _, v := range s {
		points = append(points, convert(v))
	}
	return WrapSeq(points)
}

func compareItems[T Number](a, b Point[T]) int {
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
		timeIndex := make(map[timeKey][]int, len(s.points))
		valueIndex := make(map[T][]int, len(s.points))
		valueSlice := s.valueOrder[:0]
		for i, v := range s.points {
			k := newTimeKey(v.Time)
			timeIndex[k] = append(timeIndex[k], i)
			valueIndex[v.Value] = append(valueIndex[v.Value], i)
			valueSlice = append(valueSlice, i)
		}
		sort.SliceStable(valueSlice, func(i, j int) bool {
			return s.points[valueSlice[i]].Value < s.points[valueSlice[j]].Value
		})
		s.timeIndex = timeIndex
		s.valueIndex = valueIndex
		s.valueOrder = valueSlice
	})
}

// Points returns a replica of inside points
func (s *Seq[T]) Points() []Point[T] {
	ret := make([]Point[T], len(s.points))
	copy(ret, s.points)
	return ret
}

// Len returns length of inside points
func (s *Seq[T]) Len() int {
	return len(s.points)
}

// Index returns point of inside points, returns zero if index is out of range
func (s *Seq[T]) Index(i int) Point[T] {
	if i < 0 || i >= len(s.points) {
		return Point[T]{}
	}
	return s.points[i]
}

// Time returns the first point with time t, returns zero if not found
func (s *Seq[T]) Time(t time.Time) Point[T] {
	got := s.MTime(t)
	if len(got) == 0 {
		return Point[T]{}
	}
	return got[0]
}

// MTime returns all points with time t, returns nil if not found
func (s *Seq[T]) MTime(t time.Time) []Point[T] {
	s.buildIndex()
	index := s.timeIndex[newTimeKey(t)]
	if len(index) == 0 {
		return nil
	}
	ret := make([]Point[T], len(index))
	for i, v := range index {
		ret[i] = s.points[v]
	}
	return ret
}

// Value returns the first point with value v, returns zero if not found
func (s *Seq[T]) Value(v T) Point[T] {
	got := s.MValue(v)
	if len(got) == 0 {
		return Point[T]{}
	}
	return got[0]
}

// MValue returns all points with value v, returns nil if not found
func (s *Seq[T]) MValue(v T) []Point[T] {
	s.buildIndex()
	index := s.valueIndex[v]
	if len(index) == 0 {
		return nil
	}
	ret := make([]Point[T], len(index))
	for i, v := range index {
		ret[i] = s.points[v]
	}
	return ret
}

// Traverse call fn for every point one by one, break if fn returns false
func (s *Seq[T]) Traverse(fn func(i int, v Point[T]) (stop bool)) {
	for i, v := range s.points {
		if !fn(i, v) {
			break
		}
	}
}

// Sum returns sum of all values
func (s *Seq[T]) Sum() T {
	var ret T
	for _, v := range s.points {
		ret += v.Value
	}
	return ret
}

// Max returns the point with max value, returns zero if empty
func (s *Seq[T]) Max() Point[T] {
	var ret Point[T]
	for _, v := range s.points {
		if ret.IsZero() {
			ret = v
		} else if v.Value > ret.Value {
			ret = v
		}
	}
	return ret
}

// Min returns the point with min value, returns zero if empty
func (s *Seq[T]) Min() Point[T] {
	var ret Point[T]
	for _, v := range s.points {
		if ret.IsZero() {
			ret = v
		} else if v.Value < ret.Value {
			ret = v
		}
	}
	return ret
}

// First returns the first point, returns zero if empty
func (s *Seq[T]) First() Point[T] {
	if len(s.points) == 0 {
		return Point[T]{}
	}
	return s.points[0]
}

// Last returns the last point, returns zero if empty
func (s *Seq[T]) Last() Point[T] {
	if len(s.points) == 0 {
		return Point[T]{}
	}
	return s.points[len(s.points)-1]
}

// Percentile returns the point matched with percentile pct, returns zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Seq[T]) Percentile(pct float64) Point[T] {
	s.buildIndex()
	if len(s.points) == 0 {
		return Point[T]{}
	}
	if pct > 1 {
		pct = 1
	}
	if pct < 0 {
		pct = 0
	}
	i := int(float64(len(s.points))*pct - 1)
	if i < 0 {
		i = 0
	}
	return s.points[s.valueOrder[i]]
}

// Range returns a sub *Seq with specified interval
func (s *Seq[T]) Range(interval Interval) *Seq[T] {
	i := 0
	if interval.NotBefore != nil {
		i, _ = slices.BinarySearchFunc(s.points, Point[T]{Time: *interval.NotBefore}, func(a, b Point[T]) int {
			return compareItems[T](a, b)
		})
	}
	j := len(s.points)
	if interval.NotAfter != nil {
		j, _ = slices.BinarySearchFunc(s.points, Point[T]{Time: *interval.NotAfter}, func(a, b Point[T]) int {
			return compareItems[T](a, b)
		})

		if j < len(s.points) && s.points[j].Time.Equal(*interval.NotAfter) {
			j++
		}
	}
	return &Seq[T]{
		points: s.points[i:j],
	}
}

// Slice returns a sub Seq with specified index,
// (1, 2) means [1:2], (-1, 2) means [:2], (-1, -1) means [:]
func (s *Seq[T]) Slice(i, j int) *Seq[T] {
	if i < 0 && j < 0 {
		return s
	}
	points := s.points
	if i < 0 {
		points = points[:j]
	} else if j < 0 {
		points = points[i:]
	} else {
		points = points[i:j]
	}
	return &Seq[T]{
		points: points,
	}
}

// Filter returns a Seq with points which make fn returns true
func (s *Seq[T]) Filter(fn func(i int, v Point[T]) bool) *Seq[T] {
	points := make([]Point[T], 0, len(s.points))
	for i, v := range s.points {
		if fn(i, v) {
			points = append(points, v)
		}
	}
	return &Seq[T]{
		points: points,
	}
}

// Merge merges two Seq into a new Seq according to the specified rule
func Merge[T Number](s1, s2 *Seq[T], fn func(t time.Time, v1, v2 *T) *T) *Seq[T] {
	ret := make([]Point[T], 0, max(len(s1.points), len(s2.points)))

	for i1, i2 := 0, 0; i1 < len(s1.points) || i2 < len(s2.points); {
		var (
			t time.Time
			v *T
		)
		switch {
		case i1 == len(s1.points) || s1.points[i1].Time.After(s2.points[i2].Time):
			t = s2.points[i2].Time
			v2 := s2.points[i2].Value
			v = fn(t, nil, &v2)
			i2++
		case i2 == len(s2.points) || s1.points[i1].Time.Before(s2.points[i2].Time):
			t = s1.points[i1].Time
			v1 := s1.points[i1].Value
			v = fn(t, &v1, nil)
			i1++
		case s1.points[i1].Time.Equal(s2.points[i2].Time):
			t = s1.points[i1].Time
			v1 := s1.points[i1].Value
			v2 := s2.points[i2].Value
			v = fn(t, &v1, &v2)
			i1++
			i2++
		}
		if v != nil {
			ret = append(ret, Point[T]{
				Time:  t,
				Value: *v,
			})
		}
	}

	return &Seq[T]{
		points: ret,
	}
}

// Aggregate returns an aggregated Seq according to the specified rule
func (s *Seq[T]) Aggregate(fn func(t time.Time, points []Point[T]) *T, duration time.Duration) *Seq[T] {
	var (
		ret  []Point[T]
		temp []Point[T]
	)

	if duration <= 0 {
		for i := 0; i < len(s.points); {
			t := s.points[i].Time
			temp = temp[:0]
			for i < len(s.points) && t.Equal(s.points[i].Time) {
				temp = append(temp, s.points[i])
				i++
			}
			v := fn(t, temp)
			if v != nil {
				ret = append(ret, Point[T]{
					Time:  t,
					Value: *v,
				})
			}
		}
		return &Seq[T]{
			points: ret,
		}
	}

	var (
		begin, end time.Time
	)
	if len(s.points) > 0 {
		begin = s.points[0].Time.Truncate(duration)
		end = s.points[len(s.points)-1].Time.Truncate(duration)
	}

	for t, i := begin, 0; !t.After(end); t = t.Add(duration) {
		temp = temp[:0]
		nextT := t.Add(duration)
		for i < len(s.points) {
			if !s.points[i].Time.Before(nextT) {
				break
			}
			temp = append(temp, s.points[i])
			i++
		}
		v := fn(t, temp)
		if v != nil {
			ret = append(ret, Point[T]{
				Time:  t,
				Value: *v,
			})
		}
	}

	return &Seq[T]{
		points: ret,
	}
}
