package timeseq

import (
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/gochore/pt"
)

func TestFloat64Sequence_Len(t *testing.T) {
	tests := []struct {
		name string
		s    Float64Sequence
		want int
	}{
		{
			s:    RandomFloat64Sequence(10),
			want: 10,
		},
		{
			s:    RandomFloat64Sequence(0),
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
				t.Errorf("Float64Sequence.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    Float64Sequence
		args args
	}{
		{
			s: RandomFloat64Sequence(10),
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
				t.Errorf("Float64Sequence.Swap() failed")
			}
		})
	}
}

func TestFloat64Sequence_Time(t *testing.T) {
	seq := RandomFloat64Sequence(10)
	seq.Sort()

	type args struct {
		i int
	}
	tests := []struct {
		name string
		s    Float64Sequence
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
				t.Errorf("Float64Sequence.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Slice(t *testing.T) {
	seq := RandomFloat64Sequence(10)

	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    Float64Sequence
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
				t.Errorf("Float64Sequence.Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Sort(t *testing.T) {
	tests := []struct {
		name string
		s    Float64Sequence
	}{
		{
			s: RandomFloat64Sequence(10),
		},
		{
			s: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Sort()
			if !sort.IsSorted(sortableSequence{tt.s}) {
				t.Error("Float64Sequence.Slice() failed")
			}
		})
	}
}

func TestFloat64Sequence_Range(t *testing.T) {
	now := time.Now()
	seq := RandomFloat64Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}

	type args struct {
		afterOrEqual  *time.Time
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Float64Sequence
		args args
		want Float64Sequence
	}{
		{
			s: seq,
			args: args{
				afterOrEqual:  pt.Time(now.Add(1 * time.Second)),
				beforeOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: seq[1 : 5+1],
		},
		{
			s: seq,
			args: args{
				afterOrEqual:  pt.Time(now),
				beforeOrEqual: pt.Time(now.Add(100 * time.Second)),
			},
			want: seq,
		},
		{
			s: seq,
			args: args{
				afterOrEqual:  pt.Time(now.Add(1*time.Second + time.Millisecond)),
				beforeOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: seq[2 : 4+1],
		},
		{
			s: seq,
			args: args{
				afterOrEqual:  pt.Time(now.Add(1*time.Second - time.Millisecond)),
				beforeOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: seq[1 : 5+1],
		},
		{
			s: seq,
			args: args{
				afterOrEqual:  nil,
				beforeOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: seq[:5+1],
		},
		{
			s: seq,
			args: args{
				afterOrEqual:  pt.Time(now.Add(1 * time.Second)),
				beforeOrEqual: nil,
			},
			want: seq[1:],
		},
		{
			s: seq,
			args: args{
				afterOrEqual:  nil,
				beforeOrEqual: nil,
			},
			want: seq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Range(tt.args.afterOrEqual, tt.args.beforeOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64Sequence.Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_First(t *testing.T) {
	now := time.Now()
	seq := RandomFloat64Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}

	type args struct {
		afterOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Float64Sequence
		args args
		want *Float64Item
	}{
		{
			s: seq,
			args: args{
				afterOrEqual: nil,
			},
			want: seq[0],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now),
			},
			want: seq[0],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(-1 * time.Second)),
			},
			want: seq[0],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(1 * time.Second)),
			},
			want: seq[1],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: seq[5],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: seq[5],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: seq[6],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(100 * time.Second)),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.First(tt.args.afterOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64Sequence.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Last(t *testing.T) {
	now := time.Now()
	seq := RandomFloat64Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}

	type args struct {
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Float64Sequence
		args args
		want *Float64Item
	}{
		{
			s: seq,
			args: args{
				beforeOrEqual: nil,
			},
			want: seq[len(seq)-1],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now),
			},
			want: seq[0],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(-1 * time.Second)),
			},
			want: nil,
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(1 * time.Second)),
			},
			want: seq[1],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: seq[5],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: seq[4],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: seq[5],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(100 * time.Second)),
			},
			want: seq[9],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Last(tt.args.beforeOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64Sequence.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Max(t *testing.T) {
	seq1 := RandomFloat64Sequence(10)
	seq1.Sort()
	for i := range seq1 {
		if i == 0 {
			seq1[i].Value = 1
		} else {
			seq1[i].Value = 0
		}
	}
	seq2 := RandomFloat64Sequence(10)
	seq2.Sort()
	for i := range seq2 {
		if i == 0 {
			seq2[i].Value = 0
		} else {
			seq2[i].Value = 1
		}
	}

	tests := []struct {
		name string
		s    Float64Sequence
		want *Float64Item
	}{
		{
			s:    seq1,
			want: seq1[0],
		},
		{
			s:    seq2,
			want: seq2[1],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Max(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64Sequence.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Min(t *testing.T) {
	seq1 := RandomFloat64Sequence(10)
	seq1.Sort()
	for i := range seq1 {
		if i == 0 {
			seq1[i].Value = 1
		} else {
			seq1[i].Value = 0
		}
	}
	seq2 := RandomFloat64Sequence(10)
	seq2.Sort()
	for i := range seq2 {
		if i == 0 {
			seq2[i].Value = 0
		} else {
			seq2[i].Value = 1
		}
	}

	tests := []struct {
		name string
		s    Float64Sequence
		want *Float64Item
	}{
		{
			s:    seq1,
			want: seq1[1],
		},
		{
			s:    seq2,
			want: seq2[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Min(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64Sequence.Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Sum(t *testing.T) {
	seq := RandomFloat64Sequence(100)
	for i := range seq {
		seq[i].Value = float64(i)
	}

	tests := []struct {
		name string
		s    Float64Sequence
		want float64
	}{
		{
			s:    seq,
			want: 4950,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Sum(); got != tt.want {
				t.Errorf("Float64Sequence.Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Average(t *testing.T) {
	seq := RandomFloat64Sequence(100)
	for i := range seq {
		seq[i].Value = float64(i)
	}

	tests := []struct {
		name string
		s    Float64Sequence
		want float64
	}{
		{
			s:    seq,
			want: 49.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Average(); got != tt.want {
				t.Errorf("Float64Sequence.Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Sequence_Percentile(t *testing.T) {
	seq := RandomFloat64Sequence(100)
	for i := range seq {
		seq[i].Value = float64(i) + 1
	}

	type args struct {
		pct float64
	}
	tests := []struct {
		name string
		s    Float64Sequence
		args args
		want float64
	}{
		{
			s: seq,
			args: args{
				pct: 0,
			},
			want: 1,
		},
		{
			s: seq,
			args: args{
				pct: 0.5,
			},
			want: 50,
		},
		{
			s: seq,
			args: args{
				pct: 0.95,
			},
			want: 95,
		},
		{
			s: seq,
			args: args{
				pct: 0.955,
			},
			want: 95,
		},
		{
			s: seq,
			args: args{
				pct: 1,
			},
			want: 100,
		},
		{
			s: seq,
			args: args{
				pct: 1.1,
			},
			want: 100,
		},
		{
			s: seq,
			args: args{
				pct: -0.1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.args.pct > 1 || tt.args.pct < 0 {
						return
					}
					t.Errorf("Float64Sequence.Percentile() failed: %v", r)
				}
			}()
			if got := tt.s.Percentile(tt.args.pct); got != tt.want {
				t.Errorf("Float64Sequence.Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFloat64Sequence_GC(b *testing.B) {
	b.ReportAllocs()

	length := 1000
	for i := 0; i < b.N; i++ {
		_ = RandomFloat64Sequence(length)
		runtime.GC()
	}
}

func RandomFloat64Sequence(length int) Float64Sequence {
	now := time.Now()
	ret := make(Float64Sequence, length)
	for i := range ret {
		ret[i] = &Float64Item{
			Time:  now.Add(time.Duration(rand.Intn(length)) * time.Second),
			Value: rand.NormFloat64(),
		}
	}
	return ret
}