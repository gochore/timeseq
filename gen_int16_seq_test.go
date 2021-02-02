// Code generated by cmd/generate. DO NOT EDIT.

package timeseq

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func RandomInt16s(length int) Int16s {
	now := time.Now()
	ret := make(Int16s, length)
	for i := range ret {
		delta := time.Duration(i) * time.Second
		if rand.Float64() < 0.5 {
			delta = -delta
		}
		value := int16(rand.Int31())
		for value == 0 || value == 1 || value == 2 { // reserved values
			value = int16(rand.Int31())
		}
		ret[i] = Int16{
			Time:  now.Add(delta),
			Value: value,
		}
	}
	return ret
}

func TestInt16_IsZero(t *testing.T) {
	type fields struct {
		Time  time.Time
		Value int16
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
			v := Int16{
				Time:  tt.fields.Time,
				Value: tt.fields.Value,
			}
			if got := v.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Int16s(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)

	tests := []struct {
		name string
		want Int16s
	}{
		{
			name: "regular",
			want: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if got := s.Int16s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Index(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)

	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want Int16
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
			want: Int16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if got := s.Index(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Time(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	lastYear := now.AddDate(-1, 0, 0)

	data := RandomInt16s(100)
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
		want Int16
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
			want: Int16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if got := s.Time(tt.args.t); !got.Equal(tt.want) {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_MTime(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	lastYear := now.AddDate(-1, 0, 0)

	data := RandomInt16s(100)
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
			s := NewInt16Seq(data)
			if got := s.MTime(tt.args.t); len(got) != tt.length {
				t.Errorf("MTime() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestInt16Seq_Value(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)

	data[0].Value = 0
	data[1].Value = 1
	data[2].Value = 1
	data[3].Value = 1

	type args struct {
		v int16
	}
	tests := []struct {
		name string
		args args
		want Int16
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
			want: Int16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if got := s.Value(tt.args.v); !got.Equal(tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_MValue(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)

	data[0].Value = 0
	data[1].Value = 1
	data[2].Value = 1
	data[3].Value = 1

	type args struct {
		v int16
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
			s := NewInt16Seq(data)
			if got := s.MValue(tt.args.v); len(got) != tt.length {
				t.Errorf("MValue() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestInt16Seq_Traverse(t *testing.T) {
	data := RandomInt16s(100)

	type args struct {
		fn func(i int, v Int16) (stop bool)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Int16) (stop bool) {
					return false
				},
			},
		},
		{
			name: "stop",
			args: args{
				fn: func(i int, v Int16) (stop bool) {
					return i > 10
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			s.Traverse(tt.args.fn)
		})
	}
}

func TestInt16Seq_Sum(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)
	var sum int16
	for _, v := range data {
		sum += v.Value
	}
	tests := []struct {
		name string
		want int16
	}{
		{
			name: "regular",
			want: sum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Max(t *testing.T) {
	data := RandomInt16s(100)
	max := data[0]
	for _, v := range data {
		if v.Value > max.Value {
			max = v
		}
	}

	tests := []struct {
		name string
		want Int16
	}{
		{
			name: "regular",
			want: max,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if got := s.Max(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Min(t *testing.T) {
	data := RandomInt16s(100)
	min := data[0]
	for _, v := range data {
		if v.Value < min.Value {
			min = v
		}
	}

	tests := []struct {
		name string
		want Int16
	}{
		{
			name: "regular",
			want: min,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if got := s.Min(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_First(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)

	tests := []struct {
		name string
		data Int16s
		want Int16
	}{
		{

			name: "regular",
			data: data,
			want: data[0],
		},
		{
			name: "emtpy",
			data: nil,
			want: Int16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(tt.data)
			if got := s.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Last(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)
	tests := []struct {
		name string
		data Int16s
		want Int16
	}{
		{

			name: "regular",
			data: data,
			want: data[len(data)-1],
		},
		{
			name: "emtpy",
			data: nil,
			want: Int16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(tt.data)
			if got := s.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Percentile(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)
	for i := range data {
		data[i].Value = int16(i)
	}

	type args struct {
		pct float64
	}
	tests := []struct {
		name string
		data Int16s
		args args
		want Int16
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
			want: Int16{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(tt.data)
			if got := s.Percentile(tt.args.pct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Range(t *testing.T) {
	data := RandomInt16s(100)
	Sort(data)

	type args struct {
		interval Interval
	}
	tests := []struct {
		name string
		data Int16s
		args args
		want Int16s
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
			s := NewInt16Seq(data)
			if got := s.Range(tt.args.interval).slice; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16Seq_Merge(t *testing.T) {
	data := RandomInt16s(10)
	Sort(data)

	type args struct {
		fn     func(t time.Time, v1, v2 *int16) *int16
		slices []Int16s
	}
	tests := []struct {
		name    string
		data    Int16s
		args    args
		want    Int16s
		wantErr bool
	}{
		{
			name: "regular",
			data: data[0:7],
			args: args{
				fn: func(t time.Time, v1, v2 *int16) *int16 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slices: []Int16s{data[3:10]},
			},
			want: data,
		},
		{
			name: "reverse",
			data: data[3:10],
			args: args{
				fn: func(t time.Time, v1, v2 *int16) *int16 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slices: []Int16s{data[0:7]},
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
				fn: func(t time.Time, v1, v2 *int16) *int16 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slices: []Int16s{
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
				fn: func(t time.Time, v1, v2 *int16) *int16 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slices: []Int16s{
					append(Int16s{data[9]}, data[3:9]...),
				},
			},
			want: data,
		},
		{
			name: "empty slices",
			data: data[0:7],
			args: args{
				fn: func(t time.Time, v1, v2 *int16) *int16 {
					if v1 != nil {
						return v1
					}
					return v2
				},
				slices: []Int16s{},
			},
			want: data[0:7],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(tt.data)
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

func TestInt16Seq_Aggregate(t *testing.T) {
	data := RandomInt16s(100)
	now := time.Now()
	begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := begin.Add(24*time.Hour - 1)

	type args struct {
		fn       func(t time.Time, slice Int16s) *int16
		duration time.Duration
		interval Interval
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "regular",
			args: args{
				fn: func(t time.Time, slice Int16s) *int16 {
					ret := int16(t.Hour())
					if len(slice) != 0 {
						ret = 0
					}
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: time.Hour,
				interval: BeginAt(begin).EndAt(end),
			},
			wantErr: false,
		},
		{
			name: "nil fn",
			args: args{
				fn:       nil,
				duration: time.Hour,
				interval: BeginAt(begin).EndAt(end),
			},
			wantErr: true,
		},
		{
			name: "zero duration",
			args: args{
				fn: func(t time.Time, slice Int16s) *int16 {
					ret := int16(t.Hour())
					if len(slice) != 0 {
						ret = 0
					}
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: 0,
				interval: BeginAt(begin).EndAt(end),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInt16Seq(data)
			if err := s.Aggregate(tt.args.fn, tt.args.duration, tt.args.interval); (err != nil) != tt.wantErr {
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

func TestInt16Seq_Trim(t *testing.T) {
	data := RandomInt16s(10)
	Sort(data)

	type args struct {
		fn func(i int, v Int16) bool
	}
	tests := []struct {
		name    string
		args    args
		want    Int16s
		wantErr bool
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Int16) bool {
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
			s := NewInt16Seq(data)
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

func TestInt16Seq_Clone(t *testing.T) {
	data := RandomInt16s(10)
	Sort(data)
	tests := []struct {
		name string
		seq  *Int16Seq
		want *Int16Seq
	}{
		{
			name: "regular",
			seq:  NewInt16Seq(data),
			want: &Int16Seq{
				slice: data,
			},
		},
		{
			name: "nil",
			seq:  nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.seq.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
