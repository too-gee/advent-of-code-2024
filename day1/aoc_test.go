package main

import "testing"

type testCase struct {
	fileName string
	function func([]int, []int) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 11},
		{"input_small.txt", PartTwo, 31},
		{"input.txt", PartOne, 2344935},
		{"input.txt", PartTwo, 27647262},
	}

	for _, c := range cases {
		input1, input2 := readInput(c.fileName)
		result := c.function(input1, input2)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
