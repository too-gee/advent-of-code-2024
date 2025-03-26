package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func([]int, int) int
	iters    int
	expSum   int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Solve, 2000, 37327623},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		sum := c.function(input, c.iters)
		if sum != c.expSum {
			t.Errorf("%s: expected sum: %d, got sum: %d", c.fileName, c.expSum, sum)
		}
	}
}
