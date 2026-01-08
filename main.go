package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	game := NewGame()

	// Input channel - sends both char and key
	type keyInput struct {
		char rune
		key  keyboard.Key
	}
	inputChan := make(chan keyInput)
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				return
			}
			inputChan <- keyInput{char: char, key: key}
		}
	}()

	ticker := time.NewTicker(baseTick)
	defer ticker.Stop()

	var (
		tickCount           = 0
		boosting            = false
		lastBoostKeyTime    time.Time
		lastDirKeyTime      time.Time // 上次按方向键的时间
		lastDirKeyDir       Point     // 上次按的方向
		consecutiveKeyCount = 0       // 连续按同方向键的次数
	)

	// 检查是否触发加速（需要连续快速按键）
	checkBoostKey := func(inputDir Point) {
		now := time.Now()

		// 检查是否是连续按同方向键
		if inputDir == lastDirKeyDir && time.Since(lastDirKeyTime) < keyRepeatWindow {
			consecutiveKeyCount++
		} else {
			// 方向变了或者间隔太长，重置计数
			consecutiveKeyCount = 1
		}

		lastDirKeyDir = inputDir
		lastDirKeyTime = now

		// 达到阈值后触发加速
		if consecutiveKeyCount >= boostThreshold && inputDir == game.direction {
			boosting = true
			lastBoostKeyTime = now
		}
	}

	game.render()

	for {
		select {
		case input := <-inputChan:
			var inputDir Point
			dirChanged := false

			// Handle arrow keys
			switch input.key {
			case keyboard.KeyArrowUp:
				inputDir = Point{x: 0, y: -1}
				if game.direction.y != 1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			case keyboard.KeyArrowDown:
				inputDir = Point{x: 0, y: 1}
				if game.direction.y != -1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			case keyboard.KeyArrowLeft:
				inputDir = Point{x: -1, y: 0}
				if game.direction.x != 1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			case keyboard.KeyArrowRight:
				inputDir = Point{x: 1, y: 0}
				if game.direction.x != -1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			}

			// Handle character keys
			switch input.char {
			case 'w', 'W':
				inputDir = Point{x: 0, y: -1}
				if game.direction.y != 1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			case 's', 'S':
				inputDir = Point{x: 0, y: 1}
				if game.direction.y != -1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			case 'a', 'A':
				inputDir = Point{x: -1, y: 0}
				if game.direction.x != 1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			case 'd', 'D':
				inputDir = Point{x: 1, y: 0}
				if game.direction.x != -1 && game.direction != inputDir {
					game.direction = inputDir
					dirChanged = true
				}
			case 'q', 'Q':
				fmt.Println("\n  Thanks for playing! 👋")
				return
			case 'r', 'R':
				if game.gameOver {
					game = NewGame()
					boosting = false
					tickCount = 0
					consecutiveKeyCount = 0
				}
			case 'p', 'P', ' ':
				if !game.gameOver {
					if !game.paused {
						// 开始暂停
						game.pauseStart = time.Now()
					} else {
						// 结束暂停，累加暂停时间
						game.pausedTime += time.Since(game.pauseStart)
					}
					game.paused = !game.paused
					game.render()
				}
			}

			// 检查是否触发加速
			if inputDir != (Point{}) {
				if dirChanged {
					// 方向改变了，重置连续按键计数
					consecutiveKeyCount = 1
					lastDirKeyDir = inputDir
					lastDirKeyTime = time.Now()
					boosting = false
				} else {
					// 按下的是当前方向，检查是否触发加速
					checkBoostKey(inputDir)
				}
			}

		case <-ticker.C:
			// 检查加速是否超时
			if boosting && time.Since(lastBoostKeyTime) > boostTimeout {
				boosting = false
			}

			tickCount++

			// 根据是否加速决定更新频率
			ticksNeeded := normalTicksPerUpdate
			if boosting {
				ticksNeeded = boostTicksPerUpdate
			}

			if tickCount >= ticksNeeded {
				tickCount = 0
				if !game.gameOver && !game.paused {
					game.update()
				}
				game.renderWithBoost(boosting)
			}
		}
	}
}
