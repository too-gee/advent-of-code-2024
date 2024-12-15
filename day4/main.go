package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	// wordSearch[y][x] is the character at row y and column x
	wordSearch := readFileToCharSlice(fileName)

	// part 1
	const matchString = "XMAS"

	rows := len(wordSearch)
	cols := len(wordSearch[0])
	matchLength := len(matchString)

	matches := getMatchesInRuns(matchString, wordSearch, runsSouthEast(rows, cols, matchLength))
	matches += getMatchesInRuns(matchString, wordSearch, runsSouthWest(rows, cols, matchLength))
	matches += getMatchesInRuns(matchString, wordSearch, runsSouth(rows, cols))
	matches += getMatchesInRuns(matchString, wordSearch, runsEast(rows, cols))

	fmt.Printf("There are %d XMAS matches.\n", matches)

	// part 2
	matchPatterns := [][][]string{
		{{"M", ".", "M"}, {".", "A", "."}, {"S", ".", "S"}},
		{{"M", ".", "S"}, {".", "A", "."}, {"M", ".", "S"}},
		{{"S", ".", "S"}, {".", "A", "."}, {"M", ".", "M"}},
		{{"S", ".", "M"}, {".", "A", "."}, {"S", ".", "M"}},
	}

	xMatches := 0
	for _, matchPattern := range matchPatterns {
		xMatches += getXMatches(matchPattern, wordSearch)
	}

	fmt.Printf("There are %d X-MAS matches.\n", xMatches)
}

func getXMatches(matchPattern [][]string, charSlice [][]string) int {
	matchWidth := len(matchPattern[0])
	matchHeight := len(matchPattern)

	searchWidth := len(charSlice[0])
	searchHeight := len(charSlice)

	matches := 0

	for x := 0; x < searchWidth-matchWidth+1; x++ {
		for y := 0; y < searchHeight-matchHeight+1; y++ {
			subSlice := make([][]string, matchHeight)

			for row := 0; row < matchHeight; row++ {
				subSlice[row] = charSlice[y+row][x : x+matchWidth]
			}

			if isMatch(matchPattern, subSlice) {
				matches++
			}
		}
	}

	return matches
}

func isMatch(matchPattern [][]string, charSlice [][]string) bool {
	for x := 0; x < len(matchPattern[0]); x++ {
		for y := 0; y < len(matchPattern); y++ {
			if matchPattern[y][x] == "." {
				continue
			}

			if matchPattern[y][x] != charSlice[y][x] {
				return false
			}
		}
	}

	return true
}

func getMatchesInRuns(matchString string, charSlice [][]string, runs [][]coordinates) int {
	matchCount := 0

	matchStringReversed := reverseString(matchString)

	for _, run := range runs {
		stringBuffer := strings.Repeat(" ", len(matchString))

		for _, coord := range run {
			stringBuffer = stringBuffer[1:] + charSlice[coord.y][coord.x]

			if stringBuffer == matchString || stringBuffer == matchStringReversed {
				matchCount++
			}
		}
	}

	return matchCount
}

func runsSouthWest(rows int, cols int, minLength int) [][]coordinates {
	var result [][]coordinates

	for i := minLength - 1; i < rows+cols-minLength; i++ {
		currentRun := []coordinates{}

		skipAhead := int(math.Max(0, float64(i-rows+1)))

		x := i - skipAhead
		for y := skipAhead; y < rows; y++ {
			currentRun = append(currentRun, coordinates{x: x, y: y})
			x--

			if x < 0 {
				break
			}
		}

		result = append(result, currentRun)
	}

	return result
}

func runsSouthEast(rows int, cols int, minLength int) [][]coordinates {
	var result [][]coordinates

	for i := minLength - rows; i < cols-minLength+1; i++ {
		currentRun := []coordinates{}

		skipAhead := int(math.Max(0, float64(-i)))

		x := i + skipAhead
		for y := skipAhead; y < rows; y++ {
			currentRun = append(currentRun, coordinates{x: x, y: y})
			x++

			if x >= cols {
				break
			}
		}

		result = append(result, currentRun)
	}

	return result
}

func runsEast(rows int, cols int) [][]coordinates {
	var result [][]coordinates

	for y := 0; y < cols; y++ {
		currentRun := []coordinates{}

		for x := 0; x < rows; x++ {
			currentRun = append(currentRun, coordinates{x: x, y: y})
		}

		result = append(result, currentRun)
	}

	return result
}

func runsSouth(rows int, cols int) [][]coordinates {
	var result [][]coordinates

	for x := 0; x < cols; x++ {
		currentRun := []coordinates{}

		for y := 0; y < rows; y++ {
			currentRun = append(currentRun, coordinates{x: x, y: y})
		}

		result = append(result, currentRun)
	}

	return result
}

func readFileToCharSlice(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var charSlice [][]string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := []string{}
		for _, ch := range line {
			row = append(row, string(ch))
		}
		charSlice = append(charSlice, row)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return charSlice
}

func reverseString(input string) string {
	chars := ""

	for i := len(input) - 1; i >= 0; i-- {
		chars = chars + string(input[i])
	}

	return chars
}

type coordinates struct {
	x int
	y int
}
