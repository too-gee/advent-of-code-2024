package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func([][]int, [][]int) int
	exp      int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 3},
		{"input.txt", Part1, 2978},
	}

	for _, c := range cases {
		locks, keys := readInput(c.fileName)
		output := c.function(locks, keys)

		if output != c.exp {
			t.Errorf("%s: expected: %d, got: %d", c.fileName, c.exp, output)
		}
	}
}
