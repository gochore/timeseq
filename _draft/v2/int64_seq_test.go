package timeseq

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func RandomInt64s(length int) Int64s {
	now := time.Now()
	ret := make(Int64s, length)
	for i := range ret {
		delta := time.Duration(i) * time.Second
		if rand.Float64() < 0.5 {
			delta = -delta
		}
		ret[i] = Int64{
			Time:  now.Add(delta),
			Value: rand.Int63(),
		}
	}
	return ret
}

func TestInt64Seq_Int64s(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)

	tests := []struct {
		name string
		want Int64s
	}{
		{
			name: "regular",
			want: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Int64s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Index(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)

	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want Int64
	}{
		{
			name: "regular",
			args: args{
				i: 1,
			},
			want: data[1],
		},
		{
			name: "less than zero",
			args: args{
				i: -1,
			},
			want: Int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Index(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Time(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	lastYear := now.AddDate(-1, 0, 0)

	data := RandomInt64s(100)
	data[0].Time = lastMonth
	data[1].Time = lastMonth
	data[2].Time = lastMonth
	data[3].Time = yesterday
	Sort(data)

	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		args   args
		length int
	}{
		{
			name: "regular",
			args: args{
				t: yesterday,
			},
			length: 1,
		},
		{
			name: "multiple",
			args: args{
				t: lastMonth,
			},
			length: 3,
		},
		{
			name: "none",
			args: args{
				t: lastYear,
			},
			length: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Time(tt.args.t); len(got) != tt.length {
				t.Errorf("Time() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestInt64Seq_Value(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)

	value1 := data[0].Value
	value2 := data[1].Value
	value3 := data[2].Value

	data[0].Value = value1
	data[1].Value = value2
	data[2].Value = value2
	data[3].Value = value2

	type args struct {
		v int64
	}
	tests := []struct {
		name   string
		args   args
		length int
	}{
		{
			name: "regular",
			args: args{
				v: value1,
			},
			length: 1,
		},
		{
			name: "multiple",
			args: args{
				v: value2,
			},
			length: 3,
		},
		{
			name: "none",
			args: args{
				v: value3,
			},
			length: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Value(tt.args.v); len(got) != tt.length {
				t.Errorf("Value() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestInt64Seq_Visit(t *testing.T) {
	data := RandomInt64s(100)

	type args struct {
		fn func(i int, v Int64) (stop bool)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Int64) (stop bool) {
					return false
				},
			},
		},
		{
			name: "stop",
			args: args{
				fn: func(i int, v Int64) (stop bool) {
					return i > 10
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			s.Visit(tt.args.fn)
		})
	}
}
