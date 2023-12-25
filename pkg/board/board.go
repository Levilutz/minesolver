package board

import (
	"fmt"

	"github.com/levilutz/minesweeper/pkg/util"
)

// The game board
type Board struct {
	size          int
	mines         [][]bool
	flags         [][]bool
	revealed      [][]bool
	neighbors     [][]int
	neighborCache [][][]util.Vec
}

// Create a new game board.
func NewBoard(size int) *Board {
	return &Board{
		size:          size,
		mines:         util.DArray[bool](size),
		flags:         util.DArray[bool](size),
		revealed:      util.DArray[bool](size),
		neighbors:     util.DArray[int](size),
		neighborCache: util.DArray[[]util.Vec](size),
	}
}

// Reset the game board.
func (b *Board) Reset() {
	for x := 0; x < b.size; x++ {
		for y := 0; y < b.size; y++ {
			b.mines[x][y] = false
			b.flags[x][y] = false
			b.revealed[x][y] = false
			b.neighbors[x][y] = 0
		}
	}
}

// Get the size of the board.
func (b *Board) GetSize() int {
	return b.size
}

// Check whether the game has any revealed tiles.
func (b *Board) HasReveals() bool {
	for x := 0; x < b.size; x++ {
		for y := 0; y < b.size; y++ {
			if b.revealed[x][y] {
				return true
			}
		}
	}
	return false
}

// Check whether the game is complete (all non-mines revealed).
func (b *Board) Complete() bool {
	for x := 0; x < b.size; x++ {
		for y := 0; y < b.size; y++ {
			if !b.mines[x][y] && !b.revealed[x][y] {
				return false
			}
		}
	}
	return true
}

// Check whether any mines have been revealed (game loss).
func (b *Board) HasRevealedMines() bool {
	for x := 0; x < b.size; x++ {
		for y := 0; y < b.size; y++ {
			if b.mines[x][y] && b.revealed[x][y] {
				return true
			}
		}
	}
	return false
}

// Count the number of remaining unflagged mines.
func (b *Board) UnflaggedMines() int {
	out := 0
	for x := 0; x < b.size; x++ {
		for y := 0; y < b.size; y++ {
			if b.mines[x][y] && !b.flags[x][y] {
				out += 1
			}
		}
	}
	return out
}

// Get data for the given tile.
func (b *Board) GetTile(x, y int) (hasMine, hasFlag, revealed bool, neighbors int) {
	return b.mines[x][y], b.flags[x][y], b.revealed[x][y], b.neighbors[x][y]
}

// Check whether the given tile has a flag.
func (b *Board) HasFlag(x, y int) bool {
	return b.flags[x][y]
}

// Check whether the given tile has a mine.
func (b *Board) HasMine(x, y int) bool {
	return b.mines[x][y]
}

// Get the number of neighbors the given tile has.
func (b *Board) GetNumNeighbors(x, y int) int {
	return b.neighbors[x][y]
}

// Check whether the given tile is revealed.
func (b *Board) Revealed(x, y int) bool {
	return b.revealed[x][y]
}

// Place a mine.
func (b *Board) PlaceMine(x, y int) {
	b.mines[x][y] = true
	for _, neighbor := range util.GetNeighbors(x, y, b.size) {
		b.neighbors[neighbor.X][neighbor.Y] += 1
	}
}

// Reveal a single tile. Returns whether the revealed tile was a mine.
func (b *Board) Reveal(x, y int) (isMine bool) {
	b.revealed[x][y] = true
	if b.mines[x][y] {
		return true
	} else {
		if b.neighbors[x][y] == 0 {
			b.clearZerosFrom(x, y)
		}
		return false
	}
}

// Set / remove flag for a single tile.
func (b *Board) Flag(x, y int, flag bool) {
	if !b.revealed[x][y] {
		b.flags[x][y] = flag
	}
}

// Spawn the given number of mines on the board. Returns err if impossible.
func (b *Board) SpawnMines(num int) error {
	open := make([]util.Vec, 0)
	for x := 0; x < b.size; x++ {
		for y := 0; y < b.size; y++ {
			if !b.mines[x][y] {
				open = append(open, util.Vec{X: x, Y: y})
			}
		}
	}
	if len(open) < num {
		return fmt.Errorf("insufficient empty squares to place this many mines")
	}
	open = util.Shuffle(open)
	for i := 0; i < num; i++ {
		b.PlaceMine(open[i].X, open[i].Y)
	}
	return nil
}

// From the given starting point, reveal everything cleared by zeros.
func (b *Board) clearZerosFrom(x, y int) {
	q := make([]util.Vec, 0)
	q = append(q, util.Vec{X: x, Y: y})
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		b.revealed[v.X][v.Y] = true
		if !b.mines[v.X][v.Y] && b.neighbors[v.X][v.Y] == 0 {
			for _, neighbor := range util.GetNeighbors(v.X, v.Y, b.size) {
				if !b.revealed[neighbor.X][neighbor.Y] {
					q = append(q, neighbor)
				}
			}
		}
	}
}
