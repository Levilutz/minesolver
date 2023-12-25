package solver

import (
	"github.com/levilutz/minesweeper/pkg/board"
)

// Compute until a single command is run (either a flag or a reveal).
// Returns whether something could be run, false if stuck.
func Pass(b *board.Board) bool {
	if revealIfFresh(b) {
		return true
	}

	if findObviousMines(b) {
		return true
	}

	if findObiousEmpty(b) {
		return true
	}

	// At this point, solves 13% of 8x8 w/ 10 mines

	if deduce(b) {
		return true
	}

	return false
}
