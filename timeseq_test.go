package timeseq

import (
	"testing"
	"time"
)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.interval.Contain(tt.args.t); got != tt.want {
				t.Errorf("Contain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeKey_Time(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		k    timeKey
		want time.Time
	}{
		{
			name: "regular",
			k:    newTimeKey(now),
			want: now,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.Time(); !got.Equal(tt.want) {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}
