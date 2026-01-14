package config

import "time"

// Game board dimensions
const (
	Width        = 25 // Reduced width for emoji double-width
	Height       = 25
	GameDuration = 30 * time.Second
)

// Food spawn settings
const (
	FoodSpawnInterval = 5 * time.Second
	MaxFoodsOnBoard   = 5 // Maximum concurrent foods on board
)

// Obstacle settings
const (
	ObstacleSpawnInterval = 10 * time.Second
	ObstacleDuration      = 30 * time.Second
	MaxObstacles          = 4
)

// Fireball settings
const (
	FireballSpeed    = 48 * time.Millisecond  // Time between fireball moves
	FireballCooldown = 300 * time.Millisecond // Time between shots
)

// Speed and boost settings
const (
	BaseTick             = 16 * time.Millisecond  // Base tick interval (~60 FPS)
	NormalTicksPerUpdate = 18                     // Normal speed: 16ms * 18 = 288ms
	BoostTicksPerUpdate  = 6                      // Boost: 16ms * 6 = 96ms
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
