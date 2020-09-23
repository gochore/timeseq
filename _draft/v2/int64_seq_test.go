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

func TestInt64Seq_Count(t *testing.T) {
	data := RandomInt64s(100)

	tests := []struct {
		name string
		want int
	}{
		{
			name: "regular",
			want: len(data),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Count(); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
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
				interval: Interval{
					NotBefore: &data[10].Time,
					NotAfter:  &data[89].Time,
				},
			},
			want: data[10:90],
		},
		{
			name: "nil NotBefore",
			data: data,
			args: args{
				interval: Interval{
					NotAfter: &data[89].Time,
				},
			},
			want: data[:90],
		},
		{
			name: "nil NotAfter",
			data: data,
			args: args{
				interval: Interval{
					NotBefore: &data[10].Time,
				},
			},
			want: data[10:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if got := s.Range(tt.args.interval).slice; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Seq_Merge(t *testing.T) {
	data := RandomInt64s(10)
	Sort(data)

	type args struct {
		fn     func(t time.Time, v1, v2 *int64) *int64
		slices []Int64s
	}
	tests := []struct {
		name    string
		data    Int64s
		args    args
		want    Int64s
		wantErr bool
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
				slices: []Int64s{data[3:10]},
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
				slices: []Int64s{data[0:7]},
			},
			want: data,
		},
		{
			name: "nil fn",
			data: nil,
			args: args{
				fn: nil,
			},
			wantErr: true,
		},
		{
			name: "multiple",
			data: nil,
			args: args{
				fn: func(t time.Time, v1, v2 *int64) *int64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slices: []Int64s{
					data[1:2],
					data[0:4],
					nil,
					data[2:9],
					data[9:],
				},
			},
			want: data,
		},
		{
			name: "not sorted",
			data: data[0:7],
			args: args{
				fn: func(t time.Time, v1, v2 *int64) *int64 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slices: []Int64s{
					append(Int64s{data[9]}, data[3:9]...),
				},
			},
			want: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(tt.data)
			if err := s.Merge(tt.args.fn, tt.args.slices...); (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got := s.slice; !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Merge() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestInt64Seq_Aggregate(t *testing.T) {
	data := RandomInt64s(100)
	now := time.Now()
	begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := begin.Add(24*time.Hour - 1)

	type args struct {
		fn       func(t time.Time, slice Int64s) *int64
		duration time.Duration
		begin    *time.Time
		end      *time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "regular",
			args: args{
				fn: func(t time.Time, slice Int64s) *int64 {
					ret := int64(t.Hour())
					if len(slice) != 0 {
						ret = 0
					}
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: time.Hour,
				begin:    &begin,
				end:      &end,
			},
			wantErr: false,
		},
		{
			name: "nil fn",
			args: args{
				fn:       nil,
				duration: time.Hour,
				begin:    &begin,
				end:      &end,
			},
			wantErr: true,
		},
		{
			name: "zero duration",
			args: args{
				fn: func(t time.Time, slice Int64s) *int64 {
					ret := int64(t.Hour())
					if len(slice) != 0 {
						ret = 0
					}
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: 0,
				begin:    &begin,
				end:      &end,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if err := s.Aggregate(tt.args.fn, tt.args.duration, tt.args.begin, tt.args.end); (err != nil) != tt.wantErr {
				t.Errorf("Aggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for i, v := range s.slice {
					fmt.Println(i, v.Value, v.Time)
				}
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
		name    string
		args    args
		want    Int64s
		wantErr bool
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Int64) bool {
					return i >= 5
				},
			},
			want:    data[:5],
			wantErr: false,
		},
		{
			name: "nil fn",
			args: args{
				fn: nil,
			},
			want:    data[:5],
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt64Seq(data)
			if err := s.Trim(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Trim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got := s.slice; !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Trim() = %v, want %v", got, tt.want)
					for i, v := range got {
						fmt.Println(i, v)
					}
					for i, v := range tt.want {
						fmt.Println(i, v)
					}
				}
			}
		})
	}
}
