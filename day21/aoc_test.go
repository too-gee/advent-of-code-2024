package main

import (
	"testing"
)

type testCase struct {
	fileName      string
	function      func([]string, int) int
	dirKeypads    int
	expComplexity int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Solve, 2, 126384},
		{"input.txt", Solve, 2, 222670},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		PopulateKeypads()
		complexity := c.function(input, c.dirKeypads)
		if complexity != c.expComplexity {
			t.Errorf("%s: expected complexity: %d, got complexity: %d", c.fileName, c.expComplexity, complexity)
		}
	}
}
