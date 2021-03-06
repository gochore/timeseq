package timeseq

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func parseTime(s string) time.Time {
	ret, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return ret
}

func parseTimeP(s string) *time.Time {
	ret := parseTime(s)
	return &ret
}

func TestInterval_Contain(t *testing.T) {
	now := time.Now()
	type args struct {
		t time.Time
	}
	tests := []struct {
		name     string
		interval Interval
		args     args
		want     bool
	}{
		{
			name:     "regular",
			interval: Interval{}.Before(now.Add(time.Hour)).After(now.Add(-time.Hour)),
			args: args{
				t: now,
			},
			want: true,
		},
		{
			name:     "not contain",
			interval: Interval{}.Before(now.Add(time.Hour)).After(now.Add(-time.Hour)),
			args: args{
				t: now.Add(2 * time.Hour),
			},
			want: false,
		},
		{
			name:     "AfterOrEqual",
			interval: Interval{}.AfterOrEqual(now),
			args: args{
				t: now,
			},
			want: true,
		},
		{
			name:     "After",
			interval: Interval{}.After(now),
			args: args{
				t: now,
			},
			want: false,
		},
		{
			name:     "BeforeOrEqual",
			interval: Interval{}.BeforeOrEqual(now),
			args: args{
				t: now,
			},
			want: true,
		},
		{
			name:     "Before",
			interval: Interval{}.Before(now),
			args: args{
				t: now,
			},
			want: false,
		},
		{
			name:     "regular",
			interval: Before(now.Add(time.Hour)).After(now.Add(-time.Hour)),
			args: args{
				t: now,
			},
			want: true,
		},
		{
			name:     "not contain",
			interval: Before(now.Add(time.Hour)).After(now.Add(-time.Hour)),
			args: args{
				t: now.Add(2 * time.Hour),
			},
			want: false,
		},
		{
			name:     "AfterOrEqual",
			interval: AfterOrEqual(now),
			args: args{
				t: now,
			},
			want: true,
		},
		{
			name:     "After",
			interval: After(now),
			args: args{
				t: now,
			},
			want: false,
		},
		{
			name:     "BeforeOrEqual",
			interval: BeforeOrEqual(now),
			args: args{
				t: now,
			},
			want: true,
		},
		{
			name:     "Before",
			interval: Before(now),
			args: args{
				t: now,
			},
			want: false,
		},
		{
			name:     "chain",
			interval: Before(now.Add(time.Hour)).After(now.Add(-time.Hour)).BeforeOrEqual(now).AfterOrEqual(now).Before(now).After(now),
			args: args{
				t: now,
			},
			want: false,
		},
		{
			name:     "EndAt",
			interval: EndAt(now),
			args: args{
				t: now,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.interval.Contain(tt.args.t); got != tt.want {
				t.Errorf("Contain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_String(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		interval Interval
		want     string
	}{
		{
			name:     "regular",
			interval: Interval{}.BeginAt(now).EndAt(now.Add(time.Hour)),
			want:     fmt.Sprintf("%v~%v", now.Format(time.RFC3339), now.Add(time.Hour).Format(time.RFC3339)),
		},
		{
			name:     "miss begin",
			interval: Interval{}.EndAt(now),
			want:     fmt.Sprintf("nil~%v", now.Format(time.RFC3339)),
		},
		{
			name:     "miss end",
			interval: Interval{}.BeginAt(now),
			want:     fmt.Sprintf("%v~nil", now.Format(time.RFC3339)),
		},
		{
			name:     "miss all",
			interval: Interval{},
			want:     "nil~nil",
		},
		{
			name:     "after",
			interval: Interval{}.After(now),
			want:     fmt.Sprintf("%v~nil", now.Format(time.RFC3339)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.interval.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_Format(t *testing.T) {
	now := time.Now()

	type args struct {
		layout string
	}
	tests := []struct {
		name     string
		interval Interval
		args     args
		want     string
	}{
		{
			name:     "regular",
			interval: Interval{}.BeginAt(now).EndAt(now.Add(time.Hour)),
			args: args{
				layout: time.RFC3339,
			},
			want: fmt.Sprintf("%v~%v", now.Format(time.RFC3339), now.Add(time.Hour).Format(time.RFC3339)),
		},
		{
			name:     "nano",
			interval: Interval{}.After(now),
			args: args{
				layout: time.RFC3339Nano,
			},
			want: fmt.Sprintf("%v~nil", now.Add(1).Format(time.RFC3339Nano)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.interval.Format(tt.args.layout); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_Truncate(t *testing.T) {
	type fields struct {
		NotBefore *time.Time
		NotAfter  *time.Time
	}
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Interval
	}{
		{
			name: "regular",
			fields: fields{
				NotBefore: parseTimeP("2021-02-02T20:34:10+08:00"),
				NotAfter:  parseTimeP("2021-02-02T21:34:10+08:00"),
			},
			args: args{
				d: time.Minute,
			},
			want: Interval{
				NotBefore: parseTimeP("2021-02-02T20:35:00+08:00"),
				NotAfter:  parseTimeP("2021-02-02T21:34:00+08:00"),
			},
		},
		{
			name: "unchanged",
			fields: fields{
				NotBefore: parseTimeP("2021-02-02T20:35:00+08:00"),
				NotAfter:  parseTimeP("2021-02-02T21:34:00+08:00"),
			},
			args: args{
				d: time.Minute,
			},
			want: Interval{
				NotBefore: parseTimeP("2021-02-02T20:35:00+08:00"),
				NotAfter:  parseTimeP("2021-02-02T21:34:00+08:00"),
			},
		},
		{
			name: "nil",
			fields: fields{
				NotBefore: nil,
				NotAfter:  nil,
			},
			args: args{
				d: time.Minute,
			},
			want: Interval{
				NotBefore: nil,
				NotAfter:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				NotBefore: tt.fields.NotBefore,
				NotAfter:  tt.fields.NotAfter,
			}
			if got := i.Truncate(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_Duration(t *testing.T) {
	type fields struct {
		NotBefore *time.Time
		NotAfter  *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "regular",
			fields: fields{
				NotBefore: parseTimeP("2021-02-02T20:35:00+08:00"),
				NotAfter:  parseTimeP("2021-02-02T20:36:00+08:00"),
			},
			want: time.Minute,
		},
		{
			name: "nil",
			fields: fields{
				NotBefore: nil,
				NotAfter:  parseTimeP("2021-02-02T20:36:00+08:00"),
			},
			want: -1,
		},
		{
			name: "zero",
			fields: fields{
				NotBefore: parseTimeP("2021-02-02T20:36:00+08:00"),
				NotAfter:  parseTimeP("2021-02-02T20:36:00+08:00"),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				NotBefore: tt.fields.NotBefore,
				NotAfter:  tt.fields.NotAfter,
			}
			if got := i.Duration(); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
