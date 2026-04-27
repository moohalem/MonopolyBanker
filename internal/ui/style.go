package ui

import (
	"fmt"
	"strings"
)

// ANSI escape codes for terminal styling.
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Red     = "\033[91m"
	Green   = "\033[92m"
	Yellow  = "\033[93m"
	Blue    = "\033[94m"
	Magenta = "\033[95m"
	Cyan    = "\033[96m"
	White   = "\033[97m"
	Gray    = "\033[90m"
)

// Banner prints the application title in a decorative box.
func Banner(showVersion bool) {
	const width = 44
	title := "M O N O P O L Y   B A N K E R"
	border := Bold + Cyan

	fmt.Println()
	fmt.Printf("  %s╔%s╗%s\n", border, strings.Repeat("═", width), Reset)
	printCentered(width, title, Bold+Yellow)
	if showVersion {
		printCentered(width, "v0.1.0", Dim+White)
	}
	fmt.Printf("  %s╚%s╝%s\n", border, strings.Repeat("═", width), Reset)
	fmt.Println()
}

// printCentered prints text centered inside the box border.
func printCentered(width int, text, style string) {
	border := Bold + Cyan
	pad := (width - len(text)) / 2
	right := width - pad - len(text)
	fmt.Printf("  %s║%s%s%s%s%s║%s\n",
		border,
		strings.Repeat(" ", pad),
		style, text, Reset,
		strings.Repeat(" ", right),
		Reset)
}

// FormatMoney formats an integer with comma separators.
func FormatMoney(n int) string {
	negative := n < 0
	if negative {
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	var b strings.Builder
	b.Grow(len(s) + len(s)/3)
	for i, ch := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	if negative {
		return "-" + b.String()
	}
	return b.String()
}

// Success prints a green success message.
func Success(msg string) {
	fmt.Printf("  %s%s✓ %s%s\n", Bold, Green, msg, Reset)
}

// Error prints a red error message.
func Error(msg string) {
	fmt.Printf("  %s%s✗ %s%s\n", Bold, Red, msg, Reset)
}

// Info prints a cyan informational message.
func Info(msg string) {
	fmt.Printf("  %s%sℹ %s%s\n", Bold, Cyan, msg, Reset)
}

// Warning prints a yellow warning message.
func Warning(msg string) {
	fmt.Printf("  %s%s⚠ %s%s\n", Bold, Yellow, msg, Reset)
}

// Separator prints a decorative divider line.
func Separator() {
	fmt.Printf("  %s%s%s\n", Dim, strings.Repeat("─", 46), Reset)
}
