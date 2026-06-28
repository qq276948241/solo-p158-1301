package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ColorSnake1 = "#22c55e"
	ColorSnake2 = "#3b82f6"
	ColorFood   = "#ef4444"
	ColorBorder = "#64748b"
	ColorText   = "#e2e8f0"
	ColorDim    = "#94a3b8"
	ColorBg     = "#0f172a"
)

const TickInterval = 120 * time.Millisecond

type tickMsg time.Time

type UI struct {
	Game      *Game
	HighScore int
}

func NewUI(g *Game, hs int) *UI {
	return &UI{Game: g, HighScore: hs}
}

func (u *UI) cellChar(cell Cell) string {
	switch cell.Kind {
	case CellSnake1Head:
		return "●"
	case CellSnake1Body:
		return "■"
	case CellSnake2Head:
		return "●"
	case CellSnake2Body:
		return "■"
	case CellFood:
		return "◆"
	default:
		return "·"
	}
}

func (u *UI) cellStyle(cell Cell) lipgloss.Style {
	style := lipgloss.NewStyle().Width(2)
	switch cell.Kind {
	case CellSnake1Head:
		return style.Foreground(lipgloss.Color(ColorSnake1)).Bold(true)
	case CellSnake1Body:
		return style.Foreground(lipgloss.Color(ColorSnake1))
	case CellSnake2Head:
		return style.Foreground(lipgloss.Color(ColorSnake2)).Bold(true)
	case CellSnake2Body:
		return style.Foreground(lipgloss.Color(ColorSnake2))
	case CellFood:
		return style.Foreground(lipgloss.Color(ColorFood)).Bold(true)
	default:
		return style.Foreground(lipgloss.Color(ColorDim))
	}
}

func (u *UI) renderBoard() string {
	grid := u.Game.Board.Render(u.Game.Snake1, u.Game.Snake2, u.Game.FoodMgr)

	header := "┌" + strings.Repeat("──", u.Game.Board.W) + "┐\n"
	footer := "└" + strings.Repeat("──", u.Game.Board.W) + "┘"

	var sb strings.Builder
	sb.WriteString(header)

	for y := 0; y < u.Game.Board.H; y++ {
		sb.WriteString("│")
		for x := 0; x < u.Game.Board.W; x++ {
			cell := grid[y][x]
			s := u.cellStyle(cell)
			sb.WriteString(s.Render(u.cellChar(cell)))
		}
		sb.WriteString("│\n")
	}

	sb.WriteString(footer)

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorder))
	return borderStyle.Render(sb.String())
}

func (u *UI) renderSidebar() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Underline(true).
		MarginBottom(1)

	p1Style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSnake1)).
		Bold(true)

	p2Style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSnake2)).
		Bold(true)

	foodStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorFood)).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorDim))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true)

	hsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f59e0b")).
		Bold(true)

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("  SNAKE BATTLE") + "\n\n")

	sb.WriteString(p1Style.Render(" P1 (WASD)") + "\n")
	sb.WriteString(labelStyle.Render("   Score: ") + valueStyle.Render(fmt.Sprintf("%d", u.Game.Snake1.Score)) + "\n")
	sb.WriteString(labelStyle.Render("   Length: ") + valueStyle.Render(fmt.Sprintf("%d", len(u.Game.Snake1.Body))) + "\n\n")

	sb.WriteString(p2Style.Render(" P2 (Arrows)") + "\n")
	sb.WriteString(labelStyle.Render("   Score: ") + valueStyle.Render(fmt.Sprintf("%d", u.Game.Snake2.Score)) + "\n")
	sb.WriteString(labelStyle.Render("   Length: ") + valueStyle.Render(fmt.Sprintf("%d", len(u.Game.Snake2.Body))) + "\n\n")

	sb.WriteString(foodStyle.Render(" FOOD") + "\n")
	sb.WriteString(labelStyle.Render("   Count: ") + valueStyle.Render(fmt.Sprintf("%d / %d", u.Game.FoodMgr.Count(), MaxFoods)) + "\n")
	sb.WriteString(labelStyle.Render("   Refresh: 3s") + "\n\n")

	sb.WriteString(hsStyle.Render(" HIGH SCORE") + "\n")
	sb.WriteString(labelStyle.Render("   Best: ") + hsStyle.Render(fmt.Sprintf("%d", u.HighScore)) + "\n\n")

	sb.WriteString(labelStyle.Render(" Controls:") + "\n")
	sb.WriteString(labelStyle.Render("  Space = Pause") + "\n")
	sb.WriteString(labelStyle.Render("  R = Restart") + "\n")
	sb.WriteString(labelStyle.Render("  Q = Quit") + "\n")

	sideStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Background(lipgloss.Color(ColorBg)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorBorder)).
		Width(28)

	return sideStyle.Render(sb.String())
}

