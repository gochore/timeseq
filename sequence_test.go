package timeseq

import (
	"testing"
)

func Test_sortableSequence_Less(t *testing.T) {
	seq := RandomInt64Sequence(10)

	type fields struct {
		Sequence Sequence
	}
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			fields: fields{
				Sequence: seq,
			},
			args: args{
				i: 2,
				j: 9,
			},
			want: seq.Time(2).Before(seq.Time(9)),
		},
		{
			fields: fields{
				Sequence: seq,
			},
			args: args{
				i: 8,
				j: 3,
			},
			want: seq.Time(8).Before(seq.Time(3)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sortableSequence{
				Sequence: tt.fields.Sequence,
			}
			if got := s.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("sortableSequence.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}
