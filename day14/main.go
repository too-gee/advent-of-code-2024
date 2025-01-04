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

	gridSize := shared.Coord{101, 103}
	center := shared.Coord{X: (gridSize.X - 1) / 2, Y: (gridSize.Y - 1) / 2}

	// part 1
	robots := readInput(fileName)

	quads := map[bool]map[bool]int{false: {false: 0, true: 0}, true: {false: 0, true: 0}}

	for i := range robots {
		robots[i].move(100, gridSize)

		if robots[i].pos.X == center.X || robots[i].pos.Y == center.Y {
			continue
		}

		quads[robots[i].pos.X < center.X][robots[i].pos.Y < center.Y]++
	}

	safetyScore := quads[false][false] * quads[false][true] * quads[true][false] * quads[true][true]
	fmt.Printf("The total safety score is %d\n", safetyScore)

	// part 2
	robots = readInput(fileName)

	minLength := 9999
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
		}
	}
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
		robots = append(robots, Robot{pos: shared.Coord{pX, pY}, vel: shared.Coord{vX, vY}})
	}

	return robots
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
