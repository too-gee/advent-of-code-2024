package main

import "testing"

type testCase struct {
	fileName string
	function func(Map) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 36},
		{"input_small.txt", PartTwo, 81},
		{"input.txt", PartOne, 776},
		{"input.txt", PartTwo, 1657},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
