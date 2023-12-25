package solver

import "github.com/levilutz/minesweeper/pkg/board"

// Run the given function for each revealed tile.
// When the passed-in function returns true, the top level function quits with true.
func forEachRevealed(b *board.Board, fn func(x, y int) bool) bool {
	size := b.GetSize()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if b.Revealed(x, y) {
				out := fn(x, y)
				if out {
					return true
				}
			}
		}
	}
	return false
}
