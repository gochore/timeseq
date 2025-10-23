package timeseq

import "time"

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
