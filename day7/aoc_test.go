package main

import "testing"

type testCase struct {
	fileName string
	function func([][]int) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 3749},
		{"input_small.txt", PartTwo, 11387},
		{"input.txt", PartOne, 14711933466277},
		{"input.txt", PartTwo, 286580387663654},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
