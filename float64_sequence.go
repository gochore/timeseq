package timeseq

import (
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

func (s Float64Sequence) Range(afterOrEqual, beforeOrEqual time.Time) Float64Sequence {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	return Range(s, afterOrEqual, beforeOrEqual).(Float64Sequence)
}

func (s Float64Sequence) First(afterOrEqual *time.Time) *Float64Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	if len(s) == 0 {
		return nil
	}
	if afterOrEqual == nil {
		return s[0]
	}
	i := First(s, *afterOrEqual)
	if i < 0 {
		return nil
	}
	return s[i]
}

func (s Float64Sequence) Last(beforeOrEqual *time.Time) *Float64Item {
	if !sort.IsSorted(sortableSequence{s}) {
		s.Sort()
	}
	if len(s) == 0 {
		return nil
	}
	if beforeOrEqual == nil {
		return s[len(s)-1]
	}
	i := Last(s, *beforeOrEqual)
	if i < 0 {
		return nil
	}
	return s[i]
}

func (s Float64Sequence) Max() (int, float64) {
	if len(s) == 0 {
		return -1, 0
	}

	index, max := 0, s[0].Value
	for i, v := range s {
		if v.Value > max {
			index = i
			max = v.Value
		}
	}

	return index, max
}

func (s Float64Sequence) Min() (int, float64) {
	if len(s) == 0 {
		return -1, 0
	}

	index, min := 0, s[0].Value
	for i, v := range s {
		if v.Value < min {
			index = i
			min = v.Value
		}
	}

	return index, min
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
	if pct > 1 {
		pct = 1
	}
	if pct < 0 {
		pct = 0
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
