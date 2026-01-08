package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// clearScreen 清屏
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// render 渲染游戏画面
func (g *Game) render() {
	clearScreen()

	// Cell types for coloring
	const (
		cellEmpty = iota
		cellWall
		cellHead
		cellBody
		cellCrash
	)

	// Create the board
	board := make([][]int, height)
	for i := range board {
		board[i] = make([]int, width)
	}

	// 创建豆子位置到 emoji 的映射
	foodEmojis := make(map[Point]string)
	timerEmojis := make(map[Point]string) // 倒计时数字的映射
	for _, food := range g.foods {
		foodEmojis[food.pos] = food.getEmojiWithTimer()
		// 游戏未结束时才显示倒计时数字
		if !g.gameOver {
			// 倒计时数字显示在豆子右侧一格
			timerEmoji := food.getTimerEmoji()
			if timerEmoji != "" {
				timerPos := Point{x: food.pos.x + 1, y: food.pos.y}
				timerEmojis[timerPos] = timerEmoji
			}
		}
	}

	// Draw walls
	for x := 0; x < width; x++ {
		board[0][x] = cellWall
		board[height-1][x] = cellWall
	}
	for y := 0; y < height; y++ {
		board[y][0] = cellWall
		board[y][width-1] = cellWall
	}

	// Draw snake
	for i, p := range g.snake {
		if i == 0 {
			board[p.y][p.x] = cellHead
		} else {
			board[p.y][p.x] = cellBody
		}
	}

	// Draw crash point if game over
	if g.gameOver {
		// 确保碰撞点在边界内才显示
		if g.crashPoint.x >= 0 && g.crashPoint.x < width && g.crashPoint.y >= 0 && g.crashPoint.y < height {
			board[g.crashPoint.y][g.crashPoint.x] = cellCrash
		}
	}

	// Emoji squares (these are typically rendered as perfect squares)
	const (
		charEmpty = "  " // Two spaces to match emoji width
		charWall  = "⬜"
		charHead  = "🟢"
		charBody  = "🟩"
		charCrash = "💥"
	)

	// 计算吃豆速度
	elapsed := time.Since(g.startTime) - g.pausedTime
	var eatingSpeed float64
	if elapsed.Seconds() > 0 {
		eatingSpeed = float64(g.foodEaten) / elapsed.Seconds()
	}

	// Print board
	fmt.Println("\n  🐍 SNAKE GAME 🐍")
	fmt.Printf("  Score: %d  |  吃豆速度: %.2f 个/秒  |  已吃: %d 个\n\n", g.score, eatingSpeed, g.foodEaten)
	for y, row := range board {
		fmt.Print("  ")
		for x, cell := range row {
			pos := Point{x: x, y: y}
			// 先检查是否有倒计时数字在这个位置
			if timer, hasTimer := timerEmojis[pos]; hasTimer && cell == cellEmpty {
				fmt.Print(timer + " ") // 添加空格补齐到2字符宽度
			} else if emoji, hasFood := foodEmojis[pos]; hasFood && cell == cellEmpty {
				// 然后检查是否有豆子在这个位置
				fmt.Print(emoji)
			} else {
				switch cell {
				case cellEmpty:
					fmt.Print(charEmpty)
				case cellWall:
					fmt.Print(charWall)
				case cellHead:
					fmt.Print(charHead)
				case cellBody:
					fmt.Print(charBody)
				case cellCrash:
					fmt.Print(charCrash)
				}
			}
		}
		fmt.Println()
	}
	fmt.Println("\n  Use WASD or Arrow keys to move, hold direction key to boost 🚀")
	fmt.Println("  P to pause, Q to quit")

	if g.paused {
		fmt.Println("\n  ⏸️  PAUSED - Press P to continue")
	}

	if g.gameOver {
		fmt.Println("\n  💀 GAME OVER! Press R to restart or Q to quit")
	}
}

