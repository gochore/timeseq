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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.interval.Contain(tt.args.t); got != tt.want {
				t.Errorf("Contain() = %v, want %v", got, tt.want)
			}
		})
	}
}
