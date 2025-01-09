package shared

import (
	"fmt"
	"math"
	"slices"
)

type Grid [][]string

func (g Grid) Width() int {
	if g == nil {
		return 0
	}

	return len(g[0])
}

func (g Grid) Height() int {
	if g == nil {
		return 0
	}

	return len(g)
}

func (g Grid) Contains(loc CoordLike) bool {
	return loc.GetX() >= 0 &&
		loc.GetX() < len(g[0]) &&
		loc.GetY() >= 0 &&
		loc.GetY() < len(g)
}

func (g Grid) LocationOf(value string) Coord {
	for y := range g {
		for x := range g[0] {
			if g[y][x] == value {
				return Coord{X: x, Y: y}
			}
		}
	}

	return Coord{X: -1, Y: -1}
}

func (g Grid) Draw(markers map[string]string, paths map[string][]Coord) {
	if g == nil {
		return
	}

	if markers == nil {
		markers = map[string]string{}
	}

	if paths == nil {
		paths = map[string][]Coord{}
	}

	yMin, yMax := g.Height()-1, 0
	xMin, xMax := g.Width()-1, 0

	for y := range g.Height() {
		for x := range g.Width() {
			if g[y][x] != "." {
				yMin = int(math.Min(float64(yMin), float64(y)))
				yMax = int(math.Max(float64(yMax), float64(y)))
				xMin = int(math.Min(float64(xMin), float64(x)))
				xMax = int(math.Max(float64(xMax), float64(x)))
			}
		}
	}

	for y := yMin - 2; y <= yMax+2; y++ {
		for x := xMin - 2; x <= xMax+2; x++ {
			loc := Coord{X: x, Y: y}

			// print the border
			if (x >= xMin-2 && x < xMin-1) ||
				(y >= yMin-2 && y < yMin-1) ||
				(x > xMax+1 && x <= xMax+2) ||
				(y > yMax+1 && y <= yMax+2) {
				fmt.Print("██")
				continue
			}

			// print a buffer
			if (x >= xMin-1 && x < xMin) ||
				(y >= yMin-1 && y < yMin) ||
				(x > xMax && x <= xMax+1) ||
				(y > yMax && y <= yMax+1) {
				fmt.Print("　")
				continue
			}

			// print the grid
			if g.Contains(loc) {
				pathPrint := false
				for cell, path := range paths {
					if slices.Contains(path, loc) {
						fmt.Print(cell)
						pathPrint = true
					}
				}
				if pathPrint {
					continue
				}

				cell, ok := markers[g[y][x]]

				if ok {
					fmt.Print(cell)
					continue
				}

				switch g[y][x] {
				case "#":
					fmt.Print("⬛")
				case ".":
					fmt.Print("　")
				default:
					fmt.Print(g[y][x])
				}
			}
		}
		fmt.Println()
	}
}
