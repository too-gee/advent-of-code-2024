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
	pt1complexity := Solve(codes, 12)
	fmt.Println(pt1complexity)

	for k, v := range MEMO {
		fmt.Printf("%s -> %v\n", k, v)
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
		solutions = findMultileg(DIR, solutions)
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

func findMultileg(keypad shared.Grid, sequences []string) []string {
	legs := []string{}

	for _, possible := range sequences {
		tripSolutions := [][]string{}
		for _, trip := range splitTrips(possible) {
			newSolutions := findAllRoundTrips(keypad, trip)
			tripSolutions = append(tripSolutions, newSolutions)
		}

		legs = append(legs, findShortestPermutation("", tripSolutions))
	}

	return shortest(legs)[0:]
}

func findAllRoundTrips(keypad shared.Grid, sequence string) []string {
	if string(sequence[0]) != PRESS || sequence[0] != sequence[len(sequence)-1] {
		panic("This only works for round trips")
	}

	memoId := fmt.Sprintf("trip-%s", sequence)

	if cached, ok := MEMO[memoId]; ok {
		return cached
	}

	roundTrips := findAllRoundTripsRaw(keypad, sequence)

	// memoize and return
	MEMO[memoId] = roundTrips
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
	memoId := fmt.Sprintf("route-%s%s", start, end)

	if cached, ok := MEMO[memoId]; ok {
		return cached
	}

	routes := findAllRoutesRaw(keypad, start, end)

	// memoize and return
	MEMO[memoId] = routes
	return routes
}

func findAllRoutesRaw(keypad shared.Grid, start string, end string) []string {
	startLoc := keypad.LocationOf(start)
	endLoc := keypad.LocationOf(end)

	var xMoves string
	var yMoves string

	if endLoc.X > startLoc.X {
		xMoves = strings.Repeat(RIGHT, endLoc.X-startLoc.X)
	} else {
		xMoves = strings.Repeat(LEFT, startLoc.X-endLoc.X)
	}

	if endLoc.Y > startLoc.Y {
		yMoves = strings.Repeat(DOWN, endLoc.Y-startLoc.Y)
	} else {
		yMoves = strings.Repeat(UP, startLoc.Y-endLoc.Y)
	}

	if xMoves != "" && yMoves != "" {
		return []string{xMoves + yMoves, yMoves + xMoves}
	} else {
		return []string{xMoves + yMoves}
	}

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

	return shortest(perms)
}

func findShortestPermutation(base string, options [][]string) string {
	if len(options) == 0 {
		return base
	}

	output := base

	for i := range options {
		output += shortest(options[i])[0]
	}

	return output
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

func shortest(strings []string) []string {
	shortest := 999999999
	count := 0

	for _, x := range strings {
		if len(x) < shortest {
			shortest = len(x)
			count = 0
		}

		if len(x) == shortest {
			count++
		}
	}

	output := make([]string, count)
	count = 0

	for _, x := range strings {
		if len(x) == shortest {
			output[count] = x
			count++
		}
	}

	return output
}

const UP = "^"
const DOWN = "v"
const LEFT = "<"
const RIGHT = ">"
const GAP = ""
const PRESS = "A"
