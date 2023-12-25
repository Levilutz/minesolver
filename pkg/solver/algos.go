package solver

import (
	"fmt"

	"github.com/levilutz/minesweeper/pkg/board"
	"github.com/levilutz/minesweeper/pkg/set"
	"github.com/levilutz/minesweeper/pkg/util"
)

const debug = false

func hitMine(x, y int) {
	panic(fmt.Sprintf("solver hit mine at (%d, %d)\n", x, y))
}

// If the game is fresh, reveal a random tile.
// Returns true if action was taken.
func revealIfFresh(b *board.Board) bool {
	if !b.HasReveals() {
		if b.Reveal(0, 0) {
			hitMine(0, 0)
		}
		return true
	}
	return false
}

// Find any tiles that are obviously a mine.
// Returns true if action was taken.
func findObviousMines(b *board.Board) bool {
	size := b.GetSize()

	findDefiniteFlags := func(x, y int) bool {
		numNeighbors := b.GetNumNeighbors(x, y)
		if numNeighbors == 0 {
			return false
		}
		numUnrevealedNeighbors := 0
		numUnrevealedUnflaggedNeighbors := 0
		neighbors := util.GetNeighbors(x, y, size)
		for _, neighbor := range neighbors {
			if !b.Revealed(neighbor.X, neighbor.Y) {
				numUnrevealedNeighbors += 1
				if !b.HasFlag(neighbor.X, neighbor.Y) {
					numUnrevealedUnflaggedNeighbors += 1
				}
			}
		}
		if numUnrevealedNeighbors == numNeighbors &&
			numUnrevealedUnflaggedNeighbors > 0 {
			for _, neighbor := range neighbors {
				if !b.Revealed(neighbor.X, neighbor.Y) &&
					!b.HasFlag(neighbor.X, neighbor.Y) {
					if debug {
						fmt.Printf("solver flagging (%d, %d)\n", neighbor.X, neighbor.Y)
					}
					b.Flag(neighbor.X, neighbor.Y, true)
					return true
				}
			}
		}
		return false
	}
	return forEachRevealed(b, findDefiniteFlags)
}

// Find any tiles that are obviously empty.
// Returns true if action was taken.
func findObiousEmpty(b *board.Board) bool {
	size := b.GetSize()

	findDefiniteEmpty := func(x, y int) bool {
		numNeighbors := b.GetNumNeighbors(x, y)
		neighbors := util.GetNeighbors(x, y, size)
		numFlaggedNeighbors := 0
		numUnrevealedNeighbors := 0
		for _, neighbor := range neighbors {
			if b.HasFlag(neighbor.X, neighbor.Y) {
				numFlaggedNeighbors += 1
			}
			if !b.Revealed(neighbor.X, neighbor.Y) {
				numUnrevealedNeighbors += 1
			}
		}
		if numNeighbors == numFlaggedNeighbors &&
			numUnrevealedNeighbors > numFlaggedNeighbors {
			for _, neighbor := range neighbors {
				if !b.Revealed(neighbor.X, neighbor.Y) &&
					!b.HasFlag(neighbor.X, neighbor.Y) {
					if debug {
						fmt.Printf(
							"solver revealing (%d, %d)\n", neighbor.X, neighbor.Y,
						)
					}
					if b.Reveal(neighbor.X, neighbor.Y) {
						hitMine(neighbor.X, neighbor.Y)
					}
					return true
				}
			}
		}
		return false
	}
	return forEachRevealed(b, findDefiniteEmpty)
}

// The fact that the given set contains the given number of mines.
type Fact struct {
	mines int
	tiles set.Set[util.Vec]
}

func (f *Fact) String() string {
	return fmt.Sprintf("{%d in %s}", f.mines, f.tiles)
}

// Knowledge graph.
type Knowledge struct {
	// A ref to the board (to take actions on)
	b *board.Board

	// The full set of all current facts
	nodes []*Fact

	// The facts relevant to each tile on the grid
	tiles [][][]*Fact

	// A queue of fact pairs that have not had deduction run yet
	unchecked [][]*Fact
}

func NewKnowledge(b *board.Board) *Knowledge {
	return &Knowledge{
		b:     b,
		nodes: make([]*Fact, 0),
		tiles: util.DArray[[]*Fact](b.GetSize(), func() []*Fact {
			return make([]*Fact, 0)
		}),
		unchecked: make([][]*Fact, 0),
	}
}

// Whether the given number of mines / tiles is already known.
func (k *Knowledge) HasFact(mines int, vecs set.Set[util.Vec]) bool {
	for _, node := range k.nodes {
		if node.mines == mines && set.IsEqual(node.tiles, vecs) {
			return true
		}
	}
	return false
}

