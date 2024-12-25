package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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
	totalCost := 0

	for i, clawMachine := range clawMachines {
		for aPresses := 0; aPresses <= 100; aPresses++ {
			xRem := float64(clawMachine.prize.x-(clawMachine.buttonA.x*aPresses)) / float64(clawMachine.buttonB.x)
			yRem := float64(clawMachine.prize.y-(clawMachine.buttonA.y*aPresses)) / float64(clawMachine.buttonB.y)

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

	fmt.Printf("Total cost: %d\n", totalCost)

	// part 2
	totalCost = 0

	for i, clawMachine := range clawMachines {
		clawMachine.prize.x += 10000000000000
		clawMachine.prize.y += 10000000000000

		// C, D → Prize
		C := float64(clawMachine.prize.x)
		D := float64(clawMachine.prize.y)

		// y = Kx + J → Button A (except we assume we're starting at 0,0)
		// y = Px + Q → Button B
		K := float64(clawMachine.buttonA.y) / float64(clawMachine.buttonA.x)

		P := float64(clawMachine.buttonB.y) / float64(clawMachine.buttonB.x)
		Q := D - (P * C)

		// intersection
		x := Q / (K - P)

		aPresses := int(math.Round(x / float64(clawMachine.buttonA.x)))
		bPresses := int(math.Round((C - x) / float64(clawMachine.buttonB.x)))

		cost, win := clawMachine.play(aPresses, bPresses)

		if win {
			totalCost += cost
			fmt.Printf("%d: %d x Button A; %d x Button B; %d tokens\n", i, aPresses, bPresses, cost)
		}
	}

	fmt.Printf("Total cost: %d\n", totalCost)
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
			clawMachines = append(clawMachines, ClawMachine{buttonA: XY{buttonAx, buttonAy}, buttonB: XY{buttonBx, buttonBy}, prize: XY{prizeX, prizeY}})
		}
	}

	return clawMachines
}

type XY struct {
	x int
	y int
}

type ClawMachine struct {
	buttonA XY
	buttonB XY
	prize   XY
}

func (c ClawMachine) play(aPresses int, bPresses int) (int, bool) {
	position := XY{0, 0}
	cost := 0

	position.x += aPresses * c.buttonA.x
	position.y += aPresses * c.buttonA.y
	cost += aPresses * 3

	position.x += bPresses * c.buttonB.x
	position.y += bPresses * c.buttonB.y
	cost += bPresses * 1

	return cost, position == c.prize
}

func (c ClawMachine) maxIters() int {
	maxDimension := math.Max(float64(c.prize.x), float64(c.prize.y))
	minIncrementX := math.Min(float64(c.buttonA.x), float64(c.buttonB.x))
	minIncrementY := math.Min(float64(c.buttonB.y), float64(c.buttonB.y))
	minIncrement := math.Min(minIncrementX, minIncrementY)

	return int(math.Trunc(maxDimension/minIncrement) + 1)
}

func epsilonCompare(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
