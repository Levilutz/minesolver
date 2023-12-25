package bashcolor

import (
	"fmt"
)

const (
	Red          = "31"
	Green        = "32"
	Yellow       = "33"
	Blue         = "34"
	Purple       = "35"
	Cyan         = "36"
	Gray         = "90"
	BrightRed    = "91"
	BrightGreen  = "92"
	BrightYellow = "93"
	BrightBlue   = "94"
	BrightPurple = "95"
	BrightCyan   = "96"
	White        = "97"
)

// Colorize text.
func Color(text string, color string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", color, text)
}
