// Package ui provides terminal input and output helpers.
package ui

import (
	"bufio"
	"fmt"
	"os"
)

func UserInput(text string) string {
	fmt.Print(text + " ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
