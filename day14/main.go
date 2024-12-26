package main

import (
	"bufio"
	"fmt"
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

	robots := readInput(fileName)

	// part 1
	gridSize := XY{}
	if strings.Contains(fileName, "small") {
		gridSize = XY{11, 7}
	} else {
		gridSize = XY{101, 103}
	}

	seconds := 100
	for i := range robots {
		robots[i].pos.x += robots[i].vel.x * seconds
		robots[i].pos.y += robots[i].vel.y * seconds

		robots[i].pos.x = robots[i].pos.x % gridSize.x
		robots[i].pos.y = robots[i].pos.y % gridSize.y

		if robots[i].pos.x < 0 {
			robots[i].pos.x += gridSize.x
		}

		if robots[i].pos.y < 0 {
			robots[i].pos.y += gridSize.y
		}
	}

	q1, q2, q3, q4 := 0, 0, 0, 0
	center := XY{x: (gridSize.x - 1) / 2, y: (gridSize.y - 1) / 2}

	for _, robot := range robots {
		if robot.pos.x == center.x && robot.pos.y == center.y {
			continue
		}

		if robot.pos.x < center.x && robot.pos.y < center.y {
			q1++
		} else if robot.pos.x < center.x && robot.pos.y > center.y {
			q2++
		} else if robot.pos.x > center.x && robot.pos.y < center.y {
			q3++
		} else if robot.pos.x > center.x && robot.pos.y > center.y {
			q4++
		}
	}

	fmt.Printf("The total safety score is %d\n", q1*q2*q3*q4)
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
