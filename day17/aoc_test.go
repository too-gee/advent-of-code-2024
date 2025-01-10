package main

import "testing"

type testCase struct {
	fileName          string
	function          func([]int, []int) ([]int, string)
	expectedRegisters []int
	expectedOutput    string
}

func TestAll(t *testing.T) {
	cases := []testCase{
		{"input_small_bst.txt", Solve, []int{0, 1, 9}, ""},
		{"input_small_out.txt", Solve, []int{10, 0, 0}, "0,1,2"},
		{"input_small_adv.txt", Solve, []int{0, 0, 0}, "4,2,5,6,7,7,7,7,3,1,0"},
		{"input_small_bxl.txt", Solve, []int{0, 26, 0}, ""},
		{"input_small_bxc.txt", Solve, []int{0, 44354, 43690}, ""},
		{"input_small.txt", Solve, []int{0, 0, 0}, "4,6,3,5,6,3,5,2,1,0"},
		{"input.txt", Solve, []int{0, 7, 0}, "4,1,5,3,1,5,3,5,7"},
	}

	for _, c := range cases {
		registers, program := readInput(c.fileName)
		registers, output := c.function(registers, program)
		if registers[0] != c.expectedRegisters[0] || registers[1] != c.expectedRegisters[1] || registers[2] != c.expectedRegisters[2] || output != c.expectedOutput {
			t.Errorf("%s: expected registers: %v - output: \"%s\", got registers: %v - output: \"%s\"", c.fileName, c.expectedRegisters, c.expectedOutput, registers, output)
		}
	}
}