func (u *UI) renderOverlay() string {
	if u.Game.State == StatePlaying {
		return ""
	}

	var title, msg string
	titleStyle := lipgloss.NewStyle().Bold(true).MarginBottom(1)

	switch u.Game.State {
	case StatePaused:
		titleStyle = titleStyle.Foreground(lipgloss.Color("#f59e0b"))
		title = "⏸  PAUSED"
		msg = "Press SPACE to continue"
	case StateGameOver:
		titleStyle = titleStyle.Foreground(lipgloss.Color("#ef4444"))
		title = "💀 GAME OVER"
		var winnerText string
		if u.Game.Winner == "DRAW" {
			winnerText = "It's a DRAW!"
		} else {
			winnerText = fmt.Sprintf("%s WINS!", u.Game.Winner)
		}
		scoreText := fmt.Sprintf("P1: %d  |  P2: %d", u.Game.Snake1.Score, u.Game.Snake2.Score)
		msg = winnerText + "\n" + scoreText + "\n\nPress R to restart"
	}

	overlayStyle := lipgloss.NewStyle().
		Padding(2, 4).
		Background(lipgloss.Color(ColorBg)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorBorder)).
		Align(lipgloss.Center).
		Width(36)

	content := titleStyle.Render(title) + "\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText)).Render(msg)

	return overlayStyle.Render(content)
}

func (u *UI) View() string {
	board := u.renderBoard()
	sidebar := u.renderSidebar()

	row := lipgloss.JoinHorizontal(lipgloss.Top, board, "  ", sidebar)

	overlay := u.renderOverlay()
	if overlay != "" {
		boardHeight := strings.Count(board, "\n") + 1
		fullBoard := lipgloss.Place(
			lipgloss.Width(board),
			boardHeight,
			lipgloss.Center,
			lipgloss.Center,
			overlay,
		)
		row = lipgloss.JoinHorizontal(lipgloss.Top, fullBoard, "  ", sidebar)
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorDim)).
		MarginTop(1).
		PaddingLeft(2)
	help := helpStyle.Render("🐍 P1: WASD  |  P2: ↑↓←→  |  SPACE: Pause  |  R: Restart  |  Q: Quit")

	return lipgloss.JoinVertical(lipgloss.Left, row, help)
}

func (u *UI) Init() tea.Cmd {
	return tea.Tick(TickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		u.Game.Tick()
		return u, tea.Tick(TickInterval, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return u, tea.Quit
		}
		for _, r := range msg.Runes {
			switch r {
			case 'q', 'Q':
				return u, tea.Quit
			case 'r', 'R':
				u.Game.Reset()
			case ' ':
				u.Game.TogglePause()
			case 'w', 'W':
				u.Game.Snake1.SetDirection(DirUp)
			case 's', 'S':
				u.Game.Snake1.SetDirection(DirDown)
			case 'a', 'A':
				u.Game.Snake1.SetDirection(DirLeft)
			case 'd', 'D':
				u.Game.Snake1.SetDirection(DirRight)
			}
		}
		switch msg.Type {
		case tea.KeyUp:
			u.Game.Snake2.SetDirection(DirUp)
		case tea.KeyDown:
			u.Game.Snake2.SetDirection(DirDown)
		case tea.KeyLeft:
			u.Game.Snake2.SetDirection(DirLeft)
		case tea.KeyRight:
			u.Game.Snake2.SetDirection(DirRight)
		case tea.KeySpace:
			u.Game.TogglePause()
		}
	}
	return u, nil
}
