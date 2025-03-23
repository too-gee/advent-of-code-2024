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

var MEMO = map[string]string{}

var NUM = shared.Grid{}
var DIR = shared.Grid{}

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	codes := readInput(fileName)

	NUM = shared.Grid{{"7", "8", "9"}, {"4", "5", "6"}, {"1", "2", "3"}, {GAP, "0", PRESS}}
	DIR = [][]string{{GAP, UP, PRESS}, {LEFT, DOWN, RIGHT}}

	// part 1
	complexity := 0
	for _, code := range codes {
		complexity += Part1(code)
	}
	fmt.Println(complexity)

	start := "029A"
	enc1 := findShortestSequence(NUM, start)
	enc2 := findShortestSequence(DIR, enc1)
	final := findShortestSequence(DIR, enc2)
	fmt.Printf("%s -> %s -> %s -> %s\n", start, enc1, enc2, final)
	dec1, _ := evalPresses(DIR, final)
	dec2, _ := evalPresses(DIR, dec1)
	proof, _ := evalPresses(NUM, dec2)
	fmt.Printf("%s -> %s -> %s -> %s\n", final, dec1, dec2, proof)

	fmt.Println()
	result := "<v<A>>^AAAvA^A<vA<AA>>^AvAA<^A>A<v<A>A>^AAAvA<^A>A<vA>^A<A>A"
	fmt.Println(result)
	result, _ = evalPresses(DIR, result)
	fmt.Println(result)
	result, _ = evalPresses(DIR, result)
	fmt.Println(result)
	result, _ = evalPresses(NUM, result)
	fmt.Println(result)
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

	route1 := findShortestSequence(numeric, code)
	route2 := findShortestSequence(directional, route1)
	final := findShortestSequence(directional, route2)

	fmt.Printf("%s :: route1: %v\n", code, route1)
	//fmt.Printf("%s :: route2: %v\n", code, route2)
	//fmt.Printf("%s :: final: %v\n", code, final)
	fmt.Printf("%s :: %v\n", code, final)
	fmt.Println(evalPresses(numeric, route1))

	numericPart, _ := strconv.Atoi(code[:3])
	lengthOfSequence := len(final)

	//fmt.Printf("%s :: %d * %d = %d\n", code, numericPart, lengthOfSequence, numericPart*lengthOfSequence)

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

	var horizontalDir int
	var horizontalMove string

	if current.X > end.X {
		horizontalDir = -1
		horizontalMove = LEFT
	} else {
		horizontalDir = 1
		horizontalMove = RIGHT
	}

	var verticalDir int
	var verticalMove string

	if current.Y > end.Y {
		verticalDir = -1
		verticalMove = UP
	} else {
		verticalDir = 1
		verticalMove = DOWN
	}

	for current != end {
		for current.X != end.X && keypad.At(shared.Coord{X: current.X + horizontalDir, Y: current.Y}) != GAP {
			moves = append(moves, horizontalMove)
			current.X += horizontalDir
		}

		for current.Y != end.Y && keypad.At(shared.Coord{X: current.X, Y: current.Y + verticalDir}) != GAP {
			moves = append(moves, verticalMove)
			current.Y += verticalDir
		}
	}

	return append(moves, PRESS)
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

func findShortestSequence(keypad shared.Grid, sequence string) string {
	solution := ""

	moves := strings.Split("A"+sequence+"A", "")
	for i := 0; i < len(moves)-1; i++ {
		moveId := moves[i] + moves[i+1]
		if saved, ok := MEMO[moveId]; ok {
			fmt.Println("CACHE HIT")
			solution += saved + "A"
		} else {
			hardWork := findShortestRoute(keypad, moveId)
			MEMO[moveId] = hardWork
			solution += hardWork + "A"
		}
	}

	return solution[:len(solution)-1]
}

func findShortestRoute(keypad shared.Grid, moveId string) string {
	queue := DumbQueue{""}

	a := string(moveId[0])
	b := string(moveId[1])

	init := strings.Repeat(UP, (keypad.Width()*keypad.Height())+2)
	best := init
	found := []string{}

	endLoc := keypad.LocationOf(b)

	for len(queue) > 0 {
		current := queue.pop()

		currentLoc := keypad.LocationOf(a)
		validPath := true
		for _, dir := range strings.Split(current, "") {
			switch dir {
			case UP:
				currentLoc.Y--
			case DOWN:
				currentLoc.Y++
			case LEFT:
				currentLoc.X--
			case RIGHT:
				currentLoc.X++
			}

			if !keypad.Contains(currentLoc) || keypad.At(currentLoc) == GAP {
				validPath = false
				break
			}
		}

		if !validPath {
			continue
		}

		fmt.Printf("queue: %d, current: %s\n", len(queue), current)

		if best != init && len(current) > len(best) {
			continue
		}

		if currentLoc == endLoc {
			best = current
			found = append(found, current)
			continue
		}

		if !strings.Contains(current, UP) {
			queue.push(current + DOWN)
		}

		if !strings.Contains(current, DOWN) {
			queue.push(current + UP)
		}

		if !strings.Contains(current, LEFT) {
			queue.push(current + RIGHT)
		}

		if !strings.Contains(current, RIGHT) {
			queue.push(current + LEFT)
		}
	}

	if len(found) > 1 {
		fmt.Printf("found %s: %v\n", moveId, found)
	}

	return best
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
