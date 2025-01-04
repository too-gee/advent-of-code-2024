package main

import "testing"

type testCase struct {
	fileName string
	function func(Grid, []string) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 2028},
		{"input_small.txt", PartTwo, 1751},
		{"input_medium.txt", PartOne, 10092},
		{"input_medium.txt", PartTwo, 9021},
		{"input.txt", PartOne, 1526018},
		{"input.txt", PartTwo, 1550677},
	}

	for _, c := range cases {
		input1, input2 := readInput(c.fileName)
		result := c.function(input1, input2)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
