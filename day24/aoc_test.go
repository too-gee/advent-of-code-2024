package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func(Connections) int
	exp      int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 4},
		{"input_medium.txt", Part1, 2024},
	}

	for _, c := range cases {
		conns := readInput(c.fileName)
		output := c.function(conns)

		if output != c.exp {
			t.Errorf("%s: expected: %d, got: %d", c.fileName, c.exp, output)
		}
	}
}
