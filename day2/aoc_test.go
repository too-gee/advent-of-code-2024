package main

import "testing"

type testCase struct {
	fileName string
	function func([][]int) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 2},
		{"input_small.txt", PartTwo, 4},
		{"input.txt", PartOne, 472},
		{"input.txt", PartTwo, 520},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