// renderWithBoost 带加速指示器的渲染
func (g *Game) renderWithBoost(boosting bool) {
	clearScreen()

	// Cell types for coloring
	const (
		cellEmpty = iota
		cellWall
		cellHead
		cellBody
		cellCrash
	)

	// Create the board
	board := make([][]int, height)
	for i := range board {
		board[i] = make([]int, width)
	}

	// 创建豆子位置到 emoji 的映射
	foodEmojis := make(map[Point]string)
	timerEmojis := make(map[Point]string) // 倒计时数字的映射
	for _, food := range g.foods {
		foodEmojis[food.pos] = food.getEmojiWithTimer()
		// 游戏未结束时才显示倒计时数字
		if !g.gameOver {
			// 倒计时数字显示在豆子右侧一格
			timerEmoji := food.getTimerEmoji()
			if timerEmoji != "" {
				timerPos := Point{x: food.pos.x + 1, y: food.pos.y}
				timerEmojis[timerPos] = timerEmoji
			}
		}
	}

	// Draw walls
	for x := 0; x < width; x++ {
		board[0][x] = cellWall
		board[height-1][x] = cellWall
	}
	for y := 0; y < height; y++ {
		board[y][0] = cellWall
		board[y][width-1] = cellWall
	}

	// Draw snake
	for i, p := range g.snake {
		if i == 0 {
			board[p.y][p.x] = cellHead
		} else {
			board[p.y][p.x] = cellBody
		}
	}

	// Draw crash point if game over
	if g.gameOver {
		// 确保碰撞点在边界内才显示
		if g.crashPoint.x >= 0 && g.crashPoint.x < width && g.crashPoint.y >= 0 && g.crashPoint.y < height {
			board[g.crashPoint.y][g.crashPoint.x] = cellCrash
		}
	}

	// Emoji squares (these are typically rendered as perfect squares)
	const (
		charEmpty = "  " // Two spaces to match emoji width
		charWall  = "⬜"
		charHead  = "🟢"
		charBody  = "🟩"
		charCrash = "💥"
	)

	// 计算吃豆速度
	elapsed := time.Since(g.startTime) - g.pausedTime
	var eatingSpeed float64
	if elapsed.Seconds() > 0 {
		eatingSpeed = float64(g.foodEaten) / elapsed.Seconds()
	}

	// Print board with boost indicator
	fmt.Println("\n  🐍 SNAKE GAME 🐍")
	if boosting {
		fmt.Printf("  Score: %d  |  吃豆速度: %.2f 个/秒  |  已吃: %d 个  |  🚀 BOOST!\n\n", g.score, eatingSpeed, g.foodEaten)
	} else {
		fmt.Printf("  Score: %d  |  吃豆速度: %.2f 个/秒  |  已吃: %d 个\n\n", g.score, eatingSpeed, g.foodEaten)
	}
	for y, row := range board {
		fmt.Print("  ")
		for x, cell := range row {
			pos := Point{x: x, y: y}
			// 先检查是否有倒计时数字在这个位置
			if timer, hasTimer := timerEmojis[pos]; hasTimer && cell == cellEmpty {
				fmt.Print(timer + " ") // 添加空格补齐到2字符宽度
			} else if emoji, hasFood := foodEmojis[pos]; hasFood && cell == cellEmpty {
				// 然后检查是否有豆子在这个位置
				fmt.Print(emoji)
			} else {
				switch cell {
				case cellEmpty:
					fmt.Print(charEmpty)
				case cellWall:
					fmt.Print(charWall)
				case cellHead:
					fmt.Print(charHead)
				case cellBody:
					fmt.Print(charBody)
				case cellCrash:
					fmt.Print(charCrash)
				}
			}
		}
		fmt.Println()
	}
	fmt.Println("\n  Use WASD or Arrow keys to move, hold direction key to boost 🚀")
	fmt.Println("  P to pause, Q to quit")

	if g.paused {
		fmt.Println("\n  ⏸️  PAUSED - Press P to continue")
	}

	if g.gameOver {
		fmt.Println("\n  💀 GAME OVER! Press R to restart or Q to quit")
	}
}
