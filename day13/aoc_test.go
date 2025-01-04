package main

import "testing"

type testCase struct {
	fileName string
	function func([]ClawMachine) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 480},
		{"input_small.txt", PartTwo, 875318608908},
		{"input.txt", PartOne, 29438},
		{"input.txt", PartTwo, 104958599303720},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
