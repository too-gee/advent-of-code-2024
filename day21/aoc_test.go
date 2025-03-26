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
		{"input_small.txt", Solve, 25, 154115708116294},
		{"input.txt", Solve, 25, 271397390297138},
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
