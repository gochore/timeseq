package timeseq

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

func randomSeq(length int, sorted ...bool) []Point[float64] {
	needSorted := false
	if len(sorted) > 0 {
		needSorted = sorted[0]
	}

	now := time.Now()
	ret := make([]Point[float64], length)
	for i := range ret {
		delta := time.Duration(i) * time.Second
		if !needSorted {
			if rand.Float64() < 0.5 {
				delta = -delta
			}
		}
		value := rand.Float64()
		for value == 0 || value == 1 || value == 2 { // reserved values
			value = rand.Float64()
		}
		ret[i] = NewPoint(now.Add(delta), value)
	}
	return ret
}

func TestPoint_IsZero(t *testing.T) {
	type fields struct {
		Time  time.Time
		Value float64
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
			v := NewPoint(tt.fields.Time, tt.fields.Value)
			if got := v.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Points(t *testing.T) {
	data := randomSeq(100)
	sortedData := make([]Point[float64], len(data))
	copy(sortedData, data)
	sort.SliceStable(sortedData, func(i, j int) bool {
		return sortedData[i].Time.Before(sortedData[j].Time)
	})

	tests := []struct {
		name string
		want []Point[float64]
	}{
		{
			name: "regular",
			want: sortedData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Points(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Points() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Index(t *testing.T) {
	data := randomSeq(100, true)

	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want Point[float64]
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
			want: Point[float64]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Index(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Time(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	lastYear := now.AddDate(-1, 0, 0)

	data := randomSeq(100, true)
	data[0].Time = lastMonth
	data[1].Time = lastMonth
	data[2].Time = lastMonth
	data[3].Time = yesterday

	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want Point[float64]
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
			want: Point[float64]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Time(tt.args.t); !got.Equal(tt.want) {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_MTime(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	lastYear := now.AddDate(-1, 0, 0)

	data := randomSeq(100, true)
	data[0].Time = lastMonth
	data[1].Time = lastMonth
	data[2].Time = lastMonth
	data[3].Time = yesterday

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
			s := NewSeq(data)
			if got := s.MTime(tt.args.t); len(got) != tt.length {
				t.Errorf("MTime() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestSeq_Value(t *testing.T) {
	data := randomSeq(100, true)

	data[0].Value = 0
	data[1].Value = 1
	data[2].Value = 1
	data[3].Value = 1

	type args struct {
		v float64
	}
	tests := []struct {
		name string
		args args
		want Point[float64]
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
			want: Point[float64]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Value(tt.args.v); !got.Equal(tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_MValue(t *testing.T) {
	data := randomSeq(100, true)

	data[0].Value = 0
	data[1].Value = 1
	data[2].Value = 1
	data[3].Value = 1

	type args struct {
		v float64
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
			s := NewSeq(data)
			if got := s.MValue(tt.args.v); len(got) != tt.length {
				t.Errorf("MValue() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestSeq_Traverse(t *testing.T) {
	data := randomSeq(100)

	type args struct {
		fn func(i int, v Point[float64]) (stop bool)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Point[float64]) (stop bool) {
					return false
				},
			},
		},
		{
			name: "stop",
			args: args{
				fn: func(i int, v Point[float64]) (stop bool) {
					return i > 10
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			s.Traverse(tt.args.fn)
		})
	}
}

func TestSeq_Sum(t *testing.T) {
	data := randomSeq(100, true)

	var sum float64
	for _, v := range data {
		sum += v.Value
	}
	tests := []struct {
		name string
		want float64
	}{
		{
			name: "regular",
			want: sum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Max(t *testing.T) {
	data := randomSeq(100)

	m := data[0]
	for _, v := range data {
		if v.Value > m.Value {
			m = v
		}
	}

	tests := []struct {
		name string
		want Point[float64]
	}{
		{
			name: "regular",
			want: m,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Max(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Min(t *testing.T) {
	data := randomSeq(100)

	m := data[0]
	for _, v := range data {
		if v.Value < m.Value {
			m = v
		}
	}

	tests := []struct {
		name string
		want Point[float64]
	}{
		{
			name: "regular",
			want: m,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Min(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_First(t *testing.T) {
	data := randomSeq(100, true)

	tests := []struct {
		name string
		data []Point[float64]
		want Point[float64]
	}{
		{

			name: "regular",
			data: data,
			want: data[0],
		},
		{
			name: "emtpy",
			data: nil,
			want: Point[float64]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(tt.data)
			if got := s.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Last(t *testing.T) {
	data := randomSeq(100, true)

	tests := []struct {
		name string
		data []Point[float64]
		want Point[float64]
	}{
		{

			name: "regular",
			data: data,
			want: data[len(data)-1],
		},
		{
			name: "emtpy",
			data: nil,
			want: Point[float64]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(tt.data)
			if got := s.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Percentile(t *testing.T) {
	data := randomSeq(100, true)

	for i := range data {
		data[i].Value = float64(i)
	}

	type args struct {
		pct float64
	}
	tests := []struct {
		name string
		data []Point[float64]
		args args
		want Point[float64]
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
			want: Point[float64]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(tt.data)
			if got := s.Percentile(tt.args.pct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Range(t *testing.T) {
	data := randomSeq(100, true)

	type args struct {
		begin time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		data []Point[float64]
		args args
		want []Point[float64]
	}{
		{
			name: "regular",
			data: data,
			args: args{
				begin: data[10].Time,
				end:   data[89].Time,
			},
			want: data[10:90],
		},
		{
			name: "nil NotBefore",
			data: data,
			args: args{
				end: data[89].Time,
			},
			want: data[:90],
		},
		{
			name: "nil NotAfter",
			data: data,
			args: args{
				begin: data[10].Time,
			},
			want: data[10:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Range(tt.args.begin, tt.args.end).Points(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Slice(t *testing.T) {
	data := randomSeq(100, true)

	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		data []Point[float64]
		args args
		want []Point[float64]
	}{
		{
			name: "regular",
			data: data,
			args: args{
				i: 1,
				j: 3,
			},
			want: data[1:3],
		},
		{
			name: "left negative",
			data: data,
			args: args{
				i: -1,
				j: 3,
			},
			want: data[:3],
		},
		{
			name: "right negative",
			data: data,
			args: args{
				i: 1,
				j: -3,
			},
			want: data[1:],
		},
		{
			name: "all negative",
			data: data,
			args: args{
				i: -1,
				j: -3,
			},
			want: data[:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Slice(tt.args.i, tt.args.j).Points(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_Filter(t *testing.T) {
	data := randomSeq(10, true)

	type args struct {
		fn func(i int, v Point[float64]) bool
	}
	tests := []struct {
		name string
		args args
		want []Point[float64]
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Point[float64]) bool {
					return i < 5
				},
			},
			want: data[:5],
		},
		{
			name: "all keep",
			args: args{
				fn: func(i int, v Point[float64]) bool {
					return true
				},
			},
			want: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(data)
			if got := s.Filter(tt.args.fn).Points(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
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

func parseTime(s string) time.Time {
	ret, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return ret
}

func TestSeq_Aggregate(t *testing.T) {
	type args struct {
		fn       func(t time.Time, slice []Point[float64]) *float64
		duration time.Duration
	}
	tests := []struct {
		name  string
		slice []Point[float64]
		args  args
		want  []Point[float64]
	}{
		{
			name: "regular",
			slice: []Point[float64]{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:19:01+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
			args: args{
				fn: func(t time.Time, slice []Point[float64]) *float64 {
					var ret float64
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: time.Hour,
			},
			want: []Point[float64]{
				{Time: parseTime("2021-02-08T18:00:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:00:00+08:00"), Value: 5},
			},
		},
		{
			name: "zero duration",
			slice: []Point[float64]{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
			args: args{
				fn: func(t time.Time, slice []Point[float64]) *float64 {
					var ret float64
					for _, v := range slice {
						ret += v.Value
					}
					return &ret
				},
				duration: 0,
			},
			want: []Point[float64]{
				{Time: parseTime("2021-02-08T18:08:00+08:00"), Value: 1},
				{Time: parseTime("2021-02-08T19:08:04+08:00"), Value: 2},
				{Time: parseTime("2021-02-08T19:09:01+08:00"), Value: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSeq(tt.slice)
			if got := s.Aggregate(tt.args.fn, tt.args.duration).Points(); !reflect.DeepEqual(got, tt.want) {
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

func TestConvertSeq(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		type P struct {
			Timestamp int64
			Value     int64
		}
		ps := []P{
			{Timestamp: 1, Value: 10},
			{Timestamp: 2, Value: 20},
			{Timestamp: 3, Value: 30},
		}
		seq := ConvertSeq(ps, func(p P) Point[int64] {
			return NewPoint(time.Unix(p.Timestamp, 0), p.Value)
		})
		want := []Point[int64]{
			{Time: time.Unix(1, 0), Value: 10},
			{Time: time.Unix(2, 0), Value: 20},
			{Time: time.Unix(3, 0), Value: 30},
		}
		if !reflect.DeepEqual(seq.Points(), want) {
			t.Errorf("ConvertSeq() = %v, want %v", seq.Points(), want)
		}
	})
}
