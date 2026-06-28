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

type FoodEffect func(*Snake)

type FoodConfig struct {
	Type   FoodType
	Name   string
	Color  string
	Weight float64
	Effect FoodEffect
}

type Food struct {
	Pos  Point
	Type FoodType
}

func (f Food) Color() string {
	return foodConfigs[f.Type].Color
}

func (f Food) ApplyEffect(s *Snake) {
	foodConfigs[f.Type].Effect(s)
}

func normalEffect(s *Snake) {
	s.Grow()
	s.Score += 10
}

func goldenEffect(s *Snake) {
	s.Grow()
	s.Grow()
	s.Score += 50
}

func poisonEffect(s *Snake) {
	if !s.Shrink(2) {
		s.Alive = false
		return
	}
	s.Score -= 20
	if s.Score < 0 {
		s.Score = 0
	}
}

var foodConfigs = map[FoodType]FoodConfig{
	FoodNormal: {
		Type:   FoodNormal,
		Name:   "Normal",
		Color:  "#ef4444",
		Weight: 0.75,
		Effect: normalEffect,
	},
	FoodGolden: {
		Type:   FoodGolden,
		Name:   "Golden",
		Color:  "#f59e0b",
		Weight: 0.10,
		Effect: goldenEffect,
	},
	FoodPoison: {
		Type:   FoodPoison,
		Name:   "Poison",
		Color:  "#a855f7",
		Weight: 0.15,
		Effect: poisonEffect,
	},
}

type FoodManager struct {
	Foods    []Food
	BoardW   int
	BoardH   int
	rng      *rand.Rand
	weighted []FoodType
	totalW   float64
}

func NewFoodManager(w, h int) *FoodManager {
	weighted := make([]FoodType, 0, len(foodConfigs))
	var totalW float64
	for _, cfg := range foodConfigs {
		weighted = append(weighted, cfg.Type)
		totalW += cfg.Weight
	}
	return &FoodManager{
		Foods:    make([]Food, 0, MaxFoods),
		BoardW:   w,
		BoardH:   h,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		weighted: weighted,
		totalW:   totalW,
	}
}

func (fm *FoodManager) rollType() FoodType {
	r := fm.rng.Float64() * fm.totalW
	var acc float64
	for _, t := range fm.weighted {
		acc += foodConfigs[t].Weight
		if r < acc {
			return t
		}
	}
	return FoodNormal
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
			fm.Foods = append(fm.Foods, Food{Pos: p, Type: fm.rollType()})
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
		for i := range fm.Foods {
			if fm.Foods[i].Type != FoodNormal {
				fm.Foods[i] = Food{Pos: fm.Foods[i].Pos, Type: FoodNormal}
				return
			}
		}
	}
	for attempts := 0; attempts < 200; attempts++ {
		p := fm.randomPoint()
		if !fm.occupied(p, snakes...) {
			fm.Foods = append(fm.Foods, Food{Pos: p, Type: FoodNormal})
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
