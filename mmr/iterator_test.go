package mmr

import (
	"testing"
)

func TestNext(t *testing.T) {
	tests := map[string]struct {
		input []interface{}
		want  int
	}{
		"length of one": {input: []interface{}{1}, want: 1},
		"length of two": {input: []interface{}{1, 2}, want: 1},
	}

	for name, test := range tests {
		iter := Iterator{item: test.input}
		got := iter.next()
		if got != test.want {
			t.Errorf("%s: expected %v  got %v", name, test.want, got)
		}
	}

	tests = map[string]struct {
		input []interface{}
		want  int
	}{
		"call next twice": {input: []interface{}{1, 2}, want: 2},
	}

	for name, test := range tests {
		iter := Iterator{item: test.input}
		_ = iter.next()
		got := iter.next()
		if got != test.want {
			t.Errorf("%s: expected %v  got %v", name, test.want, got)
		}
	}
}
