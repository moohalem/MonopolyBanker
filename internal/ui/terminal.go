// Package ui provides terminal input and output helpers.
package ui

import (
	"bufio"
	"fmt"
	"os"
)

// scanner is reused across all input calls to avoid repeated allocations.
var scanner = bufio.NewScanner(os.Stdin)

// UserInput displays a styled prompt and returns the user's input.
func UserInput(prompt string) string {
	fmt.Printf("  %s▸%s %s ", Cyan, Reset, prompt)
	scanner.Scan()
	return scanner.Text()
}

// ClearScreen resets the terminal viewport using ANSI escape codes.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
