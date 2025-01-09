package main

import (
	"bufio"
	"container/heap"
	"crypto/rand"
	"encoding/hex"
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
	maze.Draw(map[string]string{"=": "🟫", "S": "🟢", "E": "🔴"}, nil)

	bestCost, optimalTiles := Solve(maze)
	fmt.Printf("Lowest cost: %d\n", bestCost)
	fmt.Printf("Optimal path count: %d\n", optimalTiles)
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

func Solve(m Maze) (int, int) {
	start := m.LocationOf(START)
	end := m.LocationOf(END)

	queue := DumbQueue{}
	queue.push(State{
		id:   "start",
		loc:  start,
		dir:  "E",
		cost: 0,
		path: []shared.Coord{start},
	})

	bestCost := math.MaxInt
	visitedCost := map[shared.Coord]int{}
	bestPaths := map[string][]shared.Coord{}

	for len(queue) > 0 {
		current := queue.pop()

		localBestCost, ok := visitedCost[current.loc]

		if !ok || current.cost < localBestCost {
			visitedCost[current.loc] = current.cost
			localBestCost = current.cost
		}

		if current.cost > bestCost || current.cost > (localBestCost+1000) {
			continue
		}

		if current.loc == end {
			if !strings.HasSuffix(current.id, fmt.Sprintf("[%d,%d]", current.loc.X, current.loc.Y)) {
				current.id = fmt.Sprintf("%s -> [%d,%d]", current.id, current.loc.X, current.loc.Y)
			}

			if current.cost < bestCost {
				bestPaths = map[string][]shared.Coord{current.id: current.path}
				bestCost = current.cost
			} else if current.cost == bestCost {
				bestPaths[current.id] = current.path
			}
			continue
		}

		// This could be done in the loop below but this makes debugging easier
		neighbors := m.Neighbors(current.loc)
		newNeighbors := map[string]shared.Coord{}
		for newDir, neighbor := range neighbors {
			if !slices.Contains(current.path, neighbor) {
				newNeighbors[newDir] = neighbor
			}
		}

		for newDir, neighbor := range newNeighbors {
			newCost := current.cost + 1

			if newDir != current.dir {
				newCost += 1000
			}

			newPath := CopyAppend(current.path, neighbor)

			var newId string
			if len(newNeighbors) > 1 {
				newId = fmt.Sprintf("%s -> [%d,%d]", current.id, current.loc.X, current.loc.Y)
			} else if newDir != current.dir {
				newId = fmt.Sprintf("%s -> [%d,%d]", current.id, current.loc.X, current.loc.Y)
			} else {
				newId = current.id
			}

			queue.push(State{
				id:   newId,
				loc:  neighbor,
				dir:  newDir,
				cost: newCost,
				path: newPath,
			})
		}
	}

	bestTiles := []shared.Coord{}
	for _, path := range bestPaths {
		for _, loc := range path {
			if !slices.Contains(bestTiles, loc) {
				bestTiles = append(bestTiles, loc)
			}
		}
	}

	m.Draw(map[string]string{"=": "🟫", "S": "🟢", "E": "🔴"}, map[string][]shared.Coord{"⏺️ ": bestTiles})

	return bestCost, len(bestTiles)
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

func generateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set the version bits to 4
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set the variant bits to RFC 4122
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return hex.EncodeToString(uuid), nil
}

func CopyAppend(a []shared.Coord, b shared.Coord) []shared.Coord {
	newSlice := make([]shared.Coord, len(a)+1)
	copy(newSlice, a)
	newSlice[len(a)] = b

	return newSlice
}
