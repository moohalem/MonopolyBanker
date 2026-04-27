package main

import (
	"fmt"

	"monopoly/internal/game"
	"monopoly/internal/ui"
)

func main() {
	ui.ClearScreen()
	fmt.Println()
	playerCount := game.Welcome()
	players := game.NewPlayers(playerCount)
	game.Play(players)
}
