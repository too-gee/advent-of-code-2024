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

	// Part 1
	registers, output := Part1(registers, program)
	fmt.Printf("Part 1: registers - %v, output - %s\n", registers, output)

	// Part 2
	registers, output = Part2(registers, program)
	fmt.Printf("Part 2: registers - %v, Initial Register A - %s\n", registers, output)
}

func readInput(filePath string) ([]int64, []int) {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var a, b, c int64
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

	return []int64{a, b, c}, program
}

type State struct {
	Program []int
	Pointer int

	RegisterA int64
	RegisterB int64
	RegisterC int64

	Output []int
}

func (s *State) Execute() {
	for (*s).Pointer < len((*s).Program) {
		opcode := (*s).Program[(*s).Pointer]
		operand := (*s).Program[(*s).Pointer+1]

		switch opcode {
			case 0:
				(*s).adv(operand)
			case 1:
				(*s).bxl(operand)
			case 2:
				(*s).bst(operand)
			case 3:
				(*s).jnz(operand)
			case 4:
				(*s).bxc(operand)
			case 5:
				(*s).out(operand)
			case 6:
				(*s).bdv(operand)
			case 7:
				(*s).cdv(operand)
		}
	}
}

func (s State) RenderOutput() string {
	return strings.Join(IntsToStrings(s.Output), ",")
}

func (s State) DebugOutput(a int64, msg string, compare []int) []int {
	s.RegisterA = a
	s.Execute()

	aOctal := strconv.FormatInt(int64(a), 8)

	display := "["
	for i, digit := range s.Output {
		if len(compare) == len(s.Output) && digit != compare[i] {
			display += fmt.Sprintf("«%s»", strconv.Itoa(digit))
		} else {
			display += fmt.Sprintf(" %s ", strconv.Itoa(digit))
		}
	}
	display += "]"

	fmt.Printf("%-15s @ %15d / %19s: %v\n", msg, a, fmt.Sprintf("0o%s" ,aOctal), display)

	return s.Output
}

// OPCODE 0
func (s *State) adv(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	numerator := float64((*s).RegisterA)
	denominator := math.Pow(2, float64(operand))

	(*s).RegisterA = int64(math.Trunc(numerator / denominator))
	(*s).Pointer += 2
}

// OPCODE 1
func (s *State) bxl(operand int) {
	(*s).RegisterB ^= int64(operand)
	(*s).Pointer += 2
}

// OPCODE 2
func (s *State) bst(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	(*s).RegisterB = int64(operand % 8)
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
	result := int(operand % 8)
	(*s).Output = append((*s).Output, result)
	(*s).Pointer += 2
}

// OPCODE 6
func (s *State) bdv(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	numerator := float64((*s).RegisterA)
	denominator := math.Pow(2, float64(operand))

	(*s).RegisterB = int64(math.Trunc(numerator / denominator))
	(*s).Pointer += 2
}

// OPCODE 7
func (s *State) cdv(comboOperand int) {
	operand := (*s).ResolveComboOperand(comboOperand)
	numerator := float64((*s).RegisterA)
	denominator := math.Pow(2, float64(operand))

	(*s).RegisterC = int64(math.Trunc(numerator / denominator))
	(*s).Pointer += 2
}

func Part2(registers []int64, program []int) ([]int64, string) {
	floor := RegisterAForLength(len(program), registers, program)
	ceiling := RegisterAForLength(len(program)+1, registers, program)-1

	state := State{
		Program: program,
		Pointer: 0,
		RegisterA: registers[0],
		RegisterB: registers[1],
		RegisterC: registers[2],
		Output: []int{},
	}

	state.DebugOutput(floor - 1, "low-miss", nil)
	state.DebugOutput(floor, "floor", nil)
	state.DebugOutput(ceiling, "ceiling", nil)
	state.DebugOutput(ceiling + 1, "high-miss", nil)

	fmt.Println("--------------")

	registerA := FindMatch(state, floor)

	return []int64{0,0,0}, strconv.FormatInt(registerA, 10)
}

/*
	Once you realize that the program output acts like an odometer (hopefully, mine did)
	and increases sequentially, you can jump ahead to find when each digit rolls over.
	This was easier for me than reverse engineering the program. In order to run in a
	reasonable amount of time, this had to be tuned a bunch. I ended up with the numbers
	here after a bunch of trial and error to reduce the number of iterations. If the
	parameters are wrong, the process will overshoot on fast iterations and waste time
	catching up on slow iterations. This approach runs several hundred million times
	faster than a fully naive approach. On my machine, that's the difference between
	~10 seconds and ~40 years.
*/
func FindMatch(state State, start int64) int64 {
	places := octalPlaces(start)

	bookmark := start
	registerA := start

	var output []int
	var inc int64

	for i := places-1; i >= 0; i-- {
		inc = int64(math.Pow(8, float64(i)))

		for {
			registerA += inc

			output = state.DebugOutput(registerA, fmt.Sprintf("digit: %d, speed: %d", i, inc), state.Program)

			match := Compare(output[i:], state.Program[i:])
			tooBig := octalPlaces(registerA) > places

			if tooBig || match {
				if tooBig {
					registerA = bookmark
				}

				if match {
					bookmark = registerA
					registerA -= inc
				}

				if inc == 1 {
					inc = 0
					continue
				} else {
					inc /= 8
				}
			}

			if inc == 0 { break }
		}

		registerA = bookmark

		if Compare(output, state.Program) { break }
	}

	return registerA
}

func RegisterAForLength(length int, registers []int64, program []int) int64 {
	var registerA int64
	registerA = 0

	inc := int64(math.Pow(8, 20))

	state := State{
		Program: program,
		Pointer: 0,
		RegisterA: registers[0],
		RegisterB: registers[1],
		RegisterC: registers[2],
		Output: []int{},
	}

	for ; inc >= 1; registerA += inc {
		output := state.DebugOutput(registerA, "finding length", nil)

		if len(output) >= length {
			registerA = int64(math.Max(float64(registerA-inc),0))
			inc /= 8
		}
	}
	fmt.Println(registerA + 1)
	return registerA + 1
}

func Part1(registers []int64, program []int) ([]int64, string) {
	state := State{
		Program: program,
		Pointer: 0,

		RegisterA: registers[0],
		RegisterB: registers[1],
		RegisterC: registers[2],

		Output: []int{},
	}

	state.Execute()

	return []int64{state.RegisterA, state.RegisterB, state.RegisterC}, state.RenderOutput()
}

func (s *State) ResolveComboOperand(operand int) int64 {
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
		return int64(operand)
	}
}

func IntsToStrings(ints []int) []string {
	output := []string{}
	for _, i := range ints {
		output = append(output, strconv.Itoa(i))
	}
	return output
}

func StringsToInts(strings []string) []int {
	ints := []int{}
	for _, i := range strings {
		result, _ := strconv.Atoi(i)
		ints = append(ints, result)
	}
	return ints
}

func Compare(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func octalPlaces(num int64) int {
    if num == 0 {
        return 1
    }

    return int(math.Floor(math.Log(float64(num)) / math.Log(8))) + 1
}