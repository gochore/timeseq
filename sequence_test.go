package timeseq

import "testing"

func Test_sortableSequence_Less(t *testing.T) {
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
		// TODO: Add test cases.
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
