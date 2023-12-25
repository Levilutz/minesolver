package util

import (
	"fmt"
	"math/rand"
)

// Copy a given list.
func ListCopy[T any](l []T) []T {
	out := make([]T, len(l))
	copy(out, l)
	return out
}

// Given a list, returns only the items that the given function considers true.
func Filter[T any](l []T, fn func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range l {
		if fn(v) {
			out = append(out, v)
		}
	}
	return out
}

func Map[T, U any](l []T, fn func(T) U) []U {
	out := make([]U, len(l))
	for i, v := range l {
		out[i] = fn(v)
	}
	return out
}

// Return a randomly shuffled permutation of a given array.
// Durstenfeld algorithm - O(n).
func Shuffle[T any](l []T) []T {
	l = ListCopy(l)
	for i := 0; i < len(l)-1; i++ {
		j := rand.Intn(len(l)-i) + i
		l[i], l[j] = l[j], l[i]
	}
	return l
}

// Flatten a double list.
func Flatten[T any](items [][]T) []T {
	outLen := 0
	for _, row := range items {
		outLen += len(row)
	}
	out := make([]T, outLen)
	i := 0
	for _, row := range items {
		for _, n := range row {
			out[i] = n
			i++
		}
	}
	return out
}

// Generate a square double array of the given type.
// If defaultValueFactor is provided, it is run to generate default
// values for each grid item.
func DArray[T any](l int, defaultValueFactory ...func() T) [][]T {
	out := make([][]T, l)
	for x := 0; x < l; x++ {
		out[x] = make([]T, l)
		if len(defaultValueFactory) > 0 {
			for y := 0; y < l; y++ {
				out[x][y] = defaultValueFactory[0]()
			}
		}
	}
	return out
}

// Generate a list of given length where values are indices.
func IndexList(l int) []int {
	out := make([]int, l)
	for i := 0; i < l; i++ {
		out[i] = i
	}
	return out
}

// A vector.
type Vec struct {
	X, Y int
}

func (v Vec) String() string {
	return fmt.Sprintf("<%d, %d>", v.X, v.Y)
}

// Get neighbors in a grid of the given coordinates.
func GetNeighbors(x, y, size int) []Vec {
	loX := x
	if loX > 0 {
		loX -= 1
	}
	loY := y
	if loY > 0 {
		loY -= 1
	}
	hiX := x
	if hiX < size-1 {
		hiX += 1
	}
	hiY := y
	if hiY < size-1 {
		hiY += 1
	}
	out := make([]Vec, (hiX-loX+1)*(hiY-loY+1)-1)
	i := 0
	for xi := loX; xi <= hiX; xi++ {
		for yi := loY; yi <= hiY; yi++ {
			if xi == x && yi == y {
				continue
			}
			out[i] = Vec{X: xi, Y: yi}
			i++
		}
	}
	return out
}