func (k *Knowledge) AddFact(mines int, vecs set.Set[util.Vec]) {
	if k.HasFact(mines, vecs) {
		return
	}
	fact := Fact{
		mines: mines,
		tiles: vecs,
	}
	if debug {
		fmt.Printf("+ %s\n", &fact)
	}
	k.nodes = append(k.nodes, &fact)
	for _, vec := range vecs.AsList() {
		relevant := k.tiles[vec.X][vec.Y]
		for _, otherFact := range relevant {
			k.unchecked = append(k.unchecked, []*Fact{&fact, otherFact})
		}
		k.tiles[vec.X][vec.Y] = append(k.tiles[vec.X][vec.Y], &fact)
	}
}

func (k *Knowledge) HasUncheckedDeductions() bool {
	return len(k.unchecked) > 0
}

// Run the next deduction from the queue.
// Returns true if an action was taken.
func (k *Knowledge) RunNextDeduction() bool {
	if len(k.unchecked) == 0 {
		return false
	}
	next := k.unchecked[0]
	k.unchecked = k.unchecked[1:]
	if len(next) == 2 {
		return k.RunDualDeduction(next[0], next[1]) ||
			k.RunDualDeduction(next[1], next[0])
	} else {
		// TODO: break up into pairs, or n-deduction (if not single fact)
		fmt.Printf("cannot run deduction on %d facts\n", len(next))
		return false
	}
}

// Run a deduction on a pair of facts.
// Returns true if an action was taken.
func (k *Knowledge) RunDualDeduction(a, b *Fact) bool {
	if set.IsEqual(a.tiles, b.tiles) && a.mines != b.mines {
		panic(fmt.Sprintf("contradiction between %s & %s", a, b))
	}
	if set.IsSubsetStrict(a.tiles, b.tiles) {
		if b.mines == a.mines {
			// One-step deduction solves 21.5% of 8x8 w/ 10 mines
			// Tiles are clearable!
			empty := set.Sub(a.tiles, b.tiles).AsList()
			if len(empty) == 0 {
				panic(fmt.Sprintf("expected empty tiles from %s - %s", a, b))
			}
			if k.b.Reveal(empty[0].X, empty[0].Y) {
				hitMine(empty[0].X, empty[0].Y)
			}
			return true
		} else if b.mines < a.mines {
			// Multi-step deduction solves 30% of 8x8 w/ 10 mines
			subZoneMines := a.mines - b.mines
			subZoneTiles := set.Sub(a.tiles, b.tiles)
			if subZoneMines > 0 && subZoneMines == len(subZoneTiles) {
				// Definite mines!
				tile := subZoneTiles.AsList()[0]
				k.b.Flag(tile.X, tile.Y, true)
				return true
			}
			if debug {
				fmt.Println("from deduction")
			}
			k.AddFact(
				subZoneMines,
				subZoneTiles,
			)
		}
	}
	return false
}

// Accumulate facts about the state of the board and deduce.
// Returns true if action was taken.
func deduce(b *board.Board) bool {
	maxDeductions := 100000

	size := b.GetSize()
	know := NewKnowledge(b)

	unrevealedNeighbors := func(x, y int) set.Set[util.Vec] {
		out := make([]util.Vec, 0)
		for _, neighbor := range util.GetNeighbors(x, y, size) {
			if !b.Revealed(neighbor.X, neighbor.Y) {
				out = append(out, neighbor)
			}
		}
		return set.FromList(out)
	}

	flaggedNeighbors := func(x, y int) set.Set[util.Vec] {
		out := make([]util.Vec, 0)
		for _, neighbor := range util.GetNeighbors(x, y, size) {
			if b.HasFlag(neighbor.X, neighbor.Y) {
				out = append(out, neighbor)
			}
		}
		return set.FromList(out)
	}

	// Accumulate a fact for each visible number.
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if b.Revealed(x, y) && b.GetNumNeighbors(x, y) > 0 {
				unrevealed := unrevealedNeighbors(x, y)
				flagged := flaggedNeighbors(x, y)
				if unrevealed.Size() > flagged.Size() {
					if debug {
						fmt.Println("from board")
					}
					know.AddFact(
						b.GetNumNeighbors(x, y)-flagged.Size(),
						set.Sub(unrevealedNeighbors(x, y), flagged),
					)
				}
			}
		}
	}

	// Continuously attempt deductions
	for i := 0; i < maxDeductions; i++ {
		if !know.HasUncheckedDeductions() {
			return false
		}
		if know.RunNextDeduction() {
			return true
		}
	}

	fmt.Println("executed max deductions")

	return false
}
