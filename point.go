package timeseq

import "time"

type Number interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Point[T Number] struct {
	Time  time.Time
	Value T
}

func NewPoint[T Number](t time.Time, v T) Point[T] {
	return Point[T]{
		Time:  t,
		Value: v,
	}
}

func (p Point[T]) IsZero() bool {
	var zero T
	return p.Time.IsZero() && p.Value == zero
}

func (p Point[T]) Equal(other Point[T]) bool {
	return p.Time.Equal(other.Time) && p.Value == other.Value
}
