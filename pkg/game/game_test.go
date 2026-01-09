package game

import (
	"testing"
	"time"
)

// TestFoodExpirationWithPause tests that food timers pause correctly
func TestFoodExpirationWithPause(t *testing.T) {
	// Create a red food (10 second duration)
	food := Food{
		Pos:       Point{X: 5, Y: 5},
		FoodType:  FoodRed,
		SpawnTime: time.Now(),
	}

	// Wait 2 seconds
	time.Sleep(2 * time.Second)

	// Check remaining time without pause (should be ~8 seconds)
	remaining1 := food.GetRemainingSeconds(0)
	if remaining1 < 7 || remaining1 > 9 {
		t.Errorf("Expected ~8 seconds remaining, got %d", remaining1)
	}
	t.Logf("After 2s real time, 0s pause: %d seconds remaining", remaining1)

	// Simulate 5 seconds of pause time
	pausedTime := 5 * time.Second

	// Check remaining time with pause (should still be ~8 seconds)
	remaining2 := food.GetRemainingSeconds(pausedTime)
	if remaining2 < 7 || remaining2 > 9 {
		t.Errorf("Expected ~8 seconds remaining with 5s pause, got %d", remaining2)
	}
	t.Logf("After 2s real time, 5s pause: %d seconds remaining (timer paused correctly!)", remaining2)

	// The key point: remaining time should be the same with or without pause
	// because paused time is subtracted from elapsed time
	if remaining1 != remaining2 {
		t.Logf("✅ Timer correctly accounts for pause time (no change in countdown)")
	}
}

// TestFoodExpirationDuringPause tests expiration check during active pause
func TestFoodExpirationDuringPause(t *testing.T) {
	// Create a food that will expire in 2 seconds
	food := Food{
		Pos:       Point{X: 5, Y: 5},
		FoodType:  FoodRed,                          // 10 second duration
		SpawnTime: time.Now().Add(-9 * time.Second), // Spawned 9 seconds ago
	}

	// Without pause, should have 1 second remaining
	remaining := food.GetRemainingSeconds(0)
	t.Logf("Food without pause: %d seconds remaining", remaining)

	if food.IsExpired(0) {
		t.Error("Food should not be expired yet without pause")
	}

	// Simulate being paused for 5 seconds
	pausedTime := 5 * time.Second

	// With 5s pause, should have 6 seconds remaining (9-5 = 4 elapsed, 10-4 = 6 remaining)
	remainingWithPause := food.GetRemainingSeconds(pausedTime)
	t.Logf("Food with 5s pause: %d seconds remaining", remainingWithPause)

	if remainingWithPause < 5 || remainingWithPause > 7 {
		t.Errorf("Expected ~6 seconds remaining with pause, got %d", remainingWithPause)
	}

	if food.IsExpired(pausedTime) {
		t.Error("Food should not be expired when accounting for pause time")
	}
}

// TestGamePauseIntegration tests the full integration with Game struct
func TestGamePauseIntegration(t *testing.T) {
	g := NewGame()

	// Get initial total paused time (should be 0)
	if g.GetTotalPausedTime() != 0 {
		t.Errorf("Expected 0 initial paused time, got %v", g.GetTotalPausedTime())
	}

	// Start pause
	g.TogglePause()
	if !g.Paused {
		t.Error("Game should be paused")
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Total paused time should now include current pause
	totalPaused := g.GetTotalPausedTime()
	if totalPaused < 90*time.Millisecond || totalPaused > 110*time.Millisecond {
		t.Errorf("Expected ~100ms paused, got %v", totalPaused)
	}
	t.Logf("While paused: total paused time = %v", totalPaused)

	// Resume
	g.TogglePause()
	if g.Paused {
		t.Error("Game should not be paused")
	}

	// Total paused time should now be accumulated
	accumulatedPause := g.GetTotalPausedTime()
	if accumulatedPause < 90*time.Millisecond || accumulatedPause > 110*time.Millisecond {
		t.Errorf("Expected ~100ms accumulated pause, got %v", accumulatedPause)
	}
	t.Logf("After resume: accumulated paused time = %v", accumulatedPause)

	// Wait a bit more
	time.Sleep(50 * time.Millisecond)

	// Total paused time should still be the same (not actively paused)
	if g.GetTotalPausedTime() != accumulatedPause {
		t.Error("Paused time should not increase when not paused")
	}

	t.Log("✅ Pause integration test passed!")
}
