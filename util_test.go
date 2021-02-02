package timeseq

import (
	"testing"
	"time"
)

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
