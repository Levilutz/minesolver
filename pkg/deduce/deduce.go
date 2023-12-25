package deduce

import (
	"github.com/levilutz/minesweeper/pkg/board"
	"github.com/levilutz/minesweeper/pkg/infer"
	"github.com/levilutz/minesweeper/pkg/set"
	"github.com/levilutz/minesweeper/pkg/util"
)

// The possible numbers of mines in the given tiles.
type Fact struct {
	tiles set.Set[util.Vec]
	count set.Set[int]
}

// Whether two facts are equal.
func (f *Fact) Eq(other *Fact) bool {
	return set.IsEqual(f.tiles, other.tiles) && set.IsEqual(f.count, other.count)
}

// Whether a fact indicates a definite mine.
func (f *Fact) DefiniteMine() bool {
	return f.tiles.Size() > 0 && set.IsEqual(f.count, set.NewSet(f.tiles.Size()))
}

// Whether a fact indicates a definite empty square.
func (f *Fact) DefiniteEmpty() bool {
	return f.tiles.Size() > 0 && set.IsEqual(f.count, set.NewSet(0))
}

type Rules struct{}

// Whether two facts are equivalent.
func (Rules) Eq(a, b *Fact) bool {
	return a.Eq(b)
}

// Perform deduction on a pair of facts.
func (Rules) DeduceDual(a, b *Fact) []*Fact {
	if set.IsSubsetStrict(a.tiles, b.tiles) {
		return deduceDualSubsetStrict(a, b)
	} else if set.IsSubsetStrict(b.tiles, a.tiles) {
		return deduceDualSubsetStrict(b, a)
	}
	return nil
}

// Perform deduction on a pair of facts, where b is a strict subset of a.
func deduceDualSubsetStrict(a, b *Fact) []*Fact {
	return nil
}

// Whether two facts should be compared.
func (Rules) Relevant(a, b *Fact) bool {
	return set.Intersection(a.tiles, b.tiles).Size() > 0
}

// Whether a fact can generate an action on the board.
func (Rules) IsConclusion(f *Fact) bool {
	return f.DefiniteMine() || f.DefiniteEmpty()
}

// Compute until a single command is run.
// Returns true if command was run, or false if stuck.
func Pass(b *board.Board, maxSteps int) bool {
	// Reveal if fresh board.
	if !b.HasReveals() {
		b.Reveal(0, 0)
		return true
	}

	// Start inference
	e := infer.NewEngine[*Fact](Rules{})

	// Add the number of total unflagged mines
	remainingMines := b.UnflaggedMines()
	unknownTiles := []util.Vec{}
	for y := 0; y < b.GetSize(); y++ {
		for x := 0; x < b.GetSize(); x++ {
			if !b.HasFlag(x, y) && !b.Revealed(x, y) {
				unknownTiles = append(unknownTiles, util.Vec{X: x, Y: y})
			}
		}
	}
	e.AddFact(&Fact{
		tiles: set.FromList(unknownTiles),
		count: set.NewSet(remainingMines),
	})

	// Add a fact for each visible number
	for y := 0; y < b.GetSize(); y++ {
		for x := 0; x < b.GetSize(); x++ {
			if !b.Revealed(x, y) || b.GetNumNeighbors(x, y) == 0 {
				continue
			}
			unknown := []util.Vec{}
			unfoundMines := b.GetNumNeighbors(x, y)
			for _, neighbor := range util.GetNeighbors(x, y, b.GetSize()) {
				if !b.Revealed(neighbor.X, neighbor.Y) &&
					!b.HasFlag(neighbor.X, neighbor.Y) {
					unknown = append(unknown, util.Vec{X: neighbor.X, Y: neighbor.Y})
				}
				if b.HasFlag(neighbor.X, neighbor.Y) {
					unfoundMines -= 1
				}
			}
			if len(unknown) > 0 {
				e.AddFact(&Fact{
					tiles: set.FromList(unknown),
					count: set.NewSet(unfoundMines),
				})
			}
		}
	}

	e.Deduce(maxSteps, true)
	if e.HasConclusion() {
		c := e.Conclusions()[0]
		if c.DefiniteMine() {
			vec := c.tiles.AsList()[0]
			b.Flag(vec.X, vec.Y, true)
		} else if c.DefiniteEmpty() {
			vec := c.tiles.AsList()[0]
			b.Reveal(vec.X, vec.Y)
		} else {
			panic("expected conclusion to indicate definite mine or empty")
		}
		return true
	}
	return false
}
