package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	// part 1
	maze := readInput(fileName)

	for y := range maze {
		for x := range maze[y] {
			loc := Coord{x: x, y: y}
			if maze.isPassable(loc) && maze.isDeadEnd(loc) {
				maze.fillDeadEnd(loc)
			}
		}
	}

	maze.draw(map[string][]Coord{})
	cost, path := maze.solve(maze.find(START), maze.find(END))
	fmt.Printf("Lowest cost: %d\n", cost)
	fmt.Printf("Path: %v\n", path)
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

func (g Grid) inGrid(loc Coord) bool {
	return loc.y >= 0 &&
		loc.y < len(g) &&
		loc.x >= 0 &&
		loc.x < len(g[loc.y])
}

func (g Grid) cellType(loc Coord) string {
	return g[loc.y][loc.x]
}

func (g Grid) neighbors(loc Coord) map[string]Coord {
	neighbors := map[string]Coord{}

	for dir, neighbor := range loc.neighbors() {
		if g.inGrid(neighbor) &&
			g.cellType(neighbor) != WALL &&
			g.cellType(neighbor) != FILL {
			neighbors[dir] = neighbor
		}
	}

	return neighbors
}

func (g Grid) isPassable(loc Coord) bool {
	return g.cellType(loc) == EMPTY ||
		g.cellType(loc) == START ||
		g.cellType(loc) == END
}

func (g Grid) isDeadEnd(loc Coord) bool {
	if g.cellType(loc) == START ||
		g.cellType(loc) == END {
		return false
	}
	options := len(g.neighbors(loc))

	return options == 1
}

type Coord struct {
	x int
	y int
}

func (c Coord) neighbors() map[string]Coord {
	return map[string]Coord{
		"W": {x: c.x - 1, y: c.y},
		"E": {x: c.x + 1, y: c.y},
		"N": {x: c.x, y: c.y - 1},
		"S": {x: c.x, y: c.y + 1},
	}
}

type State struct {
	loc  Coord
	dir  string
	cost int
	path []Coord
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

func (g Grid) solve(start, end Coord) (int, []Coord) {
	queue := PriorityQueue{}
	heap.Init(&queue)
	heap.Push(&queue, &State{
		loc:  start,
		dir:  "E",
		cost: 0,
		path: []Coord{start},
	})

	visited := map[string]bool{}

	for {
		if queue.Len() == 0 {
			break
		}

		current := heap.Pop(&queue).(*State)

		if current.loc == end {
			return current.cost, current.path
		}

		id := fmt.Sprintf("%d.%d.%s", current.loc.x, current.loc.y, current.dir)

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

	return -1, nil
}

func (g Grid) draw(paths map[string][]Coord) {
	for y := range g {
		for x := range g[y] {
			drawLoc := Coord{x: x, y: y}

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

func (g Grid) find(value string) Coord {

	for y, row := range g {
		for x := range row {
			if g[y][x] == value {

				return Coord{x: x, y: y}
			}
		}
	}

	return Coord{x: -1, y: -1}
}

func (g *Grid) fillDeadEnd(loc Coord) {
	for {
		neighbors := (*g).neighbors(loc)

		if len(neighbors) != 1 {
			break
		}

		for _, neighbor := range neighbors {
			(*g)[loc.y][loc.x] = FILL
			loc = neighbor
		}
	}
}
