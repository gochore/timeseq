package timeseq

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestInt64Sequence_Len(t *testing.T) {
	tests := []struct {
		name string
		s    Int64Sequence
		want int
	}{
		{
			s:    RandomInt64Sequence(10),
			want: 10,
		},
		{
			s:    RandomInt64Sequence(0),
			want: 0,
		},
		{
			s:    nil,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("Int64Sequence.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
	}{
		{
			s: RandomInt64Sequence(10),
			args: args{
				i: 0,
				j: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := tt.s.Time(tt.args.i)
			tj := tt.s.Time(tt.args.j)
			tt.s.Swap(tt.args.i, tt.args.j)
			if ti != tt.s.Time(tt.args.j) || tj != tt.s.Time(tt.args.i) {
				t.Errorf("Int64Sequence.Swap() failed")
			}
		})
	}
}

func TestInt64Sequence_Time(t *testing.T) {
	seq := RandomInt64Sequence(10)
	seq.Sort()

	type args struct {
		i int
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
		want time.Time
	}{
		{
			s: seq,
			args: args{
				i: 9,
			},
			want: seq[9].Time,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Time(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64Sequence.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Slice(t *testing.T) {
	seq := RandomInt64Sequence(10)

	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
		want Sequence
	}{
		{
			s: seq,
			args: args{
				i: 2,
				j: 10,
			},
			want: seq[2:10],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Slice(tt.args.i, tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64Sequence.Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Sort(t *testing.T) {
	tests := []struct {
		name string
		s    Int64Sequence
	}{
		{
			s: RandomInt64Sequence(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Sort()
			if !sort.IsSorted(sortableSequence{tt.s}) {
				t.Error("Int64Sequence.Slice() failed")
			}
		})
	}
}

func TestInt64Sequence_Range(t *testing.T) {
	type args struct {
		afterOrEqual  *time.Time
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
		want Int64Sequence
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Range(tt.args.afterOrEqual, tt.args.beforeOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64Sequence.Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_First(t *testing.T) {
	type args struct {
		afterOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
		want *Int64Item
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.First(tt.args.afterOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64Sequence.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Last(t *testing.T) {
	type args struct {
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
		want *Int64Item
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Last(tt.args.beforeOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64Sequence.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Max(t *testing.T) {
	tests := []struct {
		name string
		s    Int64Sequence
		want *Int64Item
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Max(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64Sequence.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Min(t *testing.T) {
	tests := []struct {
		name string
		s    Int64Sequence
		want *Int64Item
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Min(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64Sequence.Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Sum(t *testing.T) {
	tests := []struct {
		name string
		s    Int64Sequence
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Sum(); got != tt.want {
				t.Errorf("Int64Sequence.Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Average(t *testing.T) {
	tests := []struct {
		name string
		s    Int64Sequence
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Average(); got != tt.want {
				t.Errorf("Int64Sequence.Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Sequence_Percentile(t *testing.T) {
	type args struct {
		pct float64
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Percentile(tt.args.pct); got != tt.want {
				t.Errorf("Int64Sequence.Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}
