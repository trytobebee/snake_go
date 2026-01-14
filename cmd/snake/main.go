package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/trytobebee/snake_go/pkg/config"
	"github.com/trytobebee/snake_go/pkg/game"
	"github.com/trytobebee/snake_go/pkg/input"
	"github.com/trytobebee/snake_go/pkg/renderer"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize input handler
	inputHandler := input.NewKeyboardHandler()
	if err := inputHandler.Start(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer inputHandler.Stop()

	// Initialize renderer
	render := renderer.NewTerminalRenderer()

	// Create new game
	g := game.NewGame()

	// Get input channel
	inputChan := inputHandler.GetInputChan()

	// Game loop ticker
	ticker := time.NewTicker(config.BaseTick)
	defer ticker.Stop()

	// Boost tracking state
	var (
		tickCount           = 0
		boosting            = false
		lastBoostKeyTime    time.Time
		lastDirKeyTime      time.Time
		lastDirKeyDir       game.Point
		consecutiveKeyCount = 0
		aiTickCount         = 0
	)

	// Check if boost should be triggered
	checkBoostKey := func(inputDir game.Point) {
		now := time.Now()

		// Check if same direction pressed consecutively
		if inputDir == lastDirKeyDir && time.Since(lastDirKeyTime) < config.KeyRepeatWindow {
			consecutiveKeyCount++
		} else {
			// Direction changed or gap too long, reset counter
			consecutiveKeyCount = 1
		}

		lastDirKeyDir = inputDir
		lastDirKeyTime = now

		// Trigger boost if threshold reached and same as current direction
		if consecutiveKeyCount >= config.BoostThreshold && inputDir == g.Direction {
			boosting = true
			lastBoostKeyTime = now
		}
	}

	// Initial render
	render.Render(g, false)

	// Main game loop
	for {
		select {
		case inputEvent := <-inputChan:
			// Handle quit
			if input.IsQuit(inputEvent) {
				fmt.Println("\n  Thanks for playing! 👋")
				return
			}

			// Handle restart
			if input.IsRestart(inputEvent) {
				if g.GameOver {
					g = game.NewGame()
					boosting = false
					tickCount = 0
					consecutiveKeyCount = 0
				}
			}

			// Handle pause
			if input.IsPause(inputEvent) {
				if !g.GameOver {
					g.TogglePause()
					render.Render(g, boosting)
				}
			}

			// Handle direction input
			if inputDir, isValid := input.ParseDirection(inputEvent); isValid {
				dirChanged := g.SetDirection(inputDir)

				if dirChanged {
					// Direction changed, reset boost
					consecutiveKeyCount = 1
					lastDirKeyDir = inputDir
					lastDirKeyTime = time.Now()
					boosting = false
				} else {
					// Same direction, check for boost
					checkBoostKey(inputDir)
				}
			}

		case <-ticker.C:
			// Check boost timeout
			if boosting && time.Since(lastBoostKeyTime) > config.BoostTimeout {
				boosting = false
			}

			tickCount++

			// Determine update frequency based on boost
			ticksNeeded := config.NormalTicksPerUpdate
			if boosting {
				ticksNeeded = config.BoostTicksPerUpdate
			}

			if tickCount >= ticksNeeded {
				tickCount = 0
				if !g.GameOver && !g.Paused {
					g.Update()
				}
				render.Render(g, boosting)
			}

			// AI Tick
			aiTickCount++
			aiTicksNeeded := 13
			if g.AIBoosting {
				aiTicksNeeded = 4
			}
			if aiTickCount >= aiTicksNeeded {
				aiTickCount = 0
				if !g.GameOver && !g.Paused {
					g.UpdateAISnake()
				}
				render.Render(g, boosting)
			}
		}
	}
}
