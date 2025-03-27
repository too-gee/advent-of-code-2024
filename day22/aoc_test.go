package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func([]int) int
	iters    int
	expSum   int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 2000, 37327623},
		{"input.txt", Part1, 2000, 20506453102},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		sum := c.function(input)
		if sum != c.expSum {
			t.Errorf("%s: expected sum: %d, got sum: %d", c.fileName, c.expSum, sum)
		}
	}
}
