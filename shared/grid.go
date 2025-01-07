package shared

type Grid [][]string

func (g Grid) Width() int {
	return len(g[0])
}

func (g Grid) Height() int {
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
