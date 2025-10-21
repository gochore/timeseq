package timeseq

import (
	"slices"
	"sort"
	"sync"
	"time"
)

type Point[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64] struct {
	Time  time.Time
	Value T
}

func (i Point[T]) IsZero() bool {
	var zero T
	return i.Time.IsZero() && i.Value == zero
}

type Seq[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64] struct {
	points []Point[T]

	indexOnce  sync.Once
	timeIndex  map[timeKey][]int
	valueIndex map[T][]int
	valueOrder []int
}

// NewSeq returns Seq with copied points inside
func NewSeq[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](points []Point[T]) *Seq[T] {
	temp := make([]Point[T], len(points))
	copy(temp, points)
	return WrapSeq(temp)
}

// WrapSeq returns Seq with origin points inside
func WrapSeq[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](points []Point[T]) *Seq[T] {
	if !slices.IsSortedFunc(points, compareItems[T]) {
		slices.SortStableFunc(points, compareItems[T])
	}
	return &Seq[T]{
		points: points,
	}
}

func compareItems[T float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](a, b Point[T]) int {
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

// Slice returns a replica of inside points
func (s *Seq[T]) Slice() []Point[T] {
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
func (s *Float32Seq) Min() Float32 {
	var min Float32
	found := false
	spoints := s.getSlice()
	for _, v := range spoints {
		if !found {
			min = v
			found = true
		} else if v.Value < min.Value {
			min = v
		}
	}
	return min
}

// First returns the first point, returns zero if empty
func (s *Float32Seq) First() Float32 {
	spoints := s.getSlice()
	if len(spoints) == 0 {
		return Float32{}
	}
	return spoints[0]
}

// Last returns the last point, returns zero if empty
func (s *Float32Seq) Last() Float32 {
	spoints := s.getSlice()
	if len(spoints) == 0 {
		return Float32{}
	}
	return spoints[len(spoints)-1]
}

// Percentile returns the point matched with percentile pct, returns zero if empty,
// the pct's valid range is be [0, 1], it will be treated as 1 if greater than 1, as 0 if smaller than 0
func (s *Float32Seq) Percentile(pct float64) Float32 {
	s.buildIndex()
	spoints := s.getSlice()
	if len(spoints) == 0 {
		return Float32{}
	}
	if pct > 1 {
		pct = 1
	}
	if pct < 0 {
		pct = 0
	}
	i := int(float64(len(spoints))*pct - 1)
	if i < 0 {
		i = 0
	}
	return spoints[s.valueOrder[i]]
}

// Range returns a sub *Float32Seq with specified interval
func (s *Float32Seq) Range(interval Interval) *Float32Seq {
	spoints := s.getSlice()
	points := Range(spoints, interval).(Float32s)
	return newFloat32Seq(points)
}

// Slice returns a sub *Float32Seq with specified index,
// (1, 2) means [1:2], (-1, 2) means [:2], (-1, -1) means [:]
func (s *Float32Seq) Slice(i, j int) *Float32Seq {
	if i < 0 && j < 0 {
		return s
	}
	spoints := s.getSlice()
	if i < 0 {
		return newFloat32Seq(spoints[:j])
	}
	if j < 0 {
		return newFloat32Seq(spoints[i:])
	}
	return newFloat32Seq(spoints[i:j])
}

// Trim returns a *Float32Seq without points which make fn returns true
func (s *Float32Seq) Trim(fn func(i int, v Float32) bool) *Float32Seq {
	spoints := s.getSlice()
	if fn == nil || len(spoints) == 0 {
		return s
	}

	removeM := map[int]struct{}{}
	for i, v := range spoints {
		if fn(i, v) {
			removeM[i] = struct{}{}
		}
	}
	if len(removeM) == 0 {
		return s
	}

	points := make(Float32s, 0, len(spoints)-len(removeM))
	for i, v := range spoints {
		if _, ok := removeM[i]; ok {
			continue
		}
		points = append(points, v)
	}

	return newFloat32Seq(points)
}

// Merge returns a new *Float32Seq with merged data according to the specified rule
func (s *Float32Seq) Merge(fn func(t time.Time, v1, v2 *float32) *float32, seq *Float32Seq) *Float32Seq {
	if fn == nil {
		return s
	}

	var ret Float32s

	points1 := s.getSlice()
	var points2 Float32s
	if seq != nil {
		points2 = seq.getSlice()
	}

	for i1, i2 := 0, 0; i1 < len(points1) || i2 < len(points2); {
		var (
			t time.Time
			v *float32
		)
		switch {
		case i1 == len(points1):
			t = points2[i2].Time
			v2 := points2[i2].Value
			v = fn(t, nil, &v2)
			i2++
		case i2 == len(points2):
			t = points1[i1].Time
			v1 := points1[i1].Value
			v = fn(t, &v1, nil)
			i1++
		case points1[i1].Time.Equal(points2[i2].Time):
			t = points1[i1].Time
			v1 := points1[i1].Value
			v2 := points2[i2].Value
			v = fn(t, &v1, &v2)
			i1++
			i2++
		case points1[i1].Time.Before(points2[i2].Time):
			t = points1[i1].Time
			v1 := points1[i1].Value
			v = fn(t, &v1, nil)
			i1++
		case points1[i1].Time.After(points2[i2].Time):
			t = points2[i2].Time
			v2 := points2[i2].Value
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
func (s *Float32Seq) Aggregate(fn func(t time.Time, points Float32s) *float32, duration time.Duration, interval Interval) *Float32Seq {
	if fn == nil {
		return s
	}

	ret := Float32s{}
	temp := Float32s{}

	spoints := s.getSlice()
	if duration <= 0 {
		for i := 0; i < s.Len(); {
			t := spoints[i].Time
			if !interval.Contain(t) {
				i++
				continue
			}
			temp = temp[:0]
			for i < spoints.Len() && t.Equal(spoints[i].Time) {
				temp = append(temp, spoints[i])
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

	if len(spoints) == 0 && interval.Duration() < 0 {
		return s
	}

	var begin time.Time
	if len(spoints) > 0 {
		begin = spoints[0].Time.Truncate(duration)
	}
	if interval.NotBefore != nil {
		begin = (*interval.NotBefore).Truncate(duration)
		if begin.Before(*interval.NotBefore) {
			begin = begin.Add(duration)
		}
	}

	var end time.Time
	if len(spoints) > 0 {
		end = spoints[len(spoints)-1].Time.Truncate(duration)
	}
	if interval.NotAfter != nil {
		end = (*interval.NotAfter).Truncate(duration)
	}

	for t, i := begin, 0; !t.After(end); t = t.Add(duration) {
		temp = temp[:0]
		itv := BeginAt(t).EndAt(t.Add(duration))
		for i < len(spoints) {
			if spoints[i].Time.After(*itv.NotAfter) {
				break
			}
			if itv.Contain(spoints[i].Time) {
				temp = append(temp, spoints[i])
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
