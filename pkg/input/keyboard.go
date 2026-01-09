package input

import (
	"github.com/eiannone/keyboard"
	"github.com/trytobebee/snake_go/pkg/game"
)

// KeyboardHandler handles keyboard input
type KeyboardHandler struct {
	inputChan chan KeyInput
}

// KeyInput represents a keyboard input event
type KeyInput struct {
	Char rune
	Key  keyboard.Key
}

// NewKeyboardHandler creates a new keyboard input handler
func NewKeyboardHandler() *KeyboardHandler {
	return &KeyboardHandler{
		inputChan: make(chan KeyInput),
	}
}

// Start begins listening for keyboard input
func (h *KeyboardHandler) Start() error {
	if err := keyboard.Open(); err != nil {
		return err
	}

	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				return
			}
			h.inputChan <- KeyInput{Char: char, Key: key}
		}
	}()

	return nil
}

// Stop stops the keyboard handler
func (h *KeyboardHandler) Stop() {
	keyboard.Close()
}

// GetInputChan returns the input channel
func (h *KeyboardHandler) GetInputChan() <-chan KeyInput {
	return h.inputChan
}

// ParseDirection parses a key input into a direction
func ParseDirection(input KeyInput) (dir game.Point, isValid bool) {
	// Handle arrow keys
	switch input.Key {
	case keyboard.KeyArrowUp:
		return game.Point{X: 0, Y: -1}, true
	case keyboard.KeyArrowDown:
		return game.Point{X: 0, Y: 1}, true
	case keyboard.KeyArrowLeft:
		return game.Point{X: -1, Y: 0}, true
	case keyboard.KeyArrowRight:
		return game.Point{X: 1, Y: 0}, true
	}

	// Handle WASD keys
	switch input.Char {
	case 'w', 'W':
		return game.Point{X: 0, Y: -1}, true
	case 's', 'S':
		return game.Point{X: 0, Y: 1}, true
	case 'a', 'A':
		return game.Point{X: -1, Y: 0}, true
	case 'd', 'D':
		return game.Point{X: 1, Y: 0}, true
	}

	return game.Point{}, false
}

// IsQuit checks if the input is a quit command
func IsQuit(input KeyInput) bool {
	return input.Char == 'q' || input.Char == 'Q'
}

// IsRestart checks if the input is a restart command
func IsRestart(input KeyInput) bool {
	return input.Char == 'r' || input.Char == 'R'
}

// IsPause checks if the input is a pause command
func IsPause(input KeyInput) bool {
	return input.Char == 'p' || input.Char == 'P' || input.Char == ' '
}
