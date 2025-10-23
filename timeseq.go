package timeseq

import (
	"slices"
	"time"
)

// NewSeq returns Seq with points inside
func NewSeq[T Number](points []Point[T]) *Seq[T] {
	if !slices.IsSortedFunc(points, compareItems[T]) {
		slices.SortStableFunc(points, compareItems[T])
	}
	return &Seq[T]{
		points: points,
	}
}

// NewSeqCopy returns Seq with copied points inside
func NewSeqCopy[T Number](points []Point[T]) *Seq[T] {
	temp := make([]Point[T], len(points))
	copy(temp, points)
	return NewSeq(temp)
}

// NewSeqConvert converts a slice of any type to Seq using the convert function
func NewSeqConvert[T Number, S []P, P any](s S, convert func(p P) Point[T]) *Seq[T] {
	points := make([]Point[T], 0, len(s))
	for _, v := range s {
		points = append(points, convert(v))
	}
	return NewSeq(points)
}

// Convert converts a Seq of one type to a Seq of another type using the convert function
func Convert[T Number, P Number](s *Seq[P], convert func(v Point[P]) Point[T]) *Seq[T] {
	return NewSeqConvert(s.points, convert)
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
		case i1 == len(s1.points):
			t = s2.points[i2].Time
			v2 := s2.points[i2].Value
			v = fn(t, nil, &v2)
			i2++
		case i2 == len(s2.points):
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
		case s1.points[i1].Time.After(s2.points[i2].Time):
			t = s2.points[i2].Time
			v2 := s2.points[i2].Value
			v = fn(t, nil, &v2)
			i2++
		case s1.points[i1].Time.Before(s2.points[i2].Time):
			t = s1.points[i1].Time
			v1 := s1.points[i1].Value
			v = fn(t, &v1, nil)
			i1++
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
