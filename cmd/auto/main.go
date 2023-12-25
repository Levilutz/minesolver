package main

import (
	"fmt"
	"time"

	"github.com/levilutz/minesweeper/pkg/board"
	"github.com/levilutz/minesweeper/pkg/solver"
	"github.com/levilutz/minesweeper/pkg/textrender"
)

func ViewOne() {
	boardSize := 16
	numMines := 40
	delay := time.Millisecond * 50

	var b *board.Board
	for {
		b = board.NewBoard(boardSize)
		b.SpawnMines(numMines)
		if !b.HasMine(0, 0) {
			break
		}
	}

	for {
		fmt.Println(textrender.RenderBoard(b))
		tookAction := solver.Pass(b)
		if b.Complete() {
			fmt.Println(textrender.RenderBoard(b))
			fmt.Println("solver won!")
			break
		} else if !tookAction {
			fmt.Println("solver stuck")
			break
		}
		time.Sleep(delay)
	}
}

// Run numRounds tests of the solver, return the number of successes.
func TestRounds(numRounds int) int {
	boardSize := 16
	numMines := 40
	allowedSteps := boardSize * boardSize * 2

	wins := 0
	for round := 0; round < numRounds; round++ {
		b := board.NewBoard(boardSize)
		b.SpawnMines(numMines)
		if b.HasMine(0, 0) {
			round--
			continue
		}
		var tookAction bool
		for i := 0; i < allowedSteps; i++ {
			tookAction = solver.Pass(b)
			if b.Complete() && !tookAction {
				break
			}
		}
		if b.Complete() {
			wins += 1
		}
	}
	return wins
}

func main() {
	// fmt.Printf("%d / 10000", TestRounds(10000))
	ViewOne()
}
