// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/gochore/pt"
)

func TestInt16Sequence_Len(t *testing.T) {
	tests := []struct {
		name string
		s    Int16Sequence
		want int
	}{
		{
			s:    RandomInt16Sequence(10),
			want: 10,
		},
		{
			s:    RandomInt16Sequence(0),
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
				t.Errorf("Int16Sequence.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    Int16Sequence
		args args
	}{
		{
			s: RandomInt16Sequence(10),
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
				t.Errorf("Int16Sequence.Swap() failed")
			}
		})
	}
}

func TestInt16Sequence_Time(t *testing.T) {
	seq := RandomInt16Sequence(10)
	seq.Sort()

	type args struct {
		i int
	}
	tests := []struct {
		name string
		s    Int16Sequence
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
				t.Errorf("Int16Sequence.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Slice(t *testing.T) {
	seq := RandomInt16Sequence(10)

	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    Int16Sequence
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
				t.Errorf("Int16Sequence.Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Sort(t *testing.T) {
	tests := []struct {
		name string
		s    Int16Sequence
	}{
		{
			s: RandomInt16Sequence(10),
		},
		{
			s: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Sort()
			if !sort.IsSorted(sortableSequence{tt.s}) {
				t.Error("Int16Sequence.Slice() failed")
			}
		})
	}
}

func TestInt16Sequence_Range(t *testing.T) {
	now := time.Now()
	seq := RandomInt16Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}
	rand.Shuffle(seq.Len(), func(i, j int) {
		seq[i], seq[j] = seq[j], seq[i]
	})

	type args struct {
		afterOrEqual  *time.Time
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Int16Sequence
		args args
		want Int16Sequence
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
				t.Errorf("Int16Sequence.Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_First(t *testing.T) {
	now := time.Now()
	seq := RandomInt16Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}
	rand.Shuffle(seq.Len(), func(i, j int) {
		seq[i], seq[j] = seq[j], seq[i]
	})

	type args struct {
		afterOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Int16Sequence
		args args
		want *Int16Item
	}{
		{
			s: seq,
			args: args{
				afterOrEqual: nil,
			},
			want: &seq[0],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now),
			},
			want: &seq[0],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(-1 * time.Second)),
			},
			want: &seq[0],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(1 * time.Second)),
			},
			want: &seq[1],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: &seq[5],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: &seq[5],
		},
		{
			s: seq,
			args: args{
				afterOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: &seq[6],
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
				t.Errorf("Int16Sequence.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Last(t *testing.T) {
	now := time.Now()
	seq := RandomInt16Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}
	rand.Shuffle(seq.Len(), func(i, j int) {
		seq[i], seq[j] = seq[j], seq[i]
	})

	type args struct {
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		s    Int16Sequence
		args args
		want *Int16Item
	}{
		{
			s: seq,
			args: args{
				beforeOrEqual: nil,
			},
			want: &seq[len(seq)-1],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now),
			},
			want: &seq[0],
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
			want: &seq[1],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: &seq[5],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: &seq[4],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: &seq[5],
		},
		{
			s: seq,
			args: args{
				beforeOrEqual: pt.Time(now.Add(100 * time.Second)),
			},
			want: &seq[9],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Last(tt.args.beforeOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16Sequence.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Max(t *testing.T) {
	seq1 := RandomInt16Sequence(10)
	seq1.Sort()
	for i := range seq1 {
		if i == 0 {
			seq1[i].Value = 1
		} else {
			seq1[i].Value = 0
		}
	}
	seq2 := RandomInt16Sequence(10)
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
		s    Int16Sequence
		want *Int16Item
	}{
		{
			s:    seq1,
			want: &seq1[0],
		},
		{
			s:    seq2,
			want: &seq2[1],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Max(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16Sequence.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Min(t *testing.T) {
	seq1 := RandomInt16Sequence(10)
	seq1.Sort()
	for i := range seq1 {
		if i == 0 {
			seq1[i].Value = 1
		} else {
			seq1[i].Value = 0
		}
	}
	seq2 := RandomInt16Sequence(10)
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
		s    Int16Sequence
		want *Int16Item
	}{
		{
			s:    seq1,
			want: &seq1[1],
		},
		{
			s:    seq2,
			want: &seq2[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Min(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16Sequence.Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Sum(t *testing.T) {
	seq := RandomInt16Sequence(10)
	for i := range seq {
		seq[i].Value = int16(i)
	}

	tests := []struct {
		name string
		s    Int16Sequence
		want int16
	}{
		{
			s:    seq,
			want: 45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Sum(); got != tt.want {
				t.Errorf("Int16Sequence.Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Average(t *testing.T) {
	seq := RandomInt16Sequence(10)
	for i := range seq {
		seq[i].Value = int16(i) * 2
	}

	tests := []struct {
		name string
		s    Int16Sequence
		want int16
	}{
		{
			s:    seq,
			want: 9,
		},
		{
			s:    nil,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Average(); got != tt.want {
				t.Errorf("Int16Sequence.Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Sequence_Percentile(t *testing.T) {
	seq := RandomInt16Sequence(100)
	for i := range seq {
		seq[i].Value = int16(i) + 1
	}

	type args struct {
		pct float64
	}
	tests := []struct {
		name string
		s    Int16Sequence
		args args
		want int16
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
		{
			s: nil,
			args: args{
				pct: 1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.args.pct > 1 || tt.args.pct < 0 {
						return
					}
					t.Errorf("Int16Sequence.Percentile() failed: %v", r)
				}
			}()
			if got := tt.s.Percentile(tt.args.pct); got != tt.want {
				t.Errorf("Int16Sequence.Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeInt16(t *testing.T) {
	seq := RandomInt16Sequence(10)
	seq.Sort()
	type args struct {
		seq1 Int16Sequence
		seq2 Int16Sequence
		fn   func(item1, item2 *Int16Item) *Int16Item
	}
	tests := []struct {
		name string
		args args
		want Int16Sequence
	}{
		{
			name: "regular",
			args: args{
				seq1: seq[0:7],
				seq2: seq[3:10],
				fn: func(item1, item2 *Int16Item) *Int16Item {
					if item1 != nil {
						return item1
					}
					return item2
				},
			},
			want: seq,
		},
		{
			name: "reverse",
			args: args{
				seq1: seq[3:10],
				seq2: seq[0:7],
				fn: func(item1, item2 *Int64Item) *Int64Item {
					if item1 != nil {
						return item1
					}
					return item2
				},
			},
			want: seq,
		},
		{
			name: "nil fn",
			args: args{
				seq1: seq[0:7],
				seq2: seq[3:10],
				fn:   nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeInt16(tt.args.seq1, tt.args.seq2, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func RandomInt16Sequence(length int) Int16Sequence {
	now := time.Now()
	ret := make(Int16Sequence, length)
	for i := range ret {
		ret[i] = Int16Item{
			Time:  now.Add(time.Duration(rand.Intn(length)) * time.Second),
			Value: int16(rand.Float64() * float64(math.MaxInt64)),
		}
	}
	return ret
}