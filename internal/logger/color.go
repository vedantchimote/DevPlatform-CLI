package logger

import "fmt"

// Color represents an ANSI color code
type Color string

const (
	ColorReset  Color = "\033[0m"
	ColorRed    Color = "\033[31m"
	ColorGreen  Color = "\033[32m"
	ColorYellow Color = "\033[33m"
	ColorBlue   Color = "\033[34m"
	ColorCyan   Color = "\033[36m"
)

// Colorize wraps text with ANSI color codes
func Colorize(text string, color Color) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorReset)
}

// Red returns red colored text
func Red(text string) string {
	return Colorize(text, ColorRed)
}

// Green returns green colored text
func Green(text string) string {
	return Colorize(text, ColorGreen)
}

// Yellow returns yellow colored text
func Yellow(text string) string {
	return Colorize(text, ColorYellow)
}

// Blue returns blue colored text
func Blue(text string) string {
	return Colorize(text, ColorBlue)
}

// Cyan returns cyan colored text
func Cyan(text string) string {
	return Colorize(text, ColorCyan)
}
