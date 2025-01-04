package main

import (
	"bufio"
	"container/heap"
	"fmt"
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

	maze := readInput(fileName).fillDeadEnds()
	maze.draw(map[string][]shared.Coord{})

	// part 1
	cost := PartOne(maze)
	fmt.Printf("Lowest cost: %d\n", cost)

	// part 2
	optimalPathCount := PartTwo(maze)
	fmt.Printf("Optimal path count: %d\n", optimalPathCount)
}

func PartTwo(g Grid) int {
	return 0
}

func readInput(filePath string) Grid {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := Grid{}

	for scanner.Scan() {
		line := scanner.Text()

		gridRow := strings.Split(line, "")
		grid = append(grid, gridRow)
	}

	return grid
}

const FILL = "="
const WALL = "#"
const EMPTY = "."
const START = "S"
const END = "E"

func (g Grid) fillDeadEnds() Grid {
	for y := range g {
		for x := range g[0] {
			loc := shared.Coord{X: x, Y: y}
			if g.isPassable(loc) && g.isDeadEnd(loc) {
				g.fillDeadEnd(loc)
			}
		}
	}

	return g
}

func (g Grid) inGrid(loc shared.Coord) bool {
	return loc.Y >= 0 &&
		loc.Y < len(g) &&
		loc.X >= 0 &&
		loc.X < len(g[loc.Y])
}

func (g Grid) cellType(loc shared.Coord) string {
	return g[loc.Y][loc.X]
}

func (g Grid) neighbors(loc shared.Coord) map[string]shared.Coord {
	neighbors := map[string]shared.Coord{}

	for dir, neighbor := range loc.Neighbors() {
		if g.inGrid(neighbor) &&
			!slices.Contains([]string{WALL, FILL}, g[neighbor.Y][neighbor.X]) {
			neighbors[dir] = neighbor
		}
	}

	return neighbors
}

func (g Grid) isPassable(loc shared.Coord) bool {
	return slices.Contains([]string{EMPTY, START, END}, g[loc.Y][loc.X])
}

func (g Grid) isDeadEnd(loc shared.Coord) bool {
	if slices.Contains([]string{START, END}, g[loc.Y][loc.X]) {
		return false
	}

	return len(g.neighbors(loc)) == 1
}

type State struct {
	loc  shared.Coord
	dir  string
	cost int
	path []shared.Coord
}

type PriorityQueue []*State

// Impliment a bunch of methods to satisfy the heap interface
func (q PriorityQueue) Len() int {
	return len(q)
}

func (q PriorityQueue) Less(i, j int) bool {
	return q[i].cost < q[j].cost
}

func (q PriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *PriorityQueue) Push(x interface{}) {
	item := x.(*State)
	*q = append(*q, item)
}

func (q *PriorityQueue) Pop() interface{} {
	n := len(*q)
	item := (*q)[n-1]
	*q = (*q)[0 : n-1]
	return item
}

type Grid [][]string

func PartOne(g Grid) int {
	start := g.find(START)
	end := g.find(END)

	queue := PriorityQueue{}
	heap.Init(&queue)
	heap.Push(&queue, &State{
		loc:  start,
		dir:  "E",
		cost: 0,
		path: []shared.Coord{start},
	})

	visited := map[string]bool{}

	for {
		if queue.Len() == 0 {
			break
		}

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

		for newDir, neighbor := range g.neighbors(current.loc) {
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

func (g Grid) draw(paths map[string][]shared.Coord) {
	for y := range g {
		for x := range g[y] {
			drawLoc := shared.Coord{X: x, Y: y}

			drawn := false

			for symbol, path := range paths {
				if slices.Contains(path, drawLoc) {
					fmt.Print(symbol)
					drawn = true
					break
				}
			}

			if !drawn {
				fmt.Print(g[y][x])
			}
		}
		fmt.Println()
	}
}

func (g Grid) find(value string) shared.Coord {

	for y, row := range g {
		for x := range row {
			if g[y][x] == value {

				return shared.Coord{X: x, Y: y}
			}
		}
	}

	return shared.Coord{X: -1, Y: -1}
}

func (g *Grid) fillDeadEnd(loc shared.Coord) {
	for {
		neighbors := (*g).neighbors(loc)

		if len(neighbors) != 1 || (*g)[loc.Y][loc.X] == START {
			break
		}

		for _, neighbor := range neighbors {
			(*g)[loc.Y][loc.X] = FILL
			loc = neighbor
		}
	}
}
