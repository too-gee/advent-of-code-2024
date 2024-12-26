package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	gridSize := XY{101, 103}
	center := XY{x: (gridSize.x - 1) / 2, y: (gridSize.y - 1) / 2}

	// part 1
	robots := readInput(fileName)

	quads := map[bool]map[bool]int{false: {false: 0, true: 0}, true: {false: 0, true: 0}}

	for i := range robots {
		robots[i].move(100, gridSize)

		if robots[i].pos.x == center.x || robots[i].pos.y == center.y {
			continue
		}

		quads[robots[i].pos.x < center.x][robots[i].pos.y < center.y]++
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

			for y := 0; y < gridSize.y; y++ {
				for x := 0; x < gridSize.x; x++ {
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
		robots = append(robots, Robot{pos: XY{pX, pY}, vel: XY{vX, vY}})
	}

	return robots
}

type XY struct {
	x int
	y int
}

type Robot struct {
	pos XY
	vel XY
}

func (r *Robot) move(seconds int, gridSize XY) {
	(*r).pos.x += (*r).vel.x * seconds
	(*r).pos.y += (*r).vel.y * seconds

	(*r).pos.x = (*r).pos.x % gridSize.x
	(*r).pos.y = (*r).pos.y % gridSize.y

	if (*r).pos.x < 0 {
		(*r).pos.x += gridSize.x
	}

	if (*r).pos.y < 0 {
		(*r).pos.y += gridSize.y
	}
}

func plot(r []Robot, gridSize XY) (int, [][]bool) {
	output := ""
	grid := make([][]bool, gridSize.y)

	for y := 0; y < gridSize.y; y++ {
		grid[y] = make([]bool, gridSize.x)
	}

	for _, robot := range r {
		grid[robot.pos.y][robot.pos.x] = true
	}

	for y := 0; y < gridSize.y; y++ {
		for x := 0; x < gridSize.x; x++ {
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
