package game

import (
	"fmt"
	"strconv"
	"time"

	"monopoly/internal/ui"
)

const (
	startingBalance = 15_000_000
	minPlayers      = 1
	maxPlayers      = 6
)

// Welcome displays the title screen and asks for the number of players.
// It loops until a valid count (1–6) is provided.
func Welcome() int {
	for {
		ui.ClearScreen()
		ui.Banner(true)
		ui.Separator()

		playerStr := ui.UserInput("Enter the number of players (1 - 6):")

		playerCount, err := strconv.Atoi(playerStr)
		if err != nil {
			ui.ClearScreen()
			ui.Error("Input must be a number between 1 - 6. Please try again.")
			time.Sleep(2 * time.Second)
			continue
		}

		if playerCount >= minPlayers && playerCount <= maxPlayers {
			return playerCount
		}

		ui.ClearScreen()
		ui.Error("Player count must be between 1 - 6. Please try again.")
		time.Sleep(2 * time.Second)
	}
}

// NewPlayers prompts for each player's name and returns the initialised roster.
func NewPlayers(count int) []Player {
	ui.ClearScreen()
	players := make([]Player, 0, count)

	for i := range count {
		ui.ClearScreen()
		ui.Banner(false)
		ui.Separator()
		players = append(players, Player{
			Name:    ui.UserInput(fmt.Sprintf("Enter name for player %d:", i+1)),
			Balance: startingBalance,
		})
		ui.ClearScreen()
	}
	return players
}
