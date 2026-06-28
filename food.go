package main

import (
	"math/rand"
	"time"
)

const MaxFoods = 5

type Food struct {
	Pos Point
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

func (fm *FoodManager) Spawn(snakes ...*Snake) bool {
	if len(fm.Foods) >= MaxFoods {
		return false
	}
	for attempts := 0; attempts < 200; attempts++ {
		p := fm.randomPoint()
		if !fm.occupied(p, snakes...) {
			fm.Foods = append(fm.Foods, Food{Pos: p})
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

func (fm *FoodManager) Refresh(snakes ...*Snake) {
	fm.Foods = fm.Foods[:0]
	fm.SpawnAll(snakes...)
}

func (fm *FoodManager) Count() int {
	return len(fm.Foods)
}
