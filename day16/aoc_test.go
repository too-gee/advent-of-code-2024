package main

import "testing"

type testCase struct {
	fileName string
	function func(Grid) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_toy.txt", PartOne, 6075},
		{"input_small.txt", PartOne, 7036},
		{"input_medium.txt", PartOne, 11048},
		{"input.txt", PartOne, 101492},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
