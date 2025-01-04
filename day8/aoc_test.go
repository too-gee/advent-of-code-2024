package main

import "testing"

type testCase struct {
	fileName string
	function func([][]string) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 14},
		{"input_small.txt", PartTwo, 34},
		{"input.txt", PartOne, 295},
		{"input.txt", PartTwo, 1034},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
