package main

import (
	"testing"
)

type testCase struct {
	fileName string
	function func([]Link) string
	exp      string
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, "7"},
		{"input.txt", Part1, "1437"},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		str := c.function(input)

		if str != c.exp {
			t.Errorf("%s: expected: %s, got: %s", c.fileName, c.exp, str)
		}
	}
}
