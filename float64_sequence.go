package timeseq

import (
	"errors"
	"sort"
	"time"
)

type Float64Item struct {
	Time  time.Time
	Value float64
}

type Float64Sequence []*Float64Item

func (s Float64Sequence) Len() int {
	return len(s)
}

func (s Float64Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Float64Sequence) Time(i int) time.Time {
	return s[i].Time
}

func (s Float64Sequence) Slice(i, j int) Sequence {
	return s[i:j]
}

func (s Float64Sequence) Sort() {
	Sort(s)
}

func (s Float64Sequence) Range(afterOrEqual, beforeOrEqual *time.Time) Float64Sequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(Float64Sequence)
}

func (s Float64Sequence) First(afterOrEqual *time.Time) *Float64Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	i := First(s, afterOrEqual)
	if i < 0 {
		return nil
	}
	return s[i]
}

func (s Float64Sequence) Last(beforeOrEqual *time.Time) *Float64Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	i := Last(s, beforeOrEqual)
	if i < 0 {
		return nil
	}
	return s[i]
}

func (s Float64Sequence) Max() *Float64Item {
	var max *Float64Item
	for _, v := range s {
		if max == nil {
			max = v
		} else if v.Value > max.Value {
			max = v
		}
	}
	return max
}

func (s Float64Sequence) Min() *Float64Item {
	var min *Float64Item
	for _, v := range s {
		if min == nil {
			min = v
		} else if v.Value < min.Value {
			min = v
		}
	}
	return min
}

func (s Float64Sequence) Sum() float64 {
	var sum float64
	for _, v := range s {
		sum += v.Value
	}
	return sum
}

func (s Float64Sequence) Average() float64 {
	if len(s) == 0 {
		return 0
	}

	return s.Sum() / float64(len(s))
}

func (s Float64Sequence) Percentile(pct float64) float64 {
	if pct > 1 || pct < 0 {
		panic(errors.New("percentile must be [0, 1]"))
	}

	var values []float64
	for _, v := range s {
		values = append(values, v.Value)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	if len(values) == 0 {
		return 0
	}

	index := int(float64(len(s))*pct - 1)
	if index < 0 {
		index = 0
	}

	return values[index]
}
