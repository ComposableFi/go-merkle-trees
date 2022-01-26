package mmr_test

import (
	"reflect"
	"testing"

	merklego_mmr "github.com/ComposableFi/go-merkle-trees/mmr"
)

func TestNext(t *testing.T) {
	tests := map[string]struct {
		input [][]byte
		want  []byte
	}{
		"length of one": {input: [][]byte{uint32ToHash(1)}, want: uint32ToHash(1)},
		"length of two": {input: [][]byte{uint32ToHash(1), uint32ToHash(2)}, want: uint32ToHash(1)},
	}

	for name, test := range tests {
		iter := merklego_mmr.Iterator{Items: test.input}
		got := iter.Next()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: expected %v  got %v", name, test.want, got)
		}
	}

	tests = map[string]struct {
		input [][]byte
		want  []byte
	}{
		"length of two": {input: [][]byte{uint32ToHash(1), uint32ToHash(2)}, want: uint32ToHash(2)},
	}

	for name, test := range tests {
		iter := merklego_mmr.Iterator{Items: test.input}
		_ = iter.Next()
		got := iter.Next()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: expected %v  got %v", name, test.want, got)
		}
	}
}
