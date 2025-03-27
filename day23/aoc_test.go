package main

import (
	"testing"
)

type testCase struct {
	fileName    string
	function    func([]Link) int
	expClusters int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Part1, 7},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		clusters := c.function(input)

		if clusters != c.expClusters {
			t.Errorf("%s: expected: %d, got: %d", c.fileName, c.expClusters, clusters)
		}
	}
}
