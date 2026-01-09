package renderer

import (
	"fmt"
	"strings"

	"github.com/trytobebee/snake_go/pkg/config"
	"github.com/trytobebee/snake_go/pkg/game"
)

// TerminalRenderer handles terminal-based rendering
type TerminalRenderer struct {
	board  [][]int
	buffer strings.Builder
}

// Cell types for the board
const (
	cellEmpty = iota
	cellWall
	cellHead
	cellBody
	cellCrash
)

// NewTerminalRenderer creates a new terminal renderer
func NewTerminalRenderer() *TerminalRenderer {
	// Pre-allocate board to reduce GC pressure
	board := make([][]int, config.Height)
	for i := range board {
		board[i] = make([]int, config.Width)
	}

	return &TerminalRenderer{
		board: board,
	}
}

// clearScreen clears the terminal using ANSI escape codes
func (r *TerminalRenderer) clearScreen() {
	// Use multiple ANSI codes for maximum compatibility:
	// \033[?25l - Hide cursor (optional, prevents flickering)
	// \033[H - Move cursor to home position (1,1)
	// \033[2J - Clear entire screen
	// \033[3J - Clear scrollback buffer (prevent screen drift)
	fmt.Print("\033[H\033[2J\033[3J")
}

// ShowCursor shows the cursor (call on exit)
func (r *TerminalRenderer) ShowCursor() {
	fmt.Print("\033[?25h")
}

// HideCursor hides the cursor (call on start)
func (r *TerminalRenderer) HideCursor() {
	fmt.Print("\033[?25l")
}

// Render renders the game state to the terminal
func (r *TerminalRenderer) Render(g *game.Game, boosting bool) {
	r.clearScreen()
	r.buffer.Reset()

	// Reset board
	for y := range r.board {
		for x := range r.board[y] {
			r.board[y][x] = cellEmpty
		}
	}

	// Create food emoji maps
	foodEmojis := make(map[game.Point]string)
	timerEmojis := make(map[game.Point]string)
	for _, food := range g.Foods {
		foodEmojis[food.Pos] = food.GetEmojiWithTimer(config.Width, config.Height)
		// Show timer only when game is not over
		if !g.GameOver {
			timerEmoji := food.GetTimerEmoji(g.GetTotalPausedTime())
			if timerEmoji != "" {
				timerPos := game.Point{X: food.Pos.X + 1, Y: food.Pos.Y}
				timerEmojis[timerPos] = timerEmoji
			}
		}
	}

	// Draw walls
	for x := 0; x < config.Width; x++ {
		r.board[0][x] = cellWall
		r.board[config.Height-1][x] = cellWall
	}
	for y := 0; y < config.Height; y++ {
		r.board[y][0] = cellWall
		r.board[y][config.Width-1] = cellWall
	}

	// Draw snake
	for i, p := range g.Snake {
		if i == 0 {
			r.board[p.Y][p.X] = cellHead
		} else {
			r.board[p.Y][p.X] = cellBody
		}
	}

	// Draw crash point if game over
	if g.GameOver {
		if g.CrashPoint.X >= 0 && g.CrashPoint.X < config.Width &&
			g.CrashPoint.Y >= 0 && g.CrashPoint.Y < config.Height {
			r.board[g.CrashPoint.Y][g.CrashPoint.X] = cellCrash
		}
	}

	// Build output using string builder (single write is faster)
	r.buffer.WriteString("\n  🐍 SNAKE GAME 🐍\n")

	// Header with stats
	if boosting {
		r.buffer.WriteString(fmt.Sprintf("  Score: %d  |  吃豆速度: %.2f 个/秒  |  已吃: %d 个  |  🚀 BOOST!\n",
			g.Score, g.GetEatingSpeed(), g.FoodEaten))
	} else {
		r.buffer.WriteString(fmt.Sprintf("  Score: %d  |  吃豆速度: %.2f 个/秒  |  已吃: %d 个\n",
			g.Score, g.GetEatingSpeed(), g.FoodEaten))
	}

	// Show congratulatory message above the game board (always reserve the line)
	if msg := g.GetMessage(); msg != "" {
		r.buffer.WriteString("  " + msg + "\n")
	} else {
		r.buffer.WriteString("\n") // Empty line when no message
	}
	r.buffer.WriteString("\n") // Extra blank line before board

	// Render board
	for y, row := range r.board {
		r.buffer.WriteString("  ")
		for x, cell := range row {
			pos := game.Point{X: x, Y: y}

			// Check for timer emoji first
			if timer, hasTimer := timerEmojis[pos]; hasTimer && cell == cellEmpty {
				r.buffer.WriteString(timer)
				r.buffer.WriteString(" ") // Add space to match 2-char width
			} else if emoji, hasFood := foodEmojis[pos]; hasFood && cell == cellEmpty {
				// Then check for food emoji
				r.buffer.WriteString(emoji)
			} else {
				// Finally, render board cell
				switch cell {
				case cellEmpty:
					r.buffer.WriteString(config.CharEmpty)
				case cellWall:
					r.buffer.WriteString(config.CharWall)
				case cellHead:
					r.buffer.WriteString(config.CharHead)
				case cellBody:
					r.buffer.WriteString(config.CharBody)
				case cellCrash:
					r.buffer.WriteString(config.CharCrash)
				}
			}
		}
		r.buffer.WriteString("\n")
	}

	// Instructions
	r.buffer.WriteString("\n  Use WASD or Arrow keys to move, hold direction key to boost 🚀\n")
	r.buffer.WriteString("  P to pause, Q to quit\n")

	// Game state messages
	if g.Paused {
		r.buffer.WriteString("\n  ⏸️  PAUSED - Press P to continue\n")
	}

	if g.GameOver {
		r.buffer.WriteString("\n  💀 GAME OVER! Press R to restart or Q to quit\n")
	}

	// Single write to stdout (much faster than multiple fmt.Print calls)
	fmt.Print(r.buffer.String())
}
