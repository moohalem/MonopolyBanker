package game

import (
	"fmt"
	"strconv"
	"time"

	"monopoly/internal/ui"
)

func NewPlayers(count int) []Player {
	ui.ClearScreen()
	players := make([]Player, 0, count)

	for i := range count {
		ui.ClearScreen()
		fmt.Println("---------------------------")
		fmt.Println("|   Monopoly Banker App   |")
		fmt.Println("---------------------------")
		fmt.Println("")
		fmt.Println("==============================")
		players = append(players, Player{
			Name:    ui.UserInput(fmt.Sprintf("Enter name for player %d:", i+1)),
			Balance: 15000000,
		})
		ui.ClearScreen()
	}
	return players
}

func Welcome() int {
	for {
		ui.ClearScreen()
		fmt.Println("---------------------------")
		fmt.Println("|   Monopoly Banker App   |")
		fmt.Println("|        v0.1.0           |")
		fmt.Println("---------------------------")
		fmt.Println("")
		fmt.Println("==============================")
		playerStr := ui.UserInput("Enter the number of players (1 - 6): ")

		playerCount, err := strconv.Atoi(playerStr)
		if err != nil {
			ui.ClearScreen()
			fmt.Println("Input must be a number between 1 - 6. Please try again.")
			time.Sleep(2 * time.Second)
			continue
		}

		if playerCount >= 1 && playerCount <= 6 {
			return playerCount
		}

		ui.ClearScreen()
		fmt.Println("Player count must be between 1 - 6. Please try again.")
		time.Sleep(2 * time.Second)
	}
}
