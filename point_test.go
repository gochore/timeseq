package timeseq

import (
	"testing"
	"time"
)

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
