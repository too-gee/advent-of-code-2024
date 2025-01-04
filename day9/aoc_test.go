package main

import "testing"

type testCase struct {
	fileName string
	function func(Disk) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 1928},
		{"input_small.txt", PartTwo, 2858},
		{"input.txt", PartOne, 6154342787400},
		{"input.txt", PartTwo, 6183632723350},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
