// Package game provides game logic for the Monopoly banker application.
package game

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"monopoly/internal/ui"
)

// Player holds a participant's name and current balance.
type Player struct {
	Name    string
	Balance int
}

type balanceChange struct {
	PlayerIdx int
	Delta     int // positive = credited, negative = debited
}

type transaction struct {
	Description string
	Changes     []balanceChange
}

const goBonus = 2_000_000

// Play runs the main game loop.
func Play(players []Player) {
	var history []transaction

	for {
		ui.ClearScreen()
		ui.Banner(false)
		printBalances(players)
		printHistory(history)

		switch strings.ToLower(showMenu()) {
		case "1":
			pay(players, &history)
		case "2":
			transfer(players, &history)
		case "3":
			receive(players, &history)
		case "g":
			passGo(players, &history)
		case "u":
			undo(players, &history)
		case "0":
			return
		}
	}
}

func printBalances(players []Player) {
	nameWidth := 10
	for _, p := range players {
		if len(p.Name) > nameWidth {
			nameWidth = len(p.Name)
		}
	}
	nameWidth += 2
	moneyWidth := 16

	fmt.Printf("  %s┌─────┬%s┬%s┐%s\n",
		ui.Dim, strings.Repeat("─", nameWidth), strings.Repeat("─", moneyWidth), ui.Reset)
	fmt.Printf("  %s│%s  #  %s│%s %-*s%s│%s %*s%s│%s\n",
		ui.Dim, ui.Reset+ui.Bold+ui.White, ui.Reset+ui.Dim,
		ui.Reset+ui.Bold+ui.White, nameWidth-1, "Player", ui.Reset+ui.Dim,
		ui.Reset+ui.Bold+ui.White, moneyWidth-1, "Balance", ui.Reset+ui.Dim, ui.Reset)
	fmt.Printf("  %s├─────┼%s┼%s┤%s\n",
		ui.Dim, strings.Repeat("─", nameWidth), strings.Repeat("─", moneyWidth), ui.Reset)

	for i, p := range players {
		money := ui.FormatMoney(p.Balance)
		color := ui.Green
		if p.Balance <= 0 {
			color = ui.Red
		}
		fmt.Printf("  %s│%s  %s%d%s  %s│%s %s%-*s%s%s│%s %s%*s%s %s│%s\n",
			ui.Dim, ui.Reset,
			ui.Yellow, i+1, ui.Reset,
			ui.Dim, ui.Reset,
			ui.Bold+ui.White, nameWidth-1, p.Name, ui.Reset,
			ui.Dim, ui.Reset,
			color, moneyWidth-2, money, ui.Reset,
			ui.Dim, ui.Reset)
	}

	fmt.Printf("  %s└─────┴%s┴%s┘%s\n",
		ui.Dim, strings.Repeat("─", nameWidth), strings.Repeat("─", moneyWidth), ui.Reset)
}

func printHistory(history []transaction) {
	if len(history) == 0 {
		return
	}
	fmt.Println()
	fmt.Printf("  %s%sRecent%s\n", ui.Dim, ui.White, ui.Reset)
	start := 0
	if len(history) > 5 {
		start = len(history) - 5
	}
	for i := len(history) - 1; i >= start; i-- {
		fmt.Printf("  %s  %s%s\n", ui.Gray, history[i].Description, ui.Reset)
	}
}

func showMenu() string {
	fmt.Println()
	mi := func(key, label string) {
		fmt.Printf("  %s[%s%s%s]%s  %s\n",
			ui.Dim, ui.Reset+ui.Bold+ui.Cyan, key, ui.Reset+ui.Dim, ui.Reset, label)
	}
	mi("1", "Pay to Bank")
	mi("2", "Transfer to Player")
	mi("3", "Receive from Bank")
	mi("G", "Pass GO (+2,000,000)")
	mi("U", "Undo Last")
	mi("0", "Exit")
	fmt.Println()
	return ui.UserInput("Enter choice:")
}

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------

func pay(players []Player, history *[]transaction) {
	ui.ClearScreen()
	ui.Banner(false)

	idx, ok := selectPlayer(players, "Which player is paying?")
	if !ok {
		showCancelled()
		return
	}

	isMortgage := strings.EqualFold(ui.UserInput("Is this a mortgage? (y/N)"), "y")

	amount := getAmount("Amount")
	if amount == 0 {
		showCancelled()
		return
	}

	final := amount
	if isMortgage {
		final = roundUpWithTax(amount)
		ui.Info(fmt.Sprintf("Mortgage + 10%% tax = %s", ui.FormatMoney(final)))
		confirm := ui.UserInput("Proceed? (Y/n)")
		if strings.ToLower(confirm) != "y" && confirm != "" {
			showCancelled()
			return
		}
	}

	p := &players[idx]
	if p.Balance < final {
		ui.Error("Insufficient balance.")
		ui.UserInput("Press enter to continue...")
		return
	}

	p.Balance -= final
	kind := "paid"
	if isMortgage {
		kind = "paid mortgage"
	}
	desc := fmt.Sprintf("%s %s %s → bank", p.Name, kind, ui.FormatMoney(final))
	*history = append(*history, transaction{
		Description: desc,
		Changes:     []balanceChange{{PlayerIdx: idx, Delta: -final}},
	})
	ui.Success(desc)
	ui.UserInput("Press enter to continue...")
}

