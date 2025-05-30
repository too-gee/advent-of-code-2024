package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func(Connections) string
	exp      string
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, "4"},
		{"input_medium.txt", Part1, "2024"},
		{"input.txt", Part1, "56729630917616"},
		{"input.txt", Part2, "bjm,hsw,nvr,skf,wkr,z07,z13,z18"},
	}

	for _, c := range cases {
		conns := readInput(c.fileName)
		output := c.function(conns)

		if output != c.exp {
			t.Errorf("%s: expected: %s, got: %s", c.fileName, c.exp, output)
		}
	}
}
