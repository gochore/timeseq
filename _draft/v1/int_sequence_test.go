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

func TestIntSequence_Len(t *testing.T) {
	tests := []struct {
		name string
		s    IntSequence
		want int
	}{
		{
			s:    RandomIntSequence(10),
			want: 10,
		},
		{
			s:    RandomIntSequence(0),
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
				t.Errorf("IntSequence.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    IntSequence
		args args
	}{
		{
			s: RandomIntSequence(10),
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
				t.Errorf("IntSequence.Swap() failed")
			}
		})
	}
}

func TestIntSequence_Time(t *testing.T) {
	seq := RandomIntSequence(10)
	seq.Sort()

	type args struct {
		i int
	}
	tests := []struct {
		name string
		s    IntSequence
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
				t.Errorf("IntSequence.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Slice(t *testing.T) {
	seq := RandomIntSequence(10)

	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    IntSequence
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
				t.Errorf("IntSequence.Interface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Sort(t *testing.T) {
	tests := []struct {
		name string
		s    IntSequence
	}{
		{
			s: RandomIntSequence(10),
		},
		{
			s: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Sort()
			if !sort.IsSorted(sortableSequence{tt.s}) {
				t.Error("IntSequence.Interface() failed")
			}
		})
	}
}

func TestIntSequence_Range(t *testing.T) {
	now := time.Now()
	seq := RandomIntSequence(10)
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
		s    IntSequence
		args args
		want IntSequence
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
				t.Errorf("IntSequence.Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_First(t *testing.T) {
	now := time.Now()
	seq := RandomIntSequence(10)
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
		s    IntSequence
		args args
		want *IntItem
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
				t.Errorf("IntSequence.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Last(t *testing.T) {
	now := time.Now()
	seq := RandomIntSequence(10)
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
		s    IntSequence
		args args
		want *IntItem
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
				t.Errorf("IntSequence.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Max(t *testing.T) {
	seq1 := RandomIntSequence(10)
	seq1.Sort()
	for i := range seq1 {
		if i == 0 {
			seq1[i].Value = 1
		} else {
			seq1[i].Value = 0
		}
	}
	seq2 := RandomIntSequence(10)
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
		s    IntSequence
		want *IntItem
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
				t.Errorf("IntSequence.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Min(t *testing.T) {
	seq1 := RandomIntSequence(10)
	seq1.Sort()
	for i := range seq1 {
		if i == 0 {
			seq1[i].Value = 1
		} else {
			seq1[i].Value = 0
		}
	}
	seq2 := RandomIntSequence(10)
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
		s    IntSequence
		want *IntItem
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
				t.Errorf("IntSequence.Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Sum(t *testing.T) {
	seq := RandomIntSequence(10)
	for i := range seq {
		seq[i].Value = int(i)
	}

	tests := []struct {
		name string
		s    IntSequence
		want int
	}{
		{
			s:    seq,
			want: 45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Sum(); got != tt.want {
				t.Errorf("IntSequence.Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Average(t *testing.T) {
	seq := RandomIntSequence(10)
	for i := range seq {
		seq[i].Value = int(i) * 2
	}

	tests := []struct {
		name string
		s    IntSequence
		want int
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
				t.Errorf("IntSequence.Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSequence_Percentile(t *testing.T) {
	seq := RandomIntSequence(100)
	for i := range seq {
		seq[i].Value = int(i) + 1
	}

	type args struct {
		pct float64
	}
	tests := []struct {
		name string
		s    IntSequence
		args args
		want int
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
					t.Errorf("IntSequence.Percentile() failed: %v", r)
				}
			}()
			if got := tt.s.Percentile(tt.args.pct); got != tt.want {
				t.Errorf("IntSequence.Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeInt(t *testing.T) {
	seq := RandomIntSequence(10)
	seq.Sort()
	type args struct {
		seq1 IntSequence
		seq2 IntSequence
		fn   func(item1, item2 *IntItem) *IntItem
	}
	tests := []struct {
		name string
		args args
		want IntSequence
	}{
		{
			name: "regular",
			args: args{
				seq1: seq[0:7],
				seq2: seq[3:10],
				fn: func(item1, item2 *IntItem) *IntItem {
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
				fn: func(item1, item2 *IntItem) *IntItem {
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
			if got := MergeInt(tt.args.seq1, tt.args.seq2, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func RandomIntSequence(length int) IntSequence {
	now := time.Now()
	ret := make(IntSequence, length)
	for i := range ret {
		delta := time.Duration(i) * time.Second
		if rand.Float64() < 0.5 {
			delta = -delta
		}
		ret[i] = IntItem{
			Time:  now.Add(delta),
			Value: int(rand.Float64() * float64(math.MaxInt64)),
		}
	}
	return ret
}