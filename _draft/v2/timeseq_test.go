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
			interval: Before(now.Add(time.Hour)),
			args: args{
				t: time.Now(),
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
