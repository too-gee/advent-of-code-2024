package main

import "testing"

type testCase struct {
	fileName      string
	function      func(Maze) (int, int)
	expectedCost  int
	expectedTiles int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_toy.txt", Solve, 6075, 76},
		{"input_small.txt", Solve, 7036, 45},
		{"input_medium.txt", Solve, 11048, 64},
		{"input.txt", Solve, 101492, 543},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		cost, tiles := c.function(input)
		if cost != c.expectedCost || tiles != c.expectedTiles {
			t.Errorf("%s: expected cost: %d - tiles: %d, got cost: %d - tiles: %d", c.fileName, c.expectedCost, c.expectedTiles, cost, tiles)
		}
	}
}
