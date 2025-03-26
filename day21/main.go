package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/too-gee/advent-of-code-2024/shared"
)

var MEMO = map[string][]string{}
var LMEMO = map[string]int{}

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

	// part 1 & part 2
	fmt.Printf("The complexity for part 1 is %d.\n", Solve(codes, 2))
	fmt.Printf("The complexity for part 2 is %d.\n", Solve(codes, 25))
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
	numTrips := getSequences(NUM, code)

	lengthOfSequence := recursiveCount(DIR, numTrips, dirKeypads)
	numericPart, _ := strconv.Atoi(code[:3])

	return numericPart * lengthOfSequence
}

// Memoized version of getSequencesRaw
func getSequences(keypad shared.Grid, sequence string) []string {
	memoId := fmt.Sprintf("trip-%s", sequence)

	if cached, ok := MEMO[memoId]; ok {
		return cached
	}

	roundTrips := getSequencesRaw(keypad, sequence)

	// memoize and return
	MEMO[memoId] = roundTrips
	return roundTrips
}

// ACCEPTS: A sequence implied to begin at A and ending explicitly at A
//
// RETURNS: A slice of sequences implied to begin at A and ending explicitly
//
//	at A that, when entered together, produce the original input
//	sequence
func getSequencesRaw(keypad shared.Grid, sequence string) []string {
	possible := [][]string{}

	paddedSequence := PRESS + sequence
	for i := 0; i < len(paddedSequence)-1; i++ {
		tmp := getPresses(keypad, string(paddedSequence[i]), string(paddedSequence[i+1]))

		newRoutes := make([]string, len(tmp))
		copy(newRoutes, tmp)
		for j := range newRoutes {
			newRoutes[j] = newRoutes[j] + PRESS
		}

		possible = append(possible, newRoutes)
	}

	found := getPermutations("", possible)
	valid := onlyValid(keypad, sequence, found)

	return valid
}

func recursiveCount(keypad shared.Grid, sequences []string, levels int) int {
	bestCount := math.MaxInt64

	if levels == 0 {
		for _, trip := range sequences {
			if len(trip) < bestCount {
				bestCount = len(trip)
			}
		}

		return bestCount
	}

	for _, sequence := range sequences {
		subs := splitSequence(sequence)
		seqLength := 0

		for _, sub := range subs {
			memoId := fmt.Sprintf("exp-%s-%d", sub, levels)

			if cached, ok := LMEMO[memoId]; ok {
				seqLength += cached
				continue
			}

			newSeqs := getSequences(keypad, sub)
			newLength := recursiveCount(keypad, newSeqs, levels-1)

			LMEMO[memoId] = newLength
			seqLength += newLength

			if seqLength > bestCount {
				break
			}
		}

		if seqLength < bestCount {
			bestCount = seqLength
		}
	}

	return bestCount
}

func getPresses(keypad shared.Grid, start string, end string) []string {
	memoId := fmt.Sprintf("route-%s%s", start, end)

	if cached, ok := MEMO[memoId]; ok {
		return cached
	}

	routes := getPressesRaw(keypad, start, end)

	// memoize and return
	MEMO[memoId] = routes
	return routes
}

func getPressesRaw(keypad shared.Grid, start string, end string) []string {
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

func getPermutations(base string, blocks [][]string) []string {
	if len(blocks) == 0 {
		return []string{base}
	}

	perms := []string{}

	for _, option := range blocks[0] {
		var newBase string
		if base == "" {
			newBase = option
		} else {
			newBase = base + option
		}

		perms = slices.Concat(perms, getPermutations(newBase, blocks[1:]))
	}

	return onlyShortest(perms)
}

func splitSequence(sequence string) []string {
	trips := []string{}

	buffer := ""

	for _, char := range sequence {
		buffer += string(char)
		if string(char) == PRESS {
			trips = append(trips, buffer)
			buffer = ""
		}
	}

	return trips
}

func onlyShortest(strings []string) []string {
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

func onlyValid(keypad shared.Grid, sequence string, trips []string) []string {
	output := []string{}

	for _, trip := range trips {
		result, _ := evalPresses(keypad, trip)

		if result == sequence {
			output = append(output, trip)
		}
	}

	return output
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

const UP = "^"
const DOWN = "v"
const LEFT = "<"
const RIGHT = ">"
const GAP = ""
const PRESS = "A"
