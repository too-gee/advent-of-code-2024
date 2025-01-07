package shared

type Coord struct {
	X int
	Y int
}

func (c Coord) Neighbors() map[string]Coord {
	return map[string]Coord{
		"W": {X: c.X - 1, Y: c.Y},
		"E": {X: c.X + 1, Y: c.Y},
		"N": {X: c.X, Y: c.Y - 1},
		"S": {X: c.X, Y: c.Y + 1},
	}
}

func (c Coord) GetX() int { return c.X }
func (c Coord) GetY() int { return c.Y }

type CoordLike interface {
	GetX() int
	GetY() int
}
