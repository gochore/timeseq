package timeseq

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func RandomInts(length int) Ints {
	now := time.Now()
	ret := make(Ints, length)
	for i := range ret {
		delta := time.Duration(i) * time.Second
		if rand.Float64() < 0.5 {
			delta = -delta
		}
		ret[i] = Int{
			Time:  now.Add(delta),
			Value: int(rand.Float64() * float64(math.MaxInt64)),
		}
	}
	return ret
}

func TestIntSeq_Ints(t *testing.T) {
	data := RandomInts(100)
	Sort(data)

	tests := []struct {
		name string
		want Ints
	}{
		{
			name: "regular",
			want: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewIntSeq(data)
			if got := s.Ints(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Index(t *testing.T) {
	data := RandomInts(100)
	Sort(data)

	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want Int
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
			want: Int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewIntSeq(data)
			if got := s.Index(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Time(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	lastYear := now.AddDate(-1, 0, 0)

	data := RandomInts(100)
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
			s := NewIntSeq(data)
			if got := s.Time(tt.args.t); len(got) != tt.length {
				t.Errorf("Time() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestIntSeq_Value(t *testing.T) {
	data := RandomInts(100)
	Sort(data)

	value1 := data[0].Value
	value2 := data[1].Value
	value3 := data[2].Value

	data[0].Value = value1
	data[1].Value = value2
	data[2].Value = value2
	data[3].Value = value2

	type args struct {
		v int
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
			s := NewIntSeq(data)
			if got := s.Value(tt.args.v); len(got) != tt.length {
				t.Errorf("Value() = %v, want %v", got, tt.length)
			}
		})
	}
}

func TestIntSeq_Visit(t *testing.T) {
	data := RandomInts(100)

	type args struct {
		fn func(i int, v Int) (stop bool)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "regular",
			args: args{
				fn: func(i int, v Int) (stop bool) {
					return false
				},
			},
		},
		{
			name: "stop",
			args: args{
				fn: func(i int, v Int) (stop bool) {
					return i > 10
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewIntSeq(data)
			s.Visit(tt.args.fn)
		})
	}
}

func TestIntSeq_Sum(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Count(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.Count(); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Max(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	tests := []struct {
		name   string
		fields fields
		want   Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.Max(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Min(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	tests := []struct {
		name   string
		fields fields
		want   Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.Min(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_First(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	tests := []struct {
		name   string
		fields fields
		want   Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Last(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	tests := []struct {
		name   string
		fields fields
		want   Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Percentile(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	type args struct {
		pct float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.Percentile(tt.args.pct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Range(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	type args struct {
		interval Interval
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *IntSeq
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if got := s.Range(tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSeq_Merge(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	type args struct {
		fn   func(t time.Time, v1, v2 *int) *int
		ints []Ints
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if err := s.Merge(tt.args.fn, tt.args.ints...); (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntSeq_Aggregate(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	type args struct {
		fn       func(t time.Time, ints Ints) *int
		duration time.Duration
		begin    *time.Time
		end      *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if err := s.Aggregate(tt.args.fn, tt.args.duration, tt.args.begin, tt.args.end); (err != nil) != tt.wantErr {
				t.Errorf("Aggregate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntSeq_Trim(t *testing.T) {
	type fields struct {
		slice      Ints
		timeIndex  map[timeKey][]int
		valueIndex map[int][]int
		valueSlice []int
	}
	type args struct {
		fn func(i int, v Int) bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntSeq{
				slice:      tt.fields.slice,
				timeIndex:  tt.fields.timeIndex,
				valueIndex: tt.fields.valueIndex,
				valueSlice: tt.fields.valueSlice,
			}
			if err := s.Trim(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Trim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
