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

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	grid := readInput(fileName)

	// part 1
	antiNodeLocationCount := PartOne(grid)
	fmt.Printf("There are %d unique locations that contain an antinode\n", antiNodeLocationCount)

	// part 2
	antiNodeLocationCount = PartTwo(grid)
	fmt.Printf("There are %d unique locations that contain a resonant antinode\n", antiNodeLocationCount)
}

func PartOne(grid shared.Grid) int {
	antiNodeLocations := coordList{}

	for _, antennas := range getAntennaLocations(grid) {
		pairs := pairs(makeRange(0, len(antennas)-1))

		for _, pair := range pairs {
			antennaA := antennas[pair[0]]
			antennaB := antennas[pair[1]]

			antinode := shared.Coord{X: (2 * antennaA.X) - antennaB.X, Y: (2 * antennaA.Y) - antennaB.Y}

			if grid.Contains(antinode) && !slices.Contains(antiNodeLocations, antinode) {
				antiNodeLocations = append(antiNodeLocations, antinode)
			}
		}
	}

	return len(antiNodeLocations)
}

func PartTwo(grid shared.Grid) int {
	antiNodeLocations := coordList{}

	for _, antennas := range getAntennaLocations(grid) {
		pairs := pairs(makeRange(0, len(antennas)-1))

		for _, pair := range pairs {
			antennaA := antennas[pair[0]]
			antennaB := antennas[pair[1]]

			if !slices.Contains(antiNodeLocations, antennaA) {
				antiNodeLocations = append(antiNodeLocations, antennaA)
			}

			rise := antennaA.Y - antennaB.Y
			run := antennaA.X - antennaB.X
			currentX := antennaA.X
			currentY := antennaA.Y

			for {
				antinode := shared.Coord{X: currentX + run, Y: currentY + rise}

				if !grid.Contains(antinode) {
					break
				}

				if !slices.Contains(antiNodeLocations, antinode) {
					antiNodeLocations = append(antiNodeLocations, antinode)
				}

				currentX += run
				currentY += rise
			}

		}
	}

	return len(antiNodeLocations)
}

func getAntennaLocations(grid shared.Grid) map[string]coordList {
	frequencies := getUniqueCharacters(grid)
	antennaLocations := map[string]coordList{}

	for _, frequency := range frequencies {
		antennaLocations[frequency] = coordList{}
		for y, row := range grid {
			for x, cell := range row {
				if cell == frequency {
					antennaLocations[frequency] = append(antennaLocations[frequency], shared.Coord{X: x, Y: y})
				}
			}
		}
	}

	return antennaLocations
}

func readInput(filePath string) shared.Grid {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	grid := shared.Grid{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		grid = append(grid, row)
	}

	return grid
}

func getUniqueCharacters(grid shared.Grid) []string {
	uniqueChars := make(map[string]bool)

	ignoredChars := []string{"."}

	for _, row := range grid {
		for _, cell := range row {
			if !slices.Contains(ignoredChars, cell) {
				uniqueChars[cell] = true
			}
		}
	}

	var uniqueCharsSlice []string

	for char := range uniqueChars {
		uniqueCharsSlice = append(uniqueCharsSlice, char)
	}

	return uniqueCharsSlice
}

func pairs(options []int) [][]int {
	permutations := [][]int{}

	for i := 0; i < int(math.Pow(float64(len(options)), float64(2))); i++ {
		binary := strconv.FormatInt(int64(i), len(options))
		binaryStr := fmt.Sprintf("%02s", binary)
		binarySlice := strings.Split(binaryStr, "")

		intSlice := []int{}

		for _, str := range binarySlice {
			intVal, _ := strconv.Atoi(str)
			intSlice = append(intSlice, options[intVal])
		}

		if intSlice[0] != intSlice[1] {
			permutations = append(permutations, intSlice)
		}
	}

	return permutations
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

type coordList []shared.Coord
