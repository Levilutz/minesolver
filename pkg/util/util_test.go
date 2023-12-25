package util_test

import (
	"fmt"
	"testing"

	"github.com/levilutz/minesweeper/pkg/util"
)

func TestShuffle(t *testing.T) {
	tblSize := 10
	rounds := 100000
	printOutput := true

	freq := make([][]int, tblSize)
	for i := 0; i < tblSize; i++ {
		freq[i] = make([]int, tblSize)
	}
	for r := 0; r < rounds; r++ {
		arr := util.IndexList(tblSize)
		arr = util.Shuffle(arr)
		for i := 0; i < tblSize; i++ {
			freq[i][arr[i]] += 1
		}
	}

	if printOutput {
		for _, row := range freq {
			for _, n := range row {
				fmt.Printf("%.2f", float64(n*tblSize)/float64(rounds))
				fmt.Print("\t")
			}
			fmt.Println()
		}
	}

	deviations := 0
	for _, row := range freq {
		for _, n := range row {
			expected := rounds / tblSize
			if float64(n) < float64(expected)*.9 || float64(n) > float64(expected)*1.1 {
				deviations++
			}
		}
	}

	if float64(deviations) > .01*float64(tblSize*tblSize) {
		t.Fatalf("got excessive deviations from even distribution: %d", deviations)
	}
}
