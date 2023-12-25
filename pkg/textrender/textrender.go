package textrender

import (
	"fmt"
	"strconv"

	"github.com/levilutz/minesweeper/pkg/board"
	"github.com/levilutz/minesweeper/pkg/util/bashcolor"
)

var colorMap = map[int]string{
	1: bashcolor.BrightBlue,
	2: bashcolor.BrightGreen,
	3: bashcolor.Red,
	4: bashcolor.Blue,
	5: bashcolor.Red,
	6: bashcolor.BrightCyan,
	7: bashcolor.Gray,
	8: bashcolor.Gray,
}

// Colorize a number.
func colorNum(num int) string {
	if num == 0 {
		return " "
	}
	color, ok := colorMap[num]
	if !ok {
		return "?"
	} else {
		return bashcolor.Color(strconv.Itoa(num), color)
	}
}

func renderTile(b *board.Board, x, y int) string {
	hasMine, hasFlag, revealed, neighbors := b.GetTile(x, y)
	if revealed {
		if hasMine {
			return "X"
		} else {
			return colorNum(neighbors)
		}
	} else {
		if hasFlag {
			return bashcolor.Color("#", bashcolor.BrightRed)
		} else {
			return "+"
		}
	}
}

func RenderBoard(b *board.Board) string {
	out := "\n"
	for y := b.GetSize() - 1; y >= 0; y-- {
		out += fmt.Sprintf("%x | ", y)
		for x := 0; x < b.GetSize(); x++ {
			out += " " + renderTile(b, x, y)
		}
		out += "\n"
	}
	out += "--+-"
	for x := 0; x < b.GetSize(); x++ {
		out += "--"
	}
	out += "\n  | "
	for x := 0; x < b.GetSize(); x++ {
		out += fmt.Sprintf(" %x", x)
	}
	out += "\n"

	return out
}
