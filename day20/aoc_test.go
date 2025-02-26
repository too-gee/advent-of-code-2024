package main

import (
	"testing"

	"github.com/too-gee/advent-of-code-2024/shared"
)

type testCase struct {
	fileName      string
	function      func(shared.Grid, int) int
	minSavings    int
	expectedCount int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 64, 1},
		{"input_small.txt", Part1, 40, 2},
		{"input_small.txt", Part1, 38, 3},
		{"input_small.txt", Part1, 36, 4},
		{"input_small.txt", Part1, 20, 5},
		{"input_small.txt", Part1, 12, 8},
		{"input_small.txt", Part1, 10, 10},
		{"input_small.txt", Part1, 8, 14},
		{"input_small.txt", Part1, 6, 16},
		{"input_small.txt", Part1, 4, 30},
		{"input_small.txt", Part1, 2, 44},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		cheats := c.function(input, c.minSavings)
		if cheats != c.expectedCount {
			t.Errorf("%s: expected count: %d, got count: %d", c.fileName, c.expectedCount, cheats)
		}
	}
}
