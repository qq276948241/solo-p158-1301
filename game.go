package main

type GameState int

const (
	StatePlaying GameState = iota
	StatePaused
	StateGameOver
)

const FoodRefreshTicks = 25

type Game struct {
	Board          *Board
	Snake1         *Snake
	Snake2         *Snake
	FoodMgr        *FoodManager
	State          GameState
	Winner         string
	FoodTickCount  int
}

func NewGame() *Game {
	board := NewBoard()
	s1 := NewSnake(Point{X: 4, Y: 2}, DirRight, "green")
	s2 := NewSnake(Point{X: 15, Y: 17}, DirLeft, "blue")
	fm := NewFoodManager(board.W, board.H)
	fm.SpawnAll(s1, s2)

	return &Game{
		Board:         board,
		Snake1:        s1,
		Snake2:        s2,
		FoodMgr:       fm,
		State:         StatePlaying,
		FoodTickCount: 0,
	}
}

func (g *Game) Reset() {
	s1 := NewSnake(Point{X: 4, Y: 2}, DirRight, "green")
	s2 := NewSnake(Point{X: 15, Y: 17}, DirLeft, "blue")
	fm := NewFoodManager(g.Board.W, g.Board.H)
	fm.SpawnAll(s1, s2)
	g.Snake1 = s1
	g.Snake2 = s2
	g.FoodMgr = fm
	g.State = StatePlaying
	g.Winner = ""
	g.FoodTickCount = 0
}

func (g *Game) TogglePause() {
	if g.State == StatePlaying {
		g.State = StatePaused
	} else if g.State == StatePaused {
		g.State = StatePlaying
	}
}

func (g *Game) checkFoodCollision(s *Snake) bool {
	head := s.Head()
	food := g.FoodMgr.GetFoodAt(head)
	if food == nil {
		return false
	}
	g.FoodMgr.RemoveAt(head)
	food.ApplyEffect(s)
	g.FoodMgr.Spawn(g.Snake1, g.Snake2)
	return true
}

func (g *Game) Tick() {
	if g.State != StatePlaying {
		return
	}

	g.Snake1.Move()
	g.Snake2.Move()

	g.checkFoodCollision(g.Snake1)
	g.checkFoodCollision(g.Snake2)

	s1Dead := !g.Snake1.Alive
	s2Dead := !g.Snake2.Alive

	if !s1Dead && g.Board.CheckWallCollision(g.Snake1) {
		s1Dead = true
	}
	if !s2Dead && g.Board.CheckWallCollision(g.Snake2) {
		s2Dead = true
	}

	if !s1Dead && g.Snake1.CollidesWithSelf() {
		s1Dead = true
	}
	if !s2Dead && g.Snake2.CollidesWithSelf() {
		s2Dead = true
	}

	if !s1Dead && g.Snake1.CollidesWith(g.Snake2) {
		s1Dead = true
	}
	if !s2Dead && g.Snake2.CollidesWith(g.Snake1) {
		s2Dead = true
	}

	if s1Dead {
		g.Snake1.Alive = false
	}
	if s2Dead {
		g.Snake2.Alive = false
	}

	if s1Dead || s2Dead {
		g.State = StateGameOver
		if s1Dead && s2Dead {
			if g.Snake1.Score > g.Snake2.Score {
				g.Winner = "P1"
			} else if g.Snake2.Score > g.Snake1.Score {
				g.Winner = "P2"
			} else {
				g.Winner = "DRAW"
			}
		} else if s1Dead {
			g.Winner = "P2"
		} else {
			g.Winner = "P1"
		}
	}

	g.FoodTickCount++
	if g.FoodTickCount >= FoodRefreshTicks {
		g.FoodTickCount = 0
		g.FoodMgr.Refresh(g.Snake1, g.Snake2)
	}
}
