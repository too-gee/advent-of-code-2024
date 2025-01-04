package main

import (
	"testing"

	"github.com/too-gee/advent-of-code-2024/shared"
)

type testCase struct {
	fileName string
	function func([]Robot, shared.Coord) int
	gridSize shared.Coord
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_medium.txt", PartOne, shared.Coord{X: 11, Y: 7}, 12},
		{"input.txt", PartOne, shared.Coord{X: 101, Y: 103}, 218965032},
		{"input.txt", PartTwo, shared.Coord{X: 101, Y: 103}, 7037},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input, c.gridSize)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
