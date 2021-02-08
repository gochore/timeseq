// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"fmt"
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
		value := rand.Int63()
		for value == 0 || value == 1 || value == 2 { // reserved values
			value = rand.Int63()
		}
		ret[i] = Int64{
			Time:  now.Add(delta),
			Value: value,
		}
	}
	return ret
}

func TestInt64_IsZero(t *testing.T) {
	type fields struct {
		Time  time.Time
		Value int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "zero",
			fields: fields{},
			want:   true,
		},
		{
			name: "not zero",
			fields: fields{
				Time:  time.Now(),
				Value: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Int64{
				Time:  tt.fields.Time,
				Value: tt.fields.Value,
			}
			if got := v.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
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
		name string
		args args
		want Int64
	}{
		{
			name: "regular",
			args: args{
				t: yesterday,
			},
			want: data[3],
		},
		{
			name: "multiple",
			args: args{
				t: lastMonth,
			},
			want: data[0],
		},
		{
			name: "none",
			args: args{
				t: lastYear,
			},
			want: Int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Time(tt.args.t); !got.Equal(tt.want) {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_MTime(t *testing.T) {
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
			if got := s.MTime(tt.args.t); len(got) != tt.length {
				t.Errorf("MTime() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestInt64Seq_Value(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)

	data[0].Value = 0
	data[1].Value = 1
	data[2].Value = 1
	data[3].Value = 1

	type args struct {
		v int64
	}
	tests := []struct {
		name string
		args args
		want Int64
	}{
		{
			name: "regular",
			args: args{
				v: 0,
			},
			want: data[0],
		},
		{
			name: "multiple",
			args: args{
				v: 1,
			},
			want: data[1],
		},
		{
			name: "none",
			args: args{
				v: 2,
			},
			want: Int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Value(tt.args.v); !got.Equal(tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_MValue(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)

	data[0].Value = 0
	data[1].Value = 1
	data[2].Value = 1
	data[3].Value = 1

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
				v: 0,
			},
			length: 1,
		},
		{
			name: "multiple",
			args: args{
				v: 1,
			},
			length: 3,
		},
		{
			name: "none",
			args: args{
				v: 2,
			},
			length: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.MValue(tt.args.v); len(got) != tt.length {
				t.Errorf("MValue() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestInt64Seq_Traverse(t *testing.T) {
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
			s.Traverse(tt.args.fn)
		})
	}
}

func TestInt64Seq_Sum(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)
	var sum int64
	for _, v := range data {
		sum += v.Value
	}
	tests := []struct {
		name string
		want int64
	}{
		{
			name: "regular",
			want: sum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Max(t *testing.T) {
	data := RandomInt64s(100)
	max := data[0]
	for _, v := range data {
		if v.Value > max.Value {
			max = v
		}
	}

	tests := []struct {
		name string
		want Int64
	}{
		{
			name: "regular",
			want: max,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Max(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Min(t *testing.T) {
	data := RandomInt64s(100)
	min := data[0]
	for _, v := range data {
		if v.Value < min.Value {
			min = v
		}
	}

	tests := []struct {
		name string
		want Int64
	}{
		{
			name: "regular",
			want: min,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Min(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_First(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)

	tests := []struct {
		name string
		data Int64s
		want Int64
	}{
		{

			name: "regular",
			data: data,
			want: data[0],
		},
		{
			name: "emtpy",
			data: nil,
			want: Int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(tt.data)
			if got := s.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Last(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)
	tests := []struct {
		name string
		data Int64s
		want Int64
	}{
		{

			name: "regular",
			data: data,
			want: data[len(data)-1],
		},
		{
			name: "emtpy",
			data: nil,
			want: Int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(tt.data)
			if got := s.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Percentile(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)
	for i := range data {
		data[i].Value = int64(i)
	}

	type args struct {
		pct float64
	}
	tests := []struct {
		name string
		data Int64s
		args args
		want Int64
	}{
		{
			name: "data[0]",
			data: data,
			args: args{
				pct: 0,
			},
			want: data[0],
		},
		{
			name: "data[49]",
			data: data,
			args: args{
				pct: 0.5,
			},
			want: data[49],
		},
		{
			name: "0.95",
			data: data,
			args: args{
				pct: 0.95,
			},
			want: data[94],
		},
		{
			name: "0.955",
			data: data,
			args: args{
				pct: 0.955,
			},
			want: data[94],
		},
		{
			name: "1",
			data: data,
			args: args{
				pct: 1,
			},
			want: data[99],
		},
		{
			name: "1.1",
			data: data,
			args: args{
				pct: 1.1,
			},
			want: data[99],
		},
		{
			name: "-0.1",
			data: data,
			args: args{
				pct: -0.1,
			},
			want: data[0],
		},
		{
			name: "empty",
			data: nil,
			args: args{
				pct: 1,
			},
			want: Int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(tt.data)
			if got := s.Percentile(tt.args.pct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Range(t *testing.T) {
	data := RandomInt64s(100)
	Sort(data)

	type args struct {
		interval Interval
	}
	tests := []struct {
		name string
		data Int64s
		args args
		want Int64s
	}{
		{
			name: "regular",
			data: data,
			args: args{
				interval: AfterOrEqual(data[10].Time).BeforeOrEqual(data[89].Time),
			},
			want: data[10:90],
		},
		{
			name: "nil NotBefore",
			data: data,
			args: args{
				interval: BeforeOrEqual(data[89].Time),
			},
			want: data[:90],
		},
		{
			name: "nil NotAfter",
			data: data,
			args: args{
				interval: AfterOrEqual(data[10].Time),
			},
			want: data[10:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Range(tt.args.interval).Int64s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Merge(t *testing.T) {
	data := RandomInt64s(10)
	Sort(data)

	type args struct {
		fn    func(t time.Time, v1, v2 *int64) *int64
		slice Int64s
	}
	tests := []struct {
		name string
		data Int64s
		args args
		want Int64s
	}{
		{
			name: "regular",
			data: data[0:7],
			args: args{
				fn: func(t time.Time, v1, v2 *int64) *int64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slice: data[3:10],
			},
			want: data,
		},
		{
			name: "reverse",
			data: data[3:10],
			args: args{
				fn: func(t time.Time, v1, v2 *int64) *int64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slice: data[0:7],
			},
			want: data,
		},
		{
			name: "nil fn",
			data: data[3:10],
			args: args{
				fn: nil,
			},
			want: data[3:10],
		},
		{
			name: "empty slices",
			data: data[0:7],
			args: args{
				fn: func(t time.Time, v1, v2 *int64) *int64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slice: Int64s{},
			},
			want: data[0:7],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(tt.data)
			if got := s.Merge(tt.args.fn, WrapInt64Seq(tt.args.slice)).Int64s(); !reflect.DeepEqual(got, tt.want) {
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

func TestInt64Seq_Aggregate(t *testing.T) {
	type args struct {
		fn       func(t time.Time, slice Int64s) *int64
		duration time.Duration
		interval Interval
	}
	tests := []struct {
		name  string
		slice Int64s
		args  args
		want  Int64s
	}{
		{
			name: "regular",
			slice: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:19:01+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
			args: args{
				fn: func(t time.Time, slice Int64s) *int64 {
					var ret int64
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: time.Hour,
				interval: BeginAt(parseTime("2021-02-08T17:00:00+08:00")).EndAt(parseTime("2021-02-08T23:59:59+08:00")),
			},
			want: Int64s{
				{Time: parseTime("2021-02-08T17:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T18:00:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:00:00+08:00"), Value: 5},
				{Time: parseTime("2021-02-08T20:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T21:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T22:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T23:00:00+08:00"), Value: 0},
			},
		},
		{
			name: "nil fn",
			slice: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:19:01+08:00"), Value: 3},
			},
			args: args{
				fn:       nil,
				duration: time.Hour,
				interval: BeginAt(parseTime("2021-02-08T17:00:00+08:00")).EndAt(parseTime("2021-02-08T23:59:59+08:00")),
			},
			want: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:19:01+08:00"), Value: 3},
			},
		},
		{
			name: "zero duration",
			slice: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
			args: args{
				fn: func(t time.Time, slice Int64s) *int64 {
					var ret int64
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: 0,
				interval: BeginAt(parseTime("2021-02-08T17:00:00+08:00")).EndAt(parseTime("2021-02-08T23:59:59+08:00")),
			},
			want: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
		},
		{
			name: "shorter interval with zero duraion",
			slice: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
			args: args{
				fn: func(t time.Time, slice Int64s) *int64 {
					var ret int64
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: 0,
				interval: BeginAt(parseTime("2021-02-08T18:30:00+08:00")).EndAt(parseTime("2021-02-08T23:59:59+08:00")),
			},
			want: Int64s{
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
		},
		{
			name: "shorter interval with non zero duraion",
			slice: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
			args: args{
				fn: func(t time.Time, slice Int64s) *int64 {
					var ret int64
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: time.Hour,
				interval: BeginAt(parseTime("2021-02-08T18:30:00+08:00")).EndAt(parseTime("2021-02-08T23:59:59+08:00")),
			},
			want: Int64s{
				{Time: parseTime("2021-02-08T19:00:00+08:00"), Value: 5},
				{Time: parseTime("2021-02-08T20:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T21:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T22:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T23:00:00+08:00"), Value: 0},
			},
		},
		{
			name: "miss begin",
			slice: Int64s{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
			args: args{
				fn: func(t time.Time, slice Int64s) *int64 {
					var ret int64
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: time.Hour,
				interval: EndAt(parseTime("2021-02-08T23:59:59+08:00")),
			},
			want: Int64s{
				{Time: parseTime("2021-02-08T18:00:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:00:00+08:00"), Value: 5},
				{Time: parseTime("2021-02-08T20:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T21:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T22:00:00+08:00"), Value: 0},
				{Time: parseTime("2021-02-08T23:00:00+08:00"), Value: 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(tt.slice)
			if got := s.Aggregate(tt.args.fn, tt.args.duration, tt.args.interval).Int64s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregate() got = %v, want %v", len(got), len(tt.want))
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

func TestInt64Seq_Trim(t *testing.T) {
	data := RandomInt64s(10)
	Sort(data)

	type args struct {
		fn func(i int, v Int64) bool
	}
	tests := []struct {
		name string
		args args
		want Int64s
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Int64) bool {
					return i >= 5
				},
			},
			want: data[:5],
		},
		{
			name: "nil fn",
			args: args{
				fn: nil,
			},
			want: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Trim(tt.args.fn).Int64s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
				for i, v := range got {
					fmt.Println(i, v)
				}
				for i, v := range tt.want {
					fmt.Println(i, v)
				}
			}
		})
	}
}
