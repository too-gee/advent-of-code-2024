package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"math"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	registers, program := readInput(fileName)
	fmt.Println(program)
	registers, output := Solve(registers, program)
	fmt.Println(registers)
	fmt.Println(output)

}

func readInput(filePath string) ([]int, []int) {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var a, b, c int
	var program []int

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Sscanf(line, "Register A: %d", &a)
		fmt.Sscanf(line, "Register B: %d", &b)
		fmt.Sscanf(line, "Register C: %d", &c)

		var tmpProgram string
		_, err := fmt.Sscanf(line, "Program: %s", &tmpProgram)

		if err == nil {
			tmpProgram := strings.Split(tmpProgram, ",")
			for _, tmpStr := range tmpProgram {
				tmpInt, _ := strconv.Atoi(tmpStr)
				program = append(program, tmpInt)
			}
		}
	}

	return []int{a, b, c}, program
}

type State struct {
	Program []int
	Pointer int

	RegisterA int
	RegisterB int
	RegisterC int

	Output []string
}

func (s State) RenderOutput() string {
	return strings.Join(s.Output, ",")
}

// OPCODE 0
func (s *State) adv(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	numerator := float64((*s).RegisterA)
	denominator := math.Pow(2, float64(operand))

	(*s).RegisterA = int(math.Trunc(numerator / denominator))
	(*s).Pointer += 2
}

// OPCODE 1
func (s *State) bxl(operand int) {
	(*s).RegisterB ^= operand
	(*s).Pointer += 2
}

// OPCODE 2
func (s *State) bst(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	(*s).RegisterB = operand % 8
	(*s).Pointer += 2
}

// OPCODE 3
func (s *State) jnz(operand int) {
	if (*s).RegisterA == 0 {
		(*s).Pointer += 2
		return
	}

	(*s).Pointer = operand
}

// OPCODE 4
func (s *State) bxc(operand int) {
	(*s).RegisterB ^= (*s).RegisterC
	(*s).Pointer += 2
}

// OPCODE 5
func (s *State) out(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	result := operand % 8
	(*s).Output = append((*s).Output, strconv.Itoa(result))
	(*s).Pointer += 2
}

// OPCODE 6
func (s *State) bdv(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	numerator := float64((*s).RegisterA)
	denominator := math.Pow(2, float64(operand))

	(*s).RegisterB = int(math.Trunc(numerator / denominator))
}

// OPCODE 7
func (s *State) cdv(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	numerator := float64((*s).RegisterA)
	denominator := math.Pow(2, float64(operand))

	(*s).RegisterC = int(math.Trunc(numerator / denominator))
}

func Solve(registers []int, program []int) ([]int, string) {
	state := State{
		Program: program,
		Pointer: 0,

		RegisterA: registers[0],
		RegisterB: registers[1],
		RegisterC: registers[2],

		Output: []string{},
	}

	for state.Pointer < len(state.Program) {
		opcode := state.Program[state.Pointer]
		operand := state.Program[state.Pointer+1]

		switch opcode {
			case 0:
				state.adv(operand)
			case 1:
				state.bxl(operand)
			case 2:
				state.bst(operand)
			case 3:
				state.jnz(operand)
			case 4:
				state.bxc(operand)
			case 5:
				state.out(operand)
			case 6:
				state.bdv(operand)
			case 7:
				state.cdv(operand)
		}
	}

	return []int{state.RegisterA, state.RegisterB, state.RegisterC}, state.RenderOutput()
}

func (s *State) ResolveComboOperand(operand int) int {
	switch operand {
	case 4:
		return s.RegisterA
	case 5:
		return s.RegisterB
	case 6:
		return s.RegisterC
	case 7:
		panic("Invalid operand value!")
	default:
		return operand
	}
}