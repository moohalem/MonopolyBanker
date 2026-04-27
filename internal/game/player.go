// Package game provides game logic for the Monopoly banker application.
package game

import (
	"fmt"
	"strconv"
	"strings"

	"monopoly/internal/ui"
)

type Player struct {
	Name    string
	Balance int
}

func Play(players []Player) {
	for {
		ui.ClearScreen()
		fmt.Println("---------------------------")
		fmt.Println("|   Monopoly Banker App   |")
		fmt.Println("---------------------------")
		fmt.Println("")
		fmt.Println("No. Name: Balance")
		fmt.Println("++++++++++++++++++++++++++++++")
		for i, player := range players {
			fmt.Printf("%d. %s: %d\n", i+1, player.Name, player.Balance)
		}
		fmt.Println("++++++++++++++++++++++++++++++")
		choice := menu()
		switch choice {
		case 1:
			buy(players)
		case 2:
			payRent(players)
		case 3:
			receive(players)
		case 4:
			payMortgage(players)
		case 5:
			PassingGo(players)
		case 0:
			return // exit
		}
	}
}

func menu() int {
	fmt.Println("")
	fmt.Println("1. Buy / Pay to The Bank")
	fmt.Println("2. Pay Rent")
	fmt.Println("3. Receive")
	fmt.Println("4. Paying Mortgage")
	fmt.Println("5. Passing GO")
	fmt.Println("0. Exit")
	choice, err := strconv.Atoi(ui.UserInput("Enter menu: "))
	if err != nil {
		return 0
	}
	return choice
}

func buy(players []Player) {
	ui.ClearScreen()
	playerIndex := getPlayerIndex(players, "Which player is buying / paying?")
	amount := getAmount("How much?")
	if players[playerIndex].Balance >= amount {
		players[playerIndex].Balance -= amount
		fmt.Printf("%s bought / paid for %d. New balance: %d\n", players[playerIndex].Name, amount, players[playerIndex].Balance)
	} else {
		fmt.Println("Insufficient balance.")
	}
	ui.UserInput("Press enter to continue...")
}

func payRent(players []Player) {
	ui.ClearScreen()
	fromIndex := getPlayerIndex(players, "Which player is paying rent?")
	toIndex := getPlayerIndex(players, "Which player is receiving rent?")
	amount := getAmount("How much rent?")
	if players[fromIndex].Balance >= amount {
		players[fromIndex].Balance -= amount
		players[toIndex].Balance += amount
		fmt.Printf("%s paid %d to %s.\n", players[fromIndex].Name, amount, players[toIndex].Name)
	} else {
		fmt.Println("Insufficient balance.")
	}
	ui.UserInput("Press enter to continue...")
}

func receive(players []Player) {
	ui.ClearScreen()
	playerIndex := getPlayerIndex(players, "Which player is receiving?")
	amount := getAmount("How much?")
	players[playerIndex].Balance += amount
	fmt.Printf("%s received %d. New balance: %d\n", players[playerIndex].Name, amount, players[playerIndex].Balance)
	ui.UserInput("Press enter to continue...")
}

func payMortgage(players []Player) {
	ui.ClearScreen()
	playerIndex := getPlayerIndex(players, "Which player is paying mortgage?")
	amount := getAmount("How much is the mortgage payment?")
	amountPlusTax := nearest10k(amount)
	fmt.Printf("The mortgage payment + 10 percent tax is: %d\n", amountPlusTax)
	confirm := ui.UserInput("Do you agree? (Y/n)")
	if strings.ToLower(confirm) == "y" || confirm == "" {
		if players[playerIndex].Balance >= amountPlusTax {
			players[playerIndex].Balance -= amountPlusTax
			fmt.Printf("%s paid mortgage of %d. New balance: %d\n", players[playerIndex].Name, amountPlusTax, players[playerIndex].Balance)
		} else {
			fmt.Println("Insufficient balance.")
		}
		ui.UserInput("Press enter to continue...")
	}
}

func PassingGo(players []Player) {
	ui.ClearScreen()
	playerIndex := getPlayerIndex(players, "Which player is passing GO?")
	players[playerIndex].Balance += 2_000_000
	fmt.Printf("%s passed GO and received 2,000,000. New balance: %d\n", players[playerIndex].Name, players[playerIndex].Balance)
	ui.UserInput("Press enter to continue...")
}

func nearest10k(amount int) int {
	// Add 10% tax using pure integer math
	amountPlusTax := (amount * 11) / 10

	// Always round UP to the next 10,000
	rounded := ((amountPlusTax + 9999) / 10_000) * 10_000

	return rounded
}

func getPlayerIndex(players []Player, prompt string) int {
	for {
		fmt.Println(prompt)
		for i, p := range players {
			fmt.Printf("%d. %s\n", i+1, p.Name)
		}
		indexStr := ui.UserInput("Enter player number: ")
		index, err := strconv.Atoi(indexStr)
		if err == nil && index >= 1 && index <= len(players) {
			return index - 1
		}
		fmt.Println("Invalid input.")
	}
}

func getAmount(prompt string) int {
	for {
		amountStr := ui.UserInput(prompt)
		amount, err := strconv.Atoi(amountStr)
		if err == nil && amount > 0 {
			return amount
		}
		fmt.Println("Invalid amount.")
	}
}
