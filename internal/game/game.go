// Package game contains core tic-tac-toe board logic.
package game

// Cell represents a single board cell value.
type Cell int

const (
	Empty   Cell = iota
	PlayerX      // X player
	PlayerO      // O player
)

// Board is a 3x3 tic-tac-toe board.
type Board [3][3]Cell

// GameState holds the current state of a game.
type GameState struct {
	Board   Board
	Current Cell // whose turn it is
	Winner  Cell // Empty if no winner yet
	Done    bool
}
