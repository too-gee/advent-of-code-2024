package main

import "testing"

type testCase struct {
	fileName string
	function func([][]int, [][]int) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 143},
		{"input_small.txt", PartTwo, 123},
		{"input.txt", PartOne, 5991},
		{"input.txt", PartTwo, 5479},
	}

	for _, c := range cases {
		input1, input2 := readInput(c.fileName)
		result := c.function(input1, input2)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
