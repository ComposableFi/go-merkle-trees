package mmr

import (
	"fmt"
	"testing"
)

func TestCountZeros(t *testing.T) {
	tests := map[string]struct{
		input uint64
		want int
	}{
		"2 zeros": {input: 5, want: 1},
		"3 zeros": {input: 20, want: 3},
	}

	for name, test := range tests {
		got := countZeros(test.input)
		if got != test.want {
			t.Errorf("%s: want %v  got %v", name, test.want, got)
		}
	}
}

func TestPosHeightInTree(t *testing.T) {
	tests := map[string]struct{
		input uint64
		want uint32
	}{
		"pos 0": {input: 0, want: 0},
		"pos 1": {input: 1, want: 0},
		"pos 2": {input: 2, want: 1},
		"pos 3": {input: 3, want: 0},
		"pos 4": {input: 4, want: 0},
		"pos 6": {input: 6, want: 2},
		"pos 7": {input: 7, want: 0},
	}

	for name, test := range tests {
		fmt.Printf("posHeight %v ", test.input)
		got := posHeightInTree(test.input)
		fmt.Printf("...done\n")
		if got != test.want {
			t.Errorf("%s: want %v  got %v", name, test.want, got)
		}
	}
}