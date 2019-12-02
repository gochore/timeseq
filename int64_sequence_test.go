package timeseq

import (
	"reflect"
	"testing"
	"time"
)

func TestInt64Sequence_Len(t *testing.T) {
	tests := []struct {
		name string
		s    Int64Sequence
		want int
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Swap(tt.args.i, tt.args.j)
		})
	}
}

func TestInt64Sequence_Time(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		s    Int64Sequence
		args args
		want time.Time
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Sort()
		})
	}
}

func TestInt64Sequence_Range(t *testing.T) {
	type args struct {
		afterOrEqual  time.Time
		beforeOrEqual time.Time
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
		name  string
		s     Int64Sequence
		want  int
		want1 int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.Max()
			if got != tt.want {
				t.Errorf("Int64Sequence.Max() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Int64Sequence.Max() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInt64Sequence_Min(t *testing.T) {
	tests := []struct {
		name  string
		s     Int64Sequence
		want  int
		want1 int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.Min()
			if got != tt.want {
				t.Errorf("Int64Sequence.Min() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Int64Sequence.Min() got1 = %v, want %v", got1, tt.want1)
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
