package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	board "github.com/levilutz/minesweeper/pkg/board"
	"github.com/levilutz/minesweeper/pkg/textrender"
)

func main() {
	boardSize := 8
	numMines := 10

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	b := board.NewBoard(boardSize)
	if err := b.SpawnMines(numMines); err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		if b.Complete() {
			fmt.Println("you win!")
		}

		fmt.Println(textrender.RenderBoard(b))
		fmt.Print(": ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input = strings.TrimRight(input, "\n")
		cmd := strings.Split(input, " ")

		if len(cmd) == 0 {
			continue

		} else if cmd[0] == "exit" {
			return

		} else if cmd[0] == "flag" || cmd[0] == "f" {
			if len(cmd) < 3 {
				fmt.Println("must provide x and y coordinates")
				continue
			}
			x, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Printf("failed to parse: %s\n", err)
				continue
			}
			y, err := strconv.Atoi(cmd[2])
			if err != nil {
				fmt.Printf("failed to parse: %s\n", err)
				continue
			}
			if b.HasFlag(x, y) {
				b.Flag(x, y, false)
				fmt.Printf("removed flag from (%d, %d)\n", x, y)
			} else {
				b.Flag(x, y, true)
				fmt.Printf("added flag to (%d, %d)\n", x, y)
			}

		} else if cmd[0] == "reveal" || cmd[0] == "r" {
			if len(cmd) < 3 {
				fmt.Println("must provide x and y coordinates")
				continue
			}
			x, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Printf("failed to parse: %s\n", err)
				continue
			}
			y, err := strconv.Atoi(cmd[2])
			if err != nil {
				fmt.Printf("failed to parse: %s\n", err)
				continue
			}
			if b.HasFlag(x, y) {
				fmt.Println("cannot reveal tile with flag")
				continue
			}
			// If fresh game, re-generate until reveal allowed.
			if !b.HasReveals() {
				for {
					isMine := b.Reveal(x, y)
					if isMine {
						b.Reset()
						b.SpawnMines(numMines)
					} else {
						break
					}
				}
			} else {
				isMine := b.Reveal(x, y)
				if isMine {
					fmt.Println("tile has mine, you lose!")
					b.Reset()
					b.SpawnMines(numMines)
				} else {
					fmt.Printf("revealed (%d, %d)\n", x, y)
				}
			}

		} else if cmd[0] == "reset" {
			b.Reset()
			b.SpawnMines(numMines)

		} else {
			fmt.Printf("unknown command: %s\n", cmd[0])
		}
	}
}