func transfer(players []Player, history *[]transaction) {
	ui.ClearScreen()
	ui.Banner(false)

	fromIdx, ok := selectPlayer(players, "Which player is paying?")
	if !ok {
		showCancelled()
		return
	}

	toIdx, ok := selectPlayer(players, "Which player is receiving?")
	if !ok {
		showCancelled()
		return
	}

	if fromIdx == toIdx {
		ui.Error("Cannot transfer to the same player.")
		ui.UserInput("Press enter to continue...")
		return
	}

	amount := getAmount("Amount")
	if amount == 0 {
		showCancelled()
		return
	}

	from, to := &players[fromIdx], &players[toIdx]
	if from.Balance < amount {
		ui.Error("Insufficient balance.")
		ui.UserInput("Press enter to continue...")
		return
	}

	from.Balance -= amount
	to.Balance += amount
	desc := fmt.Sprintf("%s paid %s to %s", from.Name, ui.FormatMoney(amount), to.Name)
	*history = append(*history, transaction{
		Description: desc,
		Changes: []balanceChange{
			{PlayerIdx: fromIdx, Delta: -amount},
			{PlayerIdx: toIdx, Delta: amount},
		},
	})
	ui.Success(desc)
	ui.UserInput("Press enter to continue...")
}

func receive(players []Player, history *[]transaction) {
	ui.ClearScreen()
	ui.Banner(false)

	idx, ok := selectPlayer(players, "Which player is receiving?")
	if !ok {
		showCancelled()
		return
	}

	amount := getAmount("Amount")
	if amount == 0 {
		showCancelled()
		return
	}

	p := &players[idx]
	p.Balance += amount
	desc := fmt.Sprintf("%s received %s from bank", p.Name, ui.FormatMoney(amount))
	*history = append(*history, transaction{
		Description: desc,
		Changes:     []balanceChange{{PlayerIdx: idx, Delta: amount}},
	})
	ui.Success(desc)
	ui.UserInput("Press enter to continue...")
}

func passGo(players []Player, history *[]transaction) {
	ui.ClearScreen()
	ui.Banner(false)

	idx, ok := selectPlayer(players, "Which player is passing GO?")
	if !ok {
		showCancelled()
		return
	}

	p := &players[idx]
	p.Balance += goBonus
	desc := fmt.Sprintf("%s passed GO (+%s)", p.Name, ui.FormatMoney(goBonus))
	*history = append(*history, transaction{
		Description: desc,
		Changes:     []balanceChange{{PlayerIdx: idx, Delta: goBonus}},
	})
	ui.Success(desc)
	ui.UserInput("Press enter to continue...")
}

func undo(players []Player, history *[]transaction) {
	if len(*history) == 0 {
		ui.Warning("Nothing to undo.")
		time.Sleep(1 * time.Second)
		return
	}

	last := (*history)[len(*history)-1]
	for _, c := range last.Changes {
		players[c.PlayerIdx].Balance -= c.Delta
	}
	*history = (*history)[:len(*history)-1]
	ui.Success(fmt.Sprintf("Undone: %s", last.Description))
	time.Sleep(1 * time.Second)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func showCancelled() {
	ui.ClearScreen()
	ui.Warning("Operation cancelled")
	time.Sleep(1 * time.Second)
}

func selectPlayer(players []Player, prompt string) (int, bool) {
	for {
		fmt.Println()
		ui.Info(prompt)
		fmt.Println()
		for i, p := range players {
			fmt.Printf("    %s%s%d%s  %s%s%s  %s%s%s\n",
				ui.Bold, ui.Yellow, i+1, ui.Reset,
				ui.White, p.Name, ui.Reset,
				ui.Dim, ui.FormatMoney(p.Balance), ui.Reset)
		}
		fmt.Println()

		input := ui.UserInput("Player number (or 'c' to cancel):")
		if strings.EqualFold(input, "c") {
			return 0, false
		}

		idx, err := strconv.Atoi(input)
		if err == nil && idx >= 1 && idx <= len(players) {
			return idx - 1, true
		}
		ui.Error("Invalid input.")
		time.Sleep(1 * time.Second)
	}
}

// getAmount prompts for an amount with inline unit support.
// Accepts: "15M" (millions), "500K" (thousands), or plain numbers.
func getAmount(prompt string) int {
	for {
		input := ui.UserInput(fmt.Sprintf("%s (e.g. 15M, 500K, or number):", prompt))
		if strings.EqualFold(input, "c") {
			return 0
		}

		amount, err := parseAmount(input)
		if err != nil || amount <= 0 {
			ui.Error("Invalid amount. Examples: 15M, 500K, 250000")
			continue
		}
		ui.Info(fmt.Sprintf("= %s", ui.FormatMoney(amount)))
		return amount
	}
}

// parseAmount interprets "15M", "500K", or plain integers.
func parseAmount(s string) (int, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	multiplier := 1
	switch strings.ToUpper(s[len(s)-1:]) {
	case "M":
		multiplier = 1_000_000
		s = s[:len(s)-1]
	case "K":
		multiplier = 1_000
		s = s[:len(s)-1]
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return n * multiplier, nil
}

// roundUpWithTax adds 10% tax and rounds up to the nearest 10,000.
func roundUpWithTax(amount int) int {
	withTax := (amount * 11) / 10
	return ((withTax + 9_999) / 10_000) * 10_000
}
