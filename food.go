package main

import (
	"math/rand"
	"time"
)

const MaxFoods = 5

type FoodType int

const (
	FoodNormal FoodType = iota
	FoodGolden
	FoodPoison
)

type Food struct {
	Pos   Point
	Type  FoodType
	Value int
	Color string
}

func newFood(p Point, ft FoodType) Food {
	switch ft {
	case FoodGolden:
		return Food{Pos: p, Type: FoodGolden, Value: 50, Color: "#f59e0b"}
	case FoodPoison:
		return Food{Pos: p, Type: FoodPoison, Value: -20, Color: "#a855f7"}
	default:
		return Food{Pos: p, Type: FoodNormal, Value: 10, Color: "#ef4444"}
	}
}

func (f Food) ApplyEffect(s *Snake) {
	switch f.Type {
	case FoodGolden:
		s.Grow()
		s.Grow()
		s.Score += f.Value
	case FoodPoison:
		if len(s.Body) <= 3 {
			s.Alive = false
			return
		}
		s.Shrink(2)
		s.Score += f.Value
		if s.Score < 0 {
			s.Score = 0
		}
	default:
		s.Grow()
		s.Score += f.Value
	}
}

type FoodManager struct {
	Foods  []Food
	BoardW int
	BoardH int
	rng    *rand.Rand
}

func NewFoodManager(w, h int) *FoodManager {
	return &FoodManager{
		Foods:  make([]Food, 0, MaxFoods),
		BoardW: w,
		BoardH: h,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (fm *FoodManager) occupied(p Point, snakes ...*Snake) bool {
	for _, f := range fm.Foods {
		if f.Pos.X == p.X && f.Pos.Y == p.Y {
			return true
		}
	}
	for _, s := range snakes {
		if s.Occupies(p) {
			return true
		}
	}
	return false
}

func (fm *FoodManager) randomPoint() Point {
	return Point{
		X: fm.rng.Intn(fm.BoardW),
		Y: fm.rng.Intn(fm.BoardH),
	}
}

func (fm *FoodManager) rollFoodType() FoodType {
	r := fm.rng.Float64()
	switch {
	case r < 0.10:
		return FoodGolden
	case r < 0.25:
		return FoodPoison
	default:
		return FoodNormal
	}
}

func (fm *FoodManager) Spawn(snakes ...*Snake) bool {
	if len(fm.Foods) >= MaxFoods {
		return false
	}
	for attempts := 0; attempts < 200; attempts++ {
		p := fm.randomPoint()
		if !fm.occupied(p, snakes...) {
			ft := fm.rollFoodType()
			fm.Foods = append(fm.Foods, newFood(p, ft))
			return true
		}
	}
	return false
}

func (fm *FoodManager) SpawnAll(snakes ...*Snake) {
	for len(fm.Foods) < MaxFoods {
		if !fm.Spawn(snakes...) {
			break
		}
	}
	fm.ensureNormalFood(snakes...)
}

func (fm *FoodManager) ensureNormalFood(snakes ...*Snake) {
	for _, f := range fm.Foods {
		if f.Type == FoodNormal {
			return
		}
	}
	if len(fm.Foods) >= MaxFoods {
		for i, f := range fm.Foods {
			if f.Type != FoodNormal {
				fm.Foods[i] = newFood(f.Pos, FoodNormal)
				return
			}
		}
	}
	for attempts := 0; attempts < 200; attempts++ {
		p := fm.randomPoint()
		if !fm.occupied(p, snakes...) {
			fm.Foods = append(fm.Foods, newFood(p, FoodNormal))
			return
		}
	}
}

func (fm *FoodManager) RemoveAt(p Point) bool {
	for i, f := range fm.Foods {
		if f.Pos.X == p.X && f.Pos.Y == p.Y {
			fm.Foods = append(fm.Foods[:i], fm.Foods[i+1:]...)
			return true
		}
	}
	return false
}

func (fm *FoodManager) HasFoodAt(p Point) bool {
	for _, f := range fm.Foods {
		if f.Pos.X == p.X && f.Pos.Y == p.Y {
			return true
		}
	}
	return false
}

func (fm *FoodManager) GetFoodAt(p Point) *Food {
	for i := range fm.Foods {
		if fm.Foods[i].Pos.X == p.X && fm.Foods[i].Pos.Y == p.Y {
			return &fm.Foods[i]
		}
	}
	return nil
}

func (fm *FoodManager) Refresh(snakes ...*Snake) {
	fm.Foods = fm.Foods[:0]
	fm.SpawnAll(snakes...)
}

func (fm *FoodManager) Count() int {
	return len(fm.Foods)
}
