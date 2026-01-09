package config

import "time"

// Game board dimensions
const (
	Width  = 25 // Reduced width for emoji double-width
	Height = 25
)

// Food spawn settings
const (
	FoodSpawnInterval = 5 * time.Second
	MaxFoodsOnBoard   = 5 // Maximum concurrent foods on board
)

// Speed and boost settings
const (
	BaseTick             = 16 * time.Millisecond  // Base tick interval (~60 FPS)
	NormalTicksPerUpdate = 9                      // Normal speed: 16ms * 9 = 144ms (~150ms)
	BoostTicksPerUpdate  = 3                      // Boost: 16ms * 3 = 48ms (~50ms)
	BoostTimeout         = 150 * time.Millisecond // Boost timeout duration
	BoostThreshold       = 2                      // Consecutive key presses to trigger boost
	KeyRepeatWindow      = 200 * time.Millisecond // Time window for consecutive key detection
)

// Emoji characters for rendering
const (
	CharEmpty = "  " // Two spaces to match emoji width
	CharWall  = "⬜"
	CharHead  = "🟢"
	CharBody  = "🟩"
	CharCrash = "💥"
)
