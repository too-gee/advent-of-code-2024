package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"slices"
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

	maze := readInput(fileName).FillDeadEnds()
	maze.DrawCustomMarkers(map[string]string{"=": "🟫", "S": "🟢", "E": "🔴"})

	// part 1
	//cost := PartOne(maze)
	//fmt.Printf("Lowest cost: %d\n", cost)

	// part 2
	optimalCosts := PartTwo(maze)
	endLoc := maze.LocationOf(END)
	endCosts := CostMap{}
	minCost := math.MaxInt32
	for _, dir := range []string{"N", "S", "E", "W"} {
		id := fmt.Sprintf("%d.%d.%s", endLoc.X, endLoc.Y, dir)
		if cost, ok := optimalCosts[id]; ok {
			if cost < minCost {
				endCosts = CostMap{id: cost}
				minCost = cost
			} else if cost == minCost {
				endCosts[id] = cost
			}
		}
	}

	fmt.Printf("Optimal path count: %d\n", len(endCosts))
}

func readInput(filePath string) Maze {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := Maze{}

	for scanner.Scan() {
		line := scanner.Text()

		gridRow := strings.Split(line, "")
		grid.Grid = append(grid.Grid, gridRow)
	}

	return grid
}

type DumbQueue []State

func (q *DumbQueue) push(item State) { *q = append(*q, item) }

func (q *DumbQueue) pop() State {
	item := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return item
}

type CostMap map[string]int

func (c *CostMap) set(metaLoc string, cost int) bool {
	_, hasKey := (*c)[metaLoc]
	if !hasKey || (*c)[metaLoc] < cost {
		(*c)[metaLoc] = cost
		return true
	}

	return false
}

func PartTwo(g Maze) CostMap {
	start := g.LocationOf(START)
	end := g.LocationOf(END)

	queue := DumbQueue{}
	queue.push(State{
		loc:  start,
		dir:  "E",
		cost: 0,
		path: []shared.Coord{start},
	})

	costs := CostMap{}

	for len(queue) > 0 {
		current := queue.pop()

		id := fmt.Sprintf("%d.%d.%s", current.loc.X, current.loc.Y, current.dir)
		fmt.Printf("Id - %9v :: Cost - %7d :: Queue - %3d\n", id, current.cost, len(queue))

		if !costs.set(id, current.cost) || current.loc == end {
			continue
		}

		for newDir, neighbor := range g.Neighbors(current.loc) {
			if slices.Contains(current.path, neighbor) {
				continue
			}

			newCost := current.cost + 1

			if newDir != current.dir {
				newCost += 1000
			}

			newPath := append(current.path, neighbor)

			queue.push(State{
				loc:  neighbor,
				dir:  newDir,
				cost: newCost,
				path: newPath,
			})
		}
	}

	return costs
}

func PartOne(m Maze) int {
	start := m.LocationOf(START)
	end := m.LocationOf(END)

	queue := PriorityQueue{}
	heap.Init(&queue)
	heap.Push(&queue, &State{
		loc:  start,
		dir:  "E",
		cost: 0,
		path: []shared.Coord{start},
	})

	visited := map[string]bool{}

	for queue.Len() > 0 {
		current := heap.Pop(&queue).(*State)

		if current.loc == end {
			return current.cost
		}

		id := fmt.Sprintf("%d.%d.%s", current.loc.X, current.loc.Y, current.dir)

		if visited[id] {
			continue
		} else {
			visited[id] = true
		}

		for newDir, neighbor := range m.Neighbors(current.loc) {
			newCost := current.cost + 1

			if newDir != current.dir {
				newCost += 1000
			}

			newPath := append(current.path, neighbor)

			heap.Push(&queue, &State{
				loc:  neighbor,
				dir:  newDir,
				cost: newCost,
				path: newPath,
			})
		}
	}

	return -1
}
