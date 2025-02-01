package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func([]string, []string) int
	expected int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 6},
		{"input.txt", Part1, 258},
	}

	for _, c := range cases {
		towels, designs := readInput(c.fileName)
		output := c.function(towels, designs)
		if output != c.expected {
			t.Errorf("%s: expected: %d, got: %d", c.fileName, c.expected, output)
		}
	}
}
