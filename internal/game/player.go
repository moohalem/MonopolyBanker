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

const goBonus = 2_000_000

// ---------------------------------------------------------------------------
// Main game loop
// ---------------------------------------------------------------------------

// Play runs the main game loop: display balances, show the menu, and
// dispatch the chosen action.
func Play(players []Player) {
	for {
		ui.ClearScreen()
		ui.Banner(false)
		printBalances(players)

		switch menu() {
		case 1:
			buy(players)
		case 2:
			payRent(players)
		case 3:
			receive(players)
		case 4:
			payMortgage(players)
		case 5:
			passingGo(players)
		case 0:
			return
		}
	}
}

// printBalances renders the player table with box-drawing borders.
func printBalances(players []Player) {
	// Compute column widths dynamically.
	nameWidth := 10
	for _, p := range players {
		if len(p.Name) > nameWidth {
			nameWidth = len(p.Name)
		}
	}
	nameWidth += 2 // padding

	moneyWidth := 16

	// Header
	fmt.Printf("  %s┌─────┬%s┬%s┐%s\n",
		ui.Dim, strings.Repeat("─", nameWidth), strings.Repeat("─", moneyWidth), ui.Reset)
	fmt.Printf("  %s│%s  #  %s│%s %-*s%s│%s %*s%s│%s\n",
		ui.Dim, ui.Reset+ui.Bold+ui.White, ui.Reset+ui.Dim,
		ui.Reset+ui.Bold+ui.White, nameWidth-1, "Player", ui.Reset+ui.Dim,
		ui.Reset+ui.Bold+ui.White, moneyWidth-1, "Balance", ui.Reset+ui.Dim, ui.Reset)
	fmt.Printf("  %s├─────┼%s┼%s┤%s\n",
		ui.Dim, strings.Repeat("─", nameWidth), strings.Repeat("─", moneyWidth), ui.Reset)

	// Rows
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

	// Footer
	fmt.Printf("  %s└─────┴%s┴%s┘%s\n",
		ui.Dim, strings.Repeat("─", nameWidth), strings.Repeat("─", moneyWidth), ui.Reset)
}

// menu displays the action menu and returns the user's choice.
func menu() int {
	fmt.Println()
	menuItem := func(key, label string) {
		fmt.Printf("  %s[%s%s%s%s]%s  %s\n",
			ui.Dim, ui.Reset+ui.Bold+ui.Cyan, key, ui.Reset+ui.Dim, "", ui.Reset, label)
	}
	menuItem("1", "Buy / Pay to The Bank")
	menuItem("2", "Pay Rent")
	menuItem("3", "Receive Money")
	menuItem("4", "Paying Mortgage")
	menuItem("5", "Passing GO")
	menuItem("0", "Exit")
	fmt.Println()

	choice, err := strconv.Atoi(ui.UserInput("Enter menu:"))
	if err != nil {
		return 0
	}
	return choice
}

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------

func buy(players []Player) {
	ui.ClearScreen()
	ui.Banner(false)

	idx, ok := selectPlayer(players, "Which player is buying / paying?")
	if !ok {
		showCancelled()
		return
	}

	amount := getAmount("How much?")
	p := &players[idx]
	if p.Balance >= amount {
		p.Balance -= amount
		ui.Success(fmt.Sprintf("%s bought / paid for %s. New balance: %s",
			p.Name, ui.FormatMoney(amount), ui.FormatMoney(p.Balance)))
	} else {
		ui.Error("Insufficient balance.")
	}
	ui.UserInput("Press enter to continue...")
}

func payRent(players []Player) {
	ui.ClearScreen()
	ui.Banner(false)

	fromIdx, ok := selectPlayer(players, "Which player is paying rent?")
	if !ok {
		showCancelled()
		return
	}

	toIdx, ok := selectPlayer(players, "Which player is receiving rent?")
	if !ok {
		showCancelled()
		return
	}

	amount := getAmount("How much rent?")
	from, to := &players[fromIdx], &players[toIdx]
	if from.Balance >= amount {
		from.Balance -= amount
		to.Balance += amount
		ui.Success(fmt.Sprintf("%s paid %s to %s.",
			from.Name, ui.FormatMoney(amount), to.Name))
	} else {
		ui.Error("Insufficient balance.")
	}
	ui.UserInput("Press enter to continue...")
}

