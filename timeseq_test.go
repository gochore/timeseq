package timeseq

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/gochore/pt"
)

func TestSort(t *testing.T) {
	seq := RandomInt64Sequence(10)

	type args struct {
		seq Sequence
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				seq: seq,
			},
		},
		{
			args: args{
				seq: nil,
			},
		},
		{
			args: args{
				seq: seq,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sort(tt.args.seq)
			if !sort.IsSorted(sortableSequence{seq}) {
				t.Errorf("Sort() failed")
			}
		})
	}
}

func TestRange(t *testing.T) {
	now := time.Now()
	seq := RandomInt64Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}

	type args struct {
		seq           Sequence
		afterOrEqual  *time.Time
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		args args
		want Sequence
	}{
		{
			args: args{
				seq:           seq,
				afterOrEqual:  pt.Time(now.Add(1 * time.Second)),
				beforeOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: seq[1 : 5+1],
		},
		{
			args: args{
				seq:           seq,
				afterOrEqual:  pt.Time(now),
				beforeOrEqual: pt.Time(now.Add(100 * time.Second)),
			},
			want: seq,
		},
		{
			args: args{
				seq:           seq,
				afterOrEqual:  pt.Time(now.Add(1*time.Second + time.Millisecond)),
				beforeOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: seq[2 : 4+1],
		},
		{
			args: args{
				seq:           seq,
				afterOrEqual:  pt.Time(now.Add(1*time.Second - time.Millisecond)),
				beforeOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: seq[1 : 5+1],
		},
		{
			args: args{
				seq:           seq,
				afterOrEqual:  nil,
				beforeOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: seq[:5+1],
		},
		{
			args: args{
				seq:           seq,
				afterOrEqual:  pt.Time(now.Add(1 * time.Second)),
				beforeOrEqual: nil,
			},
			want: seq[1:],
		},
		{
			args: args{
				seq:           seq,
				afterOrEqual:  nil,
				beforeOrEqual: nil,
			},
			want: seq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Range(tt.args.seq, tt.args.afterOrEqual, tt.args.beforeOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	now := time.Now()
	seq := RandomInt64Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}

	type args struct {
		seq          Sequence
		afterOrEqual *time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				seq:          seq,
				afterOrEqual: nil,
			},
			want: 0,
		},
		{
			args: args{
				seq:          seq,
				afterOrEqual: pt.Time(now),
			},
			want: 0,
		},
		{
			args: args{
				seq:          seq,
				afterOrEqual: pt.Time(now.Add(-1 * time.Second)),
			},
			want: 0,
		},
		{
			args: args{
				seq:          seq,
				afterOrEqual: pt.Time(now.Add(1 * time.Second)),
			},
			want: 1,
		},
		{
			args: args{
				seq:          seq,
				afterOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: 5,
		},
		{
			args: args{
				seq:          seq,
				afterOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: 5,
		},
		{
			args: args{
				seq:          seq,
				afterOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: 6,
		},
		{
			args: args{
				seq:          seq,
				afterOrEqual: pt.Time(now.Add(100 * time.Second)),
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := First(tt.args.seq, tt.args.afterOrEqual); got != tt.want {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLast(t *testing.T) {
	now := time.Now()
	seq := RandomInt64Sequence(10)
	for i := range seq {
		seq[i].Time = now.Add(time.Duration(i) * time.Second)
	}

	type args struct {
		seq           Sequence
		beforeOrEqual *time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				seq:           seq,
				beforeOrEqual: nil,
			},
			want: len(seq) - 1,
		},
		{
			args: args{
				seq:           seq,
				beforeOrEqual: pt.Time(now),
			},
			want: 0,
		},
		{
			args: args{
				seq:           seq,
				beforeOrEqual: pt.Time(now.Add(-1 * time.Second)),
			},
			want: -1,
		},
		{
			args: args{
				seq:           seq,
				beforeOrEqual: pt.Time(now.Add(1 * time.Second)),
			},
			want: 1,
		},
		{
			args: args{
				seq:           seq,
				beforeOrEqual: pt.Time(now.Add(5 * time.Second)),
			},
			want: 5,
		},
		{
			args: args{
				seq:           seq,
				beforeOrEqual: pt.Time(now.Add(5*time.Second - time.Millisecond)),
			},
			want: 4,
		},
		{
			args: args{
				seq:           seq,
				beforeOrEqual: pt.Time(now.Add(5*time.Second + time.Millisecond)),
			},
			want: 5,
		},
		{
			args: args{
				seq:           seq,
				beforeOrEqual: pt.Time(now.Add(100 * time.Second)),
			},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Last(tt.args.seq, tt.args.beforeOrEqual); got != tt.want {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func RandomInt64Sequence(length int) Int64Sequence {
	now := time.Now()
	ret := make(Int64Sequence, length)
	for i := range ret {
		ret[i] = &Int64Item{
			Time:  now.Add(time.Duration(rand.Intn(length)) * time.Second),
			Value: rand.Int63(),
		}
	}
	return ret
}
