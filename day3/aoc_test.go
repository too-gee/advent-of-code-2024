package main

import (
	"slices"
	"testing"
)

type testCase struct {
	fileName string
	function func([][]string) (int, int)
	expected []int
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small.txt", Solve, []int{161, 48}},
		{"input.txt", Solve, []int{164730528, 70478672}},
	}

	for _, c := range cases {
		input := readInput(c.fileName)
		result1, result2 := c.function(input)
		if !slices.Equal([]int{result1, result2}, c.expected) {
			t.Errorf("%s: expected %d, got %d", c.fileName, c.expected, []int{result1, result2})
		}
	}
}
