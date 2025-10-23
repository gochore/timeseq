package timeseq

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func Test_Merge(t *testing.T) {
	data := randomSeq(10, true)

	type args struct {
		s1, s2 []Point[float64]
		fn     func(t time.Time, v1, v2 *float64) *float64
	}
	tests := []struct {
		name string
		args args
		want []Point[float64]
	}{
		{
			name: "regular",
			args: args{
				s1: data[0:7],
				s2: data[3:10],
				fn: func(t time.Time, v1, v2 *float64) *float64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
			},
			want: data,
		},
		{
			name: "reverse",
			args: args{
				s1: data[3:10],
				s2: data[0:7],
				fn: func(t time.Time, v1, v2 *float64) *float64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
			},
			want: data,
		},
		{
			name: "empty slices",
			args: args{
				s1: data[0:7],
				s2: []Point[float64]{},
				fn: func(t time.Time, v1, v2 *float64) *float64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
			},
			want: data[0:7],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSeq(tt.args.s1)
			s2 := NewSeq(tt.args.s2)
			if got := Merge(s1, s2, tt.args.fn).Points(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", len(got), len(tt.want))
				for i, v := range got {
					t.Errorf("got  %v %v", i, v)
				}
				for i, v := range tt.want {
					t.Errorf("want %v %v", i, v)
				}
				return
			}
		})
	}
}

func TestConvert(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		seq1 := NewSeq(randomSeq(10))
		seq2 := Convert(seq1, func(v Point[float64]) Point[int] {
			return NewPoint(v.Time, int(math.Round(v.Value)))
		})

		points := seq1.Points()
		seq2.Traverse(func(i int, v Point[int]) bool {
			expected := int(math.Round(points[i].Value))
			if v.Time != points[i].Time { // Use != instead of !Time.Equal
				t.Errorf("Convert() time %d = %v, want %v", i, v.Time, points[i].Time)
				return false
			}
			if v.Value != expected {
				t.Errorf("Convert() value %d = %v, want %v", i, v.Value, expected)
				return false
			}
			return true
		})
	})
}
