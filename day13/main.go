package main

import (
	"bufio"
	"fmt"
	"math"
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

	clawMachines := readInput(fileName)

	// part 1
	totalCost := PartOne(clawMachines)
	fmt.Printf("Total cost: %d\n", totalCost)

	// part 2
	totalCost = PartTwo(clawMachines)
	fmt.Printf("Total cost: %d\n", totalCost)
}

func PartOne(clawMachines []ClawMachine) int {
	totalCost := 0

	for i, clawMachine := range clawMachines {
		for aPresses := 0; aPresses <= 100; aPresses++ {
			xRem := float64(clawMachine.prize.X-(clawMachine.buttonA.X*aPresses)) / float64(clawMachine.buttonB.X)
			yRem := float64(clawMachine.prize.Y-(clawMachine.buttonA.Y*aPresses)) / float64(clawMachine.buttonB.Y)

			if xRem != yRem || xRem != math.Trunc(xRem) {
				continue
			}

			bPresses := int(xRem)

			cost, win := clawMachine.play(aPresses, bPresses)

			if win {
				totalCost += cost
				fmt.Printf("%d: %d x Button A; %d x Button B; %d tokens\n", i, aPresses, bPresses, cost)
				break
			}
		}
	}

	return totalCost
}

func PartTwo(clawMachines []ClawMachine) int {
	totalCost := 0

	for i, clawMachine := range clawMachines {
		clawMachine.prize.X += 10000000000000
		clawMachine.prize.Y += 10000000000000

		// C, D → Prize
		C := float64(clawMachine.prize.X)
		D := float64(clawMachine.prize.Y)

		// y = Kx + J → Button A (except we assume we're starting at 0,0)
		// y = Px + Q → Button B
		K := float64(clawMachine.buttonA.Y) / float64(clawMachine.buttonA.X)

		P := float64(clawMachine.buttonB.Y) / float64(clawMachine.buttonB.X)
		Q := D - (P * C)

		// intersection
		x := Q / (K - P)

		aPresses := int(math.Round(x / float64(clawMachine.buttonA.X)))
		bPresses := int(math.Round((C - x) / float64(clawMachine.buttonB.X)))

		cost, win := clawMachine.play(aPresses, bPresses)

		if win {
			totalCost += cost
			fmt.Printf("%d: %d x Button A; %d x Button B; %d tokens\n", i, aPresses, bPresses, cost)
		}
	}

	return totalCost
}

func readInput(filePath string) []ClawMachine {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	clawMachines := []ClawMachine{}
	var buttonAx, buttonAy, buttonBx, buttonBy, prizeX, prizeY int

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Sscanf(line, "Button A: X+%d, Y+%d", &buttonAx, &buttonAy)
		fmt.Sscanf(line, "Button B: X+%d, Y+%d", &buttonBx, &buttonBy)
		i, err := fmt.Sscanf(line, "Prize: X=%d, Y=%d", &prizeX, &prizeY)

		if err == nil && i == 2 {
			clawMachines = append(clawMachines, ClawMachine{buttonA: shared.Coord{X: buttonAx, Y: buttonAy}, buttonB: shared.Coord{X: buttonBx, Y: buttonBy}, prize: shared.Coord{X: prizeX, Y: prizeY}})
		}
	}

	return clawMachines
}

type ClawMachine struct {
	buttonA shared.Coord
	buttonB shared.Coord
	prize   shared.Coord
}

func (c ClawMachine) play(aPresses int, bPresses int) (int, bool) {
	position := shared.Coord{X: 0, Y: 0}
	cost := 0

	position.X += aPresses * c.buttonA.X
	position.Y += aPresses * c.buttonA.Y
	cost += aPresses * 3

	position.X += bPresses * c.buttonB.X
	position.Y += bPresses * c.buttonB.Y
	cost += bPresses * 1

	return cost, position == c.prize
}

func (c ClawMachine) maxIters() int {
	maxDimension := math.Max(float64(c.prize.X), float64(c.prize.Y))
	minIncrementX := math.Min(float64(c.buttonA.X), float64(c.buttonB.X))
	minIncrementY := math.Min(float64(c.buttonB.Y), float64(c.buttonB.Y))
	minIncrement := math.Min(minIncrementX, minIncrementY)

	return int(math.Trunc(maxDimension/minIncrement) + 1)
}

func epsilonCompare(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
