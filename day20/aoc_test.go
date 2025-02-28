package main

import (
	"testing"

	"github.com/too-gee/advent-of-code-2024/shared"
)

type testCase struct {
	fileName       string
	function       func(shared.Grid, int, int) int
	minSavings     int
	maxCheatLength int
	expectedCount  int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Solve, 64, 2, 1},
		{"input_small.txt", Solve, 40, 2, 2},
		{"input_small.txt", Solve, 38, 2, 3},
		{"input_small.txt", Solve, 36, 2, 4},
		{"input_small.txt", Solve, 20, 2, 5},
		{"input_small.txt", Solve, 12, 2, 8},
		{"input_small.txt", Solve, 10, 2, 10},
		{"input_small.txt", Solve, 8, 2, 14},
		{"input_small.txt", Solve, 6, 2, 16},
		{"input_small.txt", Solve, 4, 2, 30},
		{"input_small.txt", Solve, 2, 2, 44},
		{"input.txt", Solve, 100, 2, 1321},
		{"input_small.txt", Solve, 50, 20, 285},
		{"input_small.txt", Solve, 52, 20, 253},
		{"input_small.txt", Solve, 54, 20, 222},
		{"input_small.txt", Solve, 56, 20, 193},
		{"input_small.txt", Solve, 58, 20, 154},
		{"input_small.txt", Solve, 60, 20, 129},
		{"input_small.txt", Solve, 62, 20, 106},
		{"input_small.txt", Solve, 64, 20, 86},
		{"input_small.txt", Solve, 66, 20, 67},
		{"input_small.txt", Solve, 68, 20, 55},
		{"input_small.txt", Solve, 70, 20, 41},
		{"input_small.txt", Solve, 72, 20, 29},
		{"input_small.txt", Solve, 74, 20, 7},
		{"input_small.txt", Solve, 76, 20, 3},
		{"input.txt", Solve, 100, 20, 971737},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		cheats := c.function(input, c.minSavings, c.maxCheatLength)
		if cheats != c.expectedCount {
			t.Errorf("%s: expected count: %d, got count: %d", c.fileName, c.expectedCount, cheats)
		}
	}
}
