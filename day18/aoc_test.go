package main

import (
	"testing"

	"github.com/too-gee/advent-of-code-2024/shared"
)

type testCase struct {
	fileName      string
	function      func([]shared.Coord, int, int) int
	initialBlocks int
	size          int
	expected      int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 12, 6, 22},
		{"input.txt", Part1, 1024, 70, 260},
		{"input_small.txt", Part2, 12, 6, 20},
		{"input.txt", Part2, 1024, 70, 2881},
	}

	for _, c := range cases {
		output := c.function(readInput(c.fileName), c.initialBlocks, c.size)
		if output != c.expected {
			t.Errorf("%s: expected: %d, got registers: %d", c.fileName, c.expected, output)
		}
	}
}
