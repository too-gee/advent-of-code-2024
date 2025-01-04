package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/too-gee/advent-of-code-2024/shared"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	gridSize := shared.Coord{X: 101, Y: 103}
	robots := readInput(fileName)

	// part 1
	safetyScore := PartOne(robots, gridSize)
	fmt.Printf("The total safety score is %d\n", safetyScore)

	// part 2
	timeToTree := PartTwo(robots, gridSize)
	fmt.Printf("Time to tree: %d\n", timeToTree)
}

func PartOne(input []Robot, gridSize shared.Coord) int {
	robots := copyRobots(input)
	center := shared.Coord{X: (gridSize.X - 1) / 2, Y: (gridSize.Y - 1) / 2}

	quads := map[bool]map[bool]int{false: {false: 0, true: 0}, true: {false: 0, true: 0}}

	for i := range robots {
		robots[i].move(100, gridSize)

		if robots[i].pos.X == center.X || robots[i].pos.Y == center.Y {
			continue
		}

		quads[robots[i].pos.X < center.X][robots[i].pos.Y < center.Y]++
	}

	safetyScore := quads[false][false] * quads[false][true] * quads[true][false] * quads[true][true]

	return safetyScore
}

func PartTwo(input []Robot, gridSize shared.Coord) int {
	robots := copyRobots(input)

	minLength := 9999
	treeTime := 0
	for i := 1; i < 10000; i++ {
		for j := range robots {
			robots[j].move(1, gridSize)
		}

		len, grid := plot(robots, gridSize)

		if len < minLength {
			fmt.Printf("Time: %d, Length: %d\n", i, len)

			for y := 0; y < gridSize.Y; y++ {
				for x := 0; x < gridSize.X; x++ {
					if grid[y][x] {
						fmt.Print("⬛")
					} else {
						fmt.Print("⬜")
					}
				}
				fmt.Println()
			}

			minLength = len
			treeTime = i
		}
	}

	return treeTime
}

func readInput(filePath string) []Robot {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	robots := []Robot{}

	for scanner.Scan() {
		line := scanner.Text()

		pX, pY, vX, vY := 0, 0, 0, 0
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &pX, &pY, &vX, &vY)
		robots = append(robots, Robot{pos: shared.Coord{X: pX, Y: pY}, vel: shared.Coord{X: vX, Y: vY}})
	}

	return robots
}

func copyRobots(robots []Robot) []Robot {
	newRobots := make([]Robot, len(robots))
	copy(newRobots, robots)
	return newRobots
}

type Robot struct {
	pos shared.Coord
	vel shared.Coord
}

func (r *Robot) move(seconds int, gridSize shared.Coord) {
	(*r).pos.X += (*r).vel.X * seconds
	(*r).pos.Y += (*r).vel.Y * seconds

	(*r).pos.X = (*r).pos.X % gridSize.X
	(*r).pos.Y = (*r).pos.Y % gridSize.Y

	if (*r).pos.X < 0 {
		(*r).pos.X += gridSize.X
	}

	if (*r).pos.Y < 0 {
		(*r).pos.Y += gridSize.Y
	}
}

func plot(r []Robot, gridSize shared.Coord) (int, [][]bool) {
	output := ""
	grid := make([][]bool, gridSize.Y)

	for y := 0; y < gridSize.Y; y++ {
		grid[y] = make([]bool, gridSize.X)
	}

	for _, robot := range r {
		grid[robot.pos.Y][robot.pos.X] = true
	}

	for y := 0; y < gridSize.Y; y++ {
		for x := 0; x < gridSize.X; x++ {
			if grid[y][x] {
				output += "#"
			} else {
				output += " "
			}
		}

		output += "\n"
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte(output))
	gzipWriter.Close()

	return buf.Len(), grid
}
