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

	blocks := readInput(fileName)

	var size, initialBlocks int

	if strings.HasSuffix(fileName, "input.txt") {
		size = 70
		initialBlocks = 1024
	} else {
		size = 6
		initialBlocks = 12
	}

	// Part 1
	bestCost := Part1(blocks, initialBlocks, size)
	fmt.Printf("The shortest path is %d steps\n", bestCost)

	// Part 2
	blockedAt := Part2(blocks, initialBlocks, size)
	fmt.Printf("Blocked at %d (%v)\n", blockedAt, blocks[blockedAt])
}

func readInput(filePath string) []shared.Coord {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var blocks []shared.Coord
	var x, y int

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Sscanf(line, "%d,%d", &x, &y)
		blocks = append(blocks, shared.Coord{X: x, Y: y})
	}

	return blocks
}

func Part1(blocks []shared.Coord, initialBlocks int, size int) int {
	grid := createGrid(blocks[:initialBlocks], size+1)
	start := shared.Coord{X: 0, Y: 0}
	end := shared.Coord{X: size, Y: size}
	// grid.Draw(nil, map[string][]shared.Coord{"🟢": {start}, "🔴": {end}})

	forwardCosts := GetCosts(grid, end)
	reverseCosts := GetCosts(grid, start)

	bestCost := math.MaxInt
	for y := range grid.Height() {
		for x := range grid.Width() {
			loc := shared.Coord{X: x, Y: y}

			forwardCost, fok := forwardCosts[loc]
			reverseCost, rok := reverseCosts[loc]

			if fok && rok {
				if bestCost > forwardCost+reverseCost {
					bestCost = forwardCost + reverseCost
				}
			}
		}
	}

	// subtract 2 since the first step in the first path is the start
	// and the first step in the second path is included in the first path
	return bestCost - 2
}

func Part2(blocks []shared.Coord, initialBlocks int, size int) int {
	grid := createGrid(blocks[:initialBlocks], size+1)
	start := shared.Coord{X: 0, Y: 0}
	end := shared.Coord{X: size, Y: size}

	// Apply the initial blocks to the grid
	for i := 0; i < len(blocks)-1; i++ {
		// drop the block
		grid[blocks[i].Y][blocks[i].X] = "#"

		// If the block is one of the initial blocks, skip ahead since we know
		// that we can reach the end.
		if i < initialBlocks {
			continue
		}

		// Initially, I was checking to see if the new block was touching
		// another block or the wall before checking to see if we could still
		// reach the end, but this actually slowed things down. On my machine,
		// I went from ~4.5 seconds to solve to ~3.2 seconds. That might be
		// different for a different puzzle input.

		// Check to see if we can still reach the end.
		if !Flood(grid, start, end) {
			return i
		}
	}

	return -1
}

func Flood(g shared.Grid, start shared.Coord, end shared.Coord) bool {
	queue := DumbQueue{}
	visited := []shared.Coord{start}

	queue.push(start)

	for len(queue) > 0 {
		current := queue.pop()

		if current == end {
			return true
		}

		neighbors := g.Neighbors(current, []string{"#"})

		for _, neighbor := range neighbors {
			if slices.Contains(visited, neighbor) {
				continue
			}

			queue.push(neighbor)
			visited = append(visited, neighbor)
		}
	}

	return false
}

func createGrid(blocks []shared.Coord, size int) shared.Grid {
	grid := shared.Grid(make([][]string, size))

	for i := range size {
		grid[i] = make([]string, size)
	}

	for y := range size {
		for x := range size {
			coord := shared.Coord{X: x, Y: y}

			if slices.Contains(blocks, coord) {
				grid[y][x] = "#"
			} else {
				grid[y][x] = "."
			}
		}
	}

	return grid
}

type DumbQueue []shared.Coord

func (q *DumbQueue) push(item shared.Coord) { *q = append(*q, item) }

func (q *DumbQueue) pop() shared.Coord {
	item := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return item
}

type State struct {
	path  []shared.Coord
	lives int
}

type Queue []*State

func (q Queue) Len() int { return len(q) }

func (q Queue) Less(i, j int) bool {
	return len(q[i].path) < len(q[j].path)
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *Queue) Push(x interface{}) {
	item := x.(*State)
	*q = append(*q, item)
}

func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*q = old[0 : n-1]
	return item
}

// NewPriorityQueue creates and initializes a new priority queue
func NewPriorityQueue() *Queue {
	q := &Queue{}
	heap.Init(q)
	return q
}

// PushState adds a new State to the priority queue
func (q *Queue) PushState(state *State) {
	heap.Push(q, state)
}

// PopState removes and returns the State with the shortest path
func (q *Queue) PopState() *State {
	return heap.Pop(q).(*State)
}

func CopyAppend(a []shared.Coord, b shared.Coord) []shared.Coord {
	newSlice := make([]shared.Coord, len(a)+1)
	copy(newSlice, a)
	newSlice[len(a)] = b

	return newSlice
}

func GetCosts(m shared.Grid, start shared.Coord) map[shared.Coord]int {
	lives := 10
	queue := Queue{}
	queue.PushState(&State{path: []shared.Coord{start}, lives: lives})

	visitedCost := map[shared.Coord]int{}

	// for debug
	counter := 0

	for len(queue) > 0 {
		current := queue.PopState()

		if current.lives == 0 {
			continue
		}

		length := len(current.path)
		loc := current.path[length-1]

		// for debug
		counter++
		if counter%10000000 == 0 {
			m.Draw(nil, map[string][]shared.Coord{"🟢": {start}, "🟨": current.path[1:]})
		}

		// fill a cost grid
		localShortest, ok := visitedCost[loc]

		if !ok || length < localShortest {
			visitedCost[loc] = length
			localShortest = length
		}

		// unless this is an improvement, skip
		if length > localShortest {
			continue
		}

		neighbors := m.Neighbors(loc, []string{"#"})

		obstacles := 4 - len(neighbors)
		if loc.X == 0 || loc.X == m.Width()-1 {
			obstacles--
		}
		if loc.Y == 0 || loc.Y == m.Height()-1 {
			obstacles--
		}

		if obstacles == 0 {
			current.lives--
		} else {
			current.lives = int(math.Min(float64(lives), float64(current.lives+1)))
		}

		for _, neighbor := range neighbors {
			if slices.Contains(current.path, neighbor) {
				continue
			}

			queue.PushState(&State{path: CopyAppend(current.path, neighbor), lives: current.lives})
		}
	}

	return visitedCost
}
