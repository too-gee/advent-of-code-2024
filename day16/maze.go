package main

import (
	"slices"

	"github.com/too-gee/advent-of-code-2024/shared"
)

const FILL = "="
const WALL = "#"
const EMPTY = "."
const START = "S"
const END = "E"

type Maze struct {
	shared.Grid
}

func (m Maze) Neighbors(loc shared.Coord) map[string]shared.Coord {
	neighbors := map[string]shared.Coord{}

	for dir, neighbor := range loc.Neighbors() {
		if m.Contains(neighbor) &&
			!slices.Contains([]string{WALL, FILL}, m.Grid[neighbor.Y][neighbor.X]) {
			neighbors[dir] = neighbor
		}
	}

	return neighbors
}

func (m Maze) IsPassable(loc shared.Coord) bool {
	return slices.Contains([]string{EMPTY, START, END}, m.Grid[loc.Y][loc.X])
}

func (m Maze) IsDeadEnd(loc shared.Coord) bool {
	if slices.Contains([]string{START, END}, m.Grid[loc.Y][loc.X]) {
		return false
	}

	return len(m.Neighbors(loc)) == 1
}

func (m Maze) FillDeadEnds() Maze {
	for y := range m.Height() {
		for x := range m.Width() {
			loc := shared.Coord{X: x, Y: y}
			if m.IsPassable(loc) && m.IsDeadEnd(loc) {
				m.FillDeadEnd(loc)
			}
		}
	}

	return m
}

func (m *Maze) FillDeadEnd(loc shared.Coord) {
	for {
		neighbors := (*m).Neighbors(loc)

		if len(neighbors) != 1 || (*m).Grid[loc.Y][loc.X] == START {
			break
		}

		for _, neighbor := range neighbors {
			(*m).Grid[loc.Y][loc.X] = FILL
			loc = neighbor
		}
	}
}
