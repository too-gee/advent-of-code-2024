package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func([]int) int
	expSum   int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 37327623},
		{"input.txt", Part1, 20506453102},
		{"input_small2.txt", Part2, 23},
		{"input_small.txt", Part2, 24},
		{"input.txt", Part2, 2423},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		sum := c.function(input)

		if sum != c.expSum {
			t.Errorf("%s: expected: %d, got: %d", c.fileName, c.expSum, sum)
		}
	}
}
