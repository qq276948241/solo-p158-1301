package main

type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

type Point struct {
	X, Y int
}

type Snake struct {
	Body      []Point
	Dir       Direction
	NextDir   Direction
	Color     string
	Alive     bool
	Score     int
	GrowQueue int
}

func NewSnake(start Point, dir Direction, color string) *Snake {
	body := make([]Point, 3)
	for i := 0; i < 3; i++ {
		body[i] = Point{start.X, start.Y + i}
	}
	return &Snake{
		Body:      body,
		Dir:       dir,
		NextDir:   dir,
		Color:     color,
		Alive:     true,
		Score:     0,
		GrowQueue: 0,
	}
}

func (s *Snake) Head() Point {
	return s.Body[0]
}

func (s *Snake) SetDirection(d Direction) {
	if (s.Dir == DirUp && d == DirDown) ||
		(s.Dir == DirDown && d == DirUp) ||
		(s.Dir == DirLeft && d == DirRight) ||
		(s.Dir == DirRight && d == DirLeft) {
		return
	}
	s.NextDir = d
}

func (s *Snake) Move() {
	if !s.Alive {
		return
	}
	s.Dir = s.NextDir
	head := s.Head()
	var newHead Point
	switch s.Dir {
	case DirUp:
		newHead = Point{head.X, head.Y - 1}
	case DirDown:
		newHead = Point{head.X, head.Y + 1}
	case DirLeft:
		newHead = Point{head.X - 1, head.Y}
	case DirRight:
		newHead = Point{head.X + 1, head.Y}
	}
	s.Body = append([]Point{newHead}, s.Body...)
	if s.GrowQueue > 0 {
		s.GrowQueue--
	} else {
		s.Body = s.Body[:len(s.Body)-1]
	}
}

func (s *Snake) Grow() {
	s.GrowQueue++
}

func (s *Snake) Shrink(n int) {
	for i := 0; i < n; i++ {
		if len(s.Body) > 1 {
			s.Body = s.Body[:len(s.Body)-1]
		}
	}
}

func (s *Snake) Occupies(p Point) bool {
	for _, seg := range s.Body {
		if seg.X == p.X && seg.Y == p.Y {
			return true
		}
	}
	return false
}

func (s *Snake) CollidesWithSelf() bool {
	if !s.Alive || len(s.Body) < 2 {
		return false
	}
	head := s.Head()
	for i := 1; i < len(s.Body); i++ {
		if s.Body[i].X == head.X && s.Body[i].Y == head.Y {
			return true
		}
	}
	return false
}

func (s *Snake) CollidesWith(other *Snake) bool {
	if !s.Alive || !other.Alive {
		return false
	}
	head := s.Head()
	for _, seg := range other.Body {
		if seg.X == head.X && seg.Y == head.Y {
			return true
		}
	}
	return false
}