func receive(players []Player) {
	ui.ClearScreen()
	ui.Banner(false)

	idx, ok := selectPlayer(players, "Which player is receiving?")
	if !ok {
		showCancelled()
		return
	}

	amount := getAmount("How much?")
	p := &players[idx]
	p.Balance += amount
	ui.Success(fmt.Sprintf("%s received %s. New balance: %s",
		p.Name, ui.FormatMoney(amount), ui.FormatMoney(p.Balance)))
	ui.UserInput("Press enter to continue...")
}

func payMortgage(players []Player) {
	ui.ClearScreen()
	ui.Banner(false)

	idx, ok := selectPlayer(players, "Which player is paying mortgage?")
	if !ok {
		showCancelled()
		return
	}

	amount := getAmount("How much is the mortgage payment?")
	total := roundUpWithTax(amount)
	ui.Info(fmt.Sprintf("Mortgage + 10%% tax = %s", ui.FormatMoney(total)))

	confirm := ui.UserInput("Do you agree? (Y/n)")
	if strings.ToLower(confirm) != "y" && confirm != "" {
		return
	}

	p := &players[idx]
	if p.Balance >= total {
		p.Balance -= total
		ui.Success(fmt.Sprintf("%s paid mortgage of %s. New balance: %s",
			p.Name, ui.FormatMoney(total), ui.FormatMoney(p.Balance)))
	} else {
		ui.Error("Insufficient balance.")
	}
	ui.UserInput("Press enter to continue...")
}

func passingGo(players []Player) {
	ui.ClearScreen()
	ui.Banner(false)

	idx, ok := selectPlayer(players, "Which player is passing GO?")
	if !ok {
		showCancelled()
		return
	}

	p := &players[idx]
	p.Balance += goBonus
	ui.Success(fmt.Sprintf("%s passed GO and received %s. New balance: %s",
		p.Name, ui.FormatMoney(goBonus), ui.FormatMoney(p.Balance)))
	ui.UserInput("Press enter to continue...")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// showCancelled displays a brief cancellation notice.
func showCancelled() {
	ui.ClearScreen()
	ui.Warning("The last operation is cancelled")
	time.Sleep(1 * time.Second)
}

// selectPlayer lists players and returns the chosen index (0-based).
// The second return value is false when the user cancels.
func selectPlayer(players []Player, prompt string) (int, bool) {
	for {
		fmt.Println()
		ui.Info(prompt)
		fmt.Println()
		for i, p := range players {
			fmt.Printf("    %s%s%d%s  %s%s%s\n",
				ui.Bold, ui.Yellow, i+1, ui.Reset,
				ui.White, p.Name, ui.Reset)
		}
		fmt.Println()

		input := ui.UserInput("Enter player number (or 'c' to cancel):")
		if strings.EqualFold(input, "c") {
			return 0, false
		}

		idx, err := strconv.Atoi(input)
		if err == nil && idx >= 1 && idx <= len(players) {
			return idx - 1, true
		}

		ui.ClearScreen()
		ui.Banner(false)
		ui.Error("Invalid input.")
		time.Sleep(2 * time.Second)
	}
}

// getAmount asks for an amount and a unit multiplier, returning the final
// value. It validates the number before prompting for the unit.
func getAmount(prompt string) int {
	for {
		amountStr := ui.UserInput(prompt)
		if strings.EqualFold(amountStr, "c") {
			return 0
		}

		amount, err := strconv.Atoi(amountStr)
		if err != nil || amount <= 0 {
			ui.Error("Invalid amount.")
			continue
		}

		fmt.Println()
		fmt.Printf("    %s[%s1%s]%s  M (×1,000,000)\n", ui.Dim, ui.Reset+ui.Bold+ui.Cyan, ui.Reset+ui.Dim, ui.Reset)
		fmt.Printf("    %s[%s2%s]%s  K (×1,000)\n", ui.Dim, ui.Reset+ui.Bold+ui.Cyan, ui.Reset+ui.Dim, ui.Reset)
		fmt.Printf("    %s[%s3%s]%s  As is\n", ui.Dim, ui.Reset+ui.Bold+ui.Cyan, ui.Reset+ui.Dim, ui.Reset)
		fmt.Println()

		unit := ui.UserInput("Enter unit choice:")
		switch unit {
		case "1":
			return amount * 1_000_000
		case "2":
			return amount * 1_000
		case "3":
			return amount
		default:
			ui.Error("Invalid unit choice.")
		}
	}
}

// roundUpWithTax adds 10% tax and rounds up to the nearest 10,000.
func roundUpWithTax(amount int) int {
	withTax := (amount * 11) / 10
	return ((withTax + 9_999) / 10_000) * 10_000
}
