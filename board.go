package main

const BoardWidth = 20
const BoardHeight = 20

type Board struct {
	W, H int
}

func NewBoard() *Board {
	return &Board{
		W: BoardWidth,
		H: BoardHeight,
	}
}

func (b *Board) InBounds(p Point) bool {
	return p.X >= 0 && p.X < b.W && p.Y >= 0 && p.Y < b.H
}

func (b *Board) CheckWallCollision(s *Snake) bool {
	if !s.Alive {
		return false
	}
	return !b.InBounds(s.Head())
}

type CellKind int

const (
	CellEmpty CellKind = iota
	CellSnake1Head
	CellSnake1Body
	CellSnake2Head
	CellSnake2Body
	CellFood
)

type Cell struct {
	Kind  CellKind
	Color string
}

func (b *Board) Render(s1, s2 *Snake, fm *FoodManager) [][]Cell {
	grid := make([][]Cell, b.H)
	for y := 0; y < b.H; y++ {
		grid[y] = make([]Cell, b.W)
		for x := 0; x < b.W; x++ {
			grid[y][x] = Cell{Kind: CellEmpty}
		}
	}

	if s1.Alive {
		for i, seg := range s1.Body {
			if i == 0 {
				grid[seg.Y][seg.X] = Cell{Kind: CellSnake1Head, Color: s1.Color}
			} else {
				grid[seg.Y][seg.X] = Cell{Kind: CellSnake1Body, Color: s1.Color}
			}
		}
	}
	if s2.Alive {
		for i, seg := range s2.Body {
			if i == 0 {
				grid[seg.Y][seg.X] = Cell{Kind: CellSnake2Head, Color: s2.Color}
			} else {
				grid[seg.Y][seg.X] = Cell{Kind: CellSnake2Body, Color: s2.Color}
			}
		}
	}
	for _, f := range fm.Foods {
		grid[f.Pos.Y][f.Pos.X] = Cell{Kind: CellFood}
	}

	return grid
}
