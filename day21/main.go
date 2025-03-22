package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/too-gee/advent-of-code-2024/shared"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	codes := readInput(fileName)

	// part 1
	complexity := 0
	for _, code := range codes {
		complexity += Part1(code)
	}
	fmt.Println(complexity)
}

func readInput(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	codes := []string{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		codes = append(codes, line)
	}

	return codes
}

func Part1(code string) int {
	numeric := shared.Grid{{"7", "8", "9"}, {"4", "5", "6"}, {"1", "2", "3"}, {GAP, "0", PRESS}}
	directional := [][]string{{GAP, UP, PRESS}, {LEFT, DOWN, RIGHT}}

	route1 := routeSequence(numeric, strings.Split(code, ""))
	route2 := routeSequence(directional, route1)
	final := routeSequence(directional, route2)

	fmt.Printf("%s :: route1: %v\n", code, route1)
	fmt.Printf("%s :: route2: %v\n", code, route2)
	fmt.Printf("%s :: final: %v\n", code, final)

	numericPart, _ := strconv.Atoi(code[:3])
	lengthOfSequence := len(final)

	fmt.Printf("%s :: %d * %d = %d\n", code, numericPart, lengthOfSequence, numericPart*lengthOfSequence)

	return numericPart * lengthOfSequence
}

func routeSequence(keypad shared.Grid, sequence []string) []string {
	route := []string{}

	previous := PRESS
	for _, current := range sequence {
		newSteps := routeButton(keypad, previous, current)
		route = slices.Concat(route, newSteps)
		previous = current
	}

	return route
}

func routeButton(keypad shared.Grid, start string, target string) []string {
	current := keypad.LocationOf(start)
	end := keypad.LocationOf(target)

	var moves []string

	for current != end {
		for current.X != end.X && keypad.At(shared.Coord{X: end.X, Y: current.Y}) != GAP {
			if current.X > end.X {
				moves = append(moves, LEFT)
				current.X--
			} else {
				moves = append(moves, RIGHT)
				current.X++
			}
		}

		for current.Y != end.Y && keypad.At(shared.Coord{X: current.X, Y: end.Y}) != GAP {
			if current.Y > end.Y {
				moves = append(moves, UP)
				current.Y--
			} else {
				moves = append(moves, DOWN)
				current.Y++
			}
		}
	}

	return append(moves, PRESS)
}

const UP = "^"
const DOWN = "v"
const LEFT = "<"
const RIGHT = ">"
const GAP = ""
const PRESS = "A"
