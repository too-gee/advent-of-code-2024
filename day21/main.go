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

var MEMO = map[string][]string{}

var NUM = shared.Grid{}
var DIR = shared.Grid{}

func PopulateKeypads() {
	NUM = shared.Grid{{"7", "8", "9"}, {"4", "5", "6"}, {"1", "2", "3"}, {GAP, "0", PRESS}}
	DIR = [][]string{{GAP, UP, PRESS}, {LEFT, DOWN, RIGHT}}
}

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	codes := readInput(fileName)

	PopulateKeypads()

	// part 1
	pt1complexity := Solve(codes, 2)
	fmt.Println(pt1complexity)

	for k, _ := range MEMO {
		fmt.Println(k)
	}
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

func Solve(codes []string, dirKeypads int) int {
	complexity := 0

	for _, code := range codes {
		complexity += getComplexity(code, dirKeypads)
	}

	return complexity
}

func getComplexity(code string, dirKeypads int) int {
	solutions := findAllRoundTrips(NUM, PRESS+code)

	for range dirKeypads {
		solutions = findAllMultileg(DIR, solutions)
	}

	final := solutions[0]

	numericPart, _ := strconv.Atoi(code[:3])
	lengthOfSequence := len(final)

	return numericPart * lengthOfSequence
}

func evalPresses(keypad shared.Grid, sequence string) (string, error) {
	result := ""

	location := keypad.LocationOf(PRESS)

	for _, code := range strings.Split(sequence, "") {
		switch code {
		case UP:
			location.Y--
		case DOWN:
			location.Y++
		case LEFT:
			location.X--
		case RIGHT:
			location.X++
		}

		if keypad.At(location) == GAP || !keypad.Contains(location) {
			return "", fmt.Errorf("invalid location %v", location)
		}

		if code == PRESS {
			result += keypad.At(location)
		}
	}

	return result, nil
}

func findAllMultileg(keypad shared.Grid, sequences []string) []string {
	legs := []string{}

	for _, possible := range sequences {
		tripSolutions := [][]string{}
		for _, trip := range splitTrips(possible) {
			newSolutions := findAllRoundTrips(keypad, trip)
			tripSolutions = append(tripSolutions, newSolutions)
		}

		legs = slices.Concat(legs, findAllPermutations("", tripSolutions))
	}

	short := []string{}
	bestLength := 9999

	for i := range legs {
		if len(legs[i]) < bestLength {
			bestLength = len(legs[i])
		}
	}

	for i := range legs {
		if len(legs[i]) <= bestLength {
			short = append(short, legs[i])
		}
	}

	return short
}

func findAllRoundTrips(keypad shared.Grid, sequence string) []string {
	if string(sequence[0]) != PRESS || sequence[0] != sequence[len(sequence)-1] {
		panic("This only works for round trips")
	}

	var roundTrips []string

	if cached, ok := MEMO[sequence]; ok {
		return cached
	} else {
		roundTrips = findAllRoundTripsRaw(keypad, sequence)
	}
	// memoize and return
	MEMO[sequence] = roundTrips
	return roundTrips
}

func findAllRoundTripsRaw(keypad shared.Grid, sequence string) []string {
	possible := [][]string{}

	for i := 0; i < len(sequence)-1; i++ {
		tmp := findAllRoutes(keypad, string(sequence[i]), string(sequence[i+1]))

		newRoutes := make([]string, len(tmp))
		copy(newRoutes, tmp)
		for j := range newRoutes {
			newRoutes[j] = newRoutes[j] + PRESS
		}

		possible = append(possible, newRoutes)
	}

	found := findAllPermutations("", possible)

	valid := []string{}

	for i := range found {
		result, err := evalPresses(keypad, found[i])
		if err == nil || result == sequence {
			valid = append(valid, found[i])
		}
	}

	return valid
}

func findAllRoutes(keypad shared.Grid, start string, end string) []string {
	startLoc := keypad.LocationOf(start)
	endLoc := keypad.LocationOf(end)

	var xDist int
	var yDist int
	var xMove string
	var yMove string

	if endLoc.X > startLoc.X {
		xDist = endLoc.X - startLoc.X
		xMove = RIGHT
	} else {
		xDist = startLoc.X - endLoc.X
		xMove = LEFT
	}

	if endLoc.Y > startLoc.Y {
		yDist = endLoc.Y - startLoc.Y
		yMove = DOWN
	} else {
		yDist = startLoc.Y - endLoc.Y
		yMove = UP
	}

	queue := DumbQueue{""}

	solutions := []string{}

	for len(queue) > 0 {
		current := queue.pop()

		xMoves := strings.Count(current, xMove)
		yMoves := strings.Count(current, yMove)

		if xMoves < xDist {
			queue.push(current + xMove)
		}

		if yMoves < yDist {
			queue.push(current + yMove)
		}

		if xMoves == xDist && yMoves == yDist {
			solutions = append(solutions, current)
		}
	}

	return solutions
}

func findAllPermutations(base string, options [][]string) []string {
	if len(options) == 0 {
		return []string{base}
	}

	perms := []string{}

	for _, option := range options[0] {
		var newBase string
		if base == "" {
			newBase = option
		} else {
			newBase = base + option
		}

		perms = slices.Concat(perms, findAllPermutations(newBase, options[1:]))
	}

	return perms
}

func splitTrips(sequence string) []string {
	trips := []string{}

	buffer := ""

	for _, char := range sequence {
		buffer += string(char)
		if string(char) == PRESS {
			trips = append(trips, PRESS+buffer)
			buffer = ""
		}
	}

	return trips
}

type DumbQueue []string

func (q *DumbQueue) push(item string) { *q = append(*q, item) }

func (q *DumbQueue) pop() string {
	item := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return item
}

const UP = "^"
const DOWN = "v"
const LEFT = "<"
const RIGHT = ">"
const GAP = ""
const PRESS = "A"
