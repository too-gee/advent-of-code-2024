package main

import "testing"

type testCase struct {
	fileName string
	function func([]Region) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small_ab.txt", PartOne, 1184},
		{"input_small_ab.txt", PartTwo, 368},
		{"input_small_abcde.txt", PartOne, 140},
		{"input_small_abcde.txt", PartTwo, 80},
		{"input_small_ex.txt", PartOne, 692},
		{"input_small_ex.txt", PartTwo, 236},
		{"input_small_xo.txt", PartOne, 772},
		{"input_small_xo.txt", PartTwo, 436},
		{"input_small.txt", PartOne, 1930},
		{"input_small.txt", PartTwo, 1206},
		{"input.txt", PartOne, 1467094},
		{"input.txt", PartTwo, 881182},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
