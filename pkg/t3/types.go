// Package t3 contains shared public types for go-t3.
package t3

// Player identifies a player in a game.
type Player string

const (
	PlayerX  Player = "X"
	PlayerO  Player = "O"
	NoPlayer Player = ""
)
