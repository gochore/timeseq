package timeseq

import (
	"reflect"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	type args struct {
		seq Sequence
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sort(tt.args.seq)
		})
	}
}

func TestRange(t *testing.T) {
	type args struct {
		seq           Sequence
		afterOrEqual  time.Time
		beforeOrEqual time.Time
	}
	tests := []struct {
		name string
		args args
		want Sequence
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Range(tt.args.seq, tt.args.afterOrEqual, tt.args.beforeOrEqual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	type args struct {
		seq          Sequence
		afterOrEqual time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := First(tt.args.seq, tt.args.afterOrEqual); got != tt.want {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLast(t *testing.T) {
	type args struct {
		seq           Sequence
		beforeOrEqual time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Last(tt.args.seq, tt.args.beforeOrEqual); got != tt.want {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}
