package main

import "testing"

type testCase struct {
	fileName string
	function func(Stones) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", PartOne, 55312},
		{"input_small.txt", PartTwo, 65601038650482},
		{"input.txt", PartOne, 218956},
		{"input.txt", PartTwo, 259593838049805},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result := c.function(input)
		if result != c.expected {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, result)
		}
	}
}
