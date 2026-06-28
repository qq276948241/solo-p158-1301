package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

const scoreFile = "score.txt"

func loadHighScore() (int, error) {
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	} else {
		exePath = filepath.Dir(exePath)
	}
	path := filepath.Join(exePath, scoreFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	val, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, nil
	}
	return val, nil
}

func saveHighScore(score int) error {
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	} else {
		exePath = filepath.Dir(exePath)
	}
	path := filepath.Join(exePath, scoreFile)

	return os.WriteFile(path, []byte(strconv.Itoa(score)), 0644)
}

func main() {
	hs, err := loadHighScore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load high score: %v\n", err)
		hs = 0
	}

	game := NewGame()
	ui := NewUI(game, hs)

	p := tea.NewProgram(ui)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	finalUI, ok := finalModel.(*UI)
	if !ok {
		return
	}

	best := hs
	if finalUI.Game.Snake1.Score > best {
		best = finalUI.Game.Snake1.Score
	}
	if finalUI.Game.Snake2.Score > best {
		best = finalUI.Game.Snake2.Score
	}
	if best > hs {
		if err := saveHighScore(best); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not save high score: %v\n", err)
		}
	}
}
