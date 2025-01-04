package main

import "testing"

type testCase struct {
	fileName string
	function func(area) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 41},
		{"input_small.txt", PartTwo, 6},
		{"input.txt", PartOne, 5534},
		{"input.txt", PartTwo, 2262},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
