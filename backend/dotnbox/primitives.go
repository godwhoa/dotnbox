package dotnbox

import (
	"fmt"
)

// Point represents a dot on a grid
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Point) String() string {
	return fmt.Sprintf("%d-%d", p.X, p.Y)
}

func (p Point) Add(x, y int) Point {
	return Point{X: p.X + x, Y: p.Y + y}
}

// InBounds does a bounds check on a point
func (p Point) InBounds(M, N int) bool {
	return p.X >= 0 && p.X <= N && p.Y >= 0 && p.Y <= M
}

// Line represents a line between two points on a grid
type Line struct {
	From Point `json:"from"`
	To   Point `json:"to"`
}

func (l Line) String() string {
	return fmt.Sprintf("from-%s-to-%s", l.From, l.To)
}

// Ordered returns a new Line with a determinstic ordering of from and to
func (l Line) Ordered() Line {
	if l.From.X > l.To.X {
		return Line{
			From: l.To,
			To:   l.From,
		}
	}
	return l
}

// IsValid does a bounds check on both the points and
// ensures the delta is only one
func (l Line) IsValid(M, N int) bool {
	if !l.From.InBounds(M, N) || !l.To.InBounds(M, N) {
		return false
	}
	deltax := l.From.X - l.To.X
	deltay := l.From.Y - l.To.Y
	return abs(deltax+deltay) == 1
}

// FindEdges returns all the edges that make up a box at the specified origin
func FindEdges(origin Point) []Line {
	topleft := origin
	topright := origin.Add(1, 0)
	bottomright := topright.Add(0, 1)
	bottomleft := topleft.Add(0, 1)
	return []Line{
		Line{From: topleft, To: topright}.Ordered(),
		Line{From: topright, To: bottomright}.Ordered(),
		Line{From: bottomleft, To: bottomright}.Ordered(),
		Line{From: topleft, To: bottomleft}.Ordered(),
	}
}

// Boxes returns a list of origin points of all the boxes
func Boxes(M, N int) []Point {
	var origins []Point
	for y := 0; y < N+1; y++ {
		for x := 0; x < M+1; x++ {
			origins = append(origins, Point{X: x, Y: y})
		}
	}
	return origins
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
