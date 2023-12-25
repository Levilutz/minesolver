package set

import "fmt"

// A set of distinct elements.
type Set[K comparable] map[K]struct{}

func NewSet[K comparable](vals ...K) Set[K] {
	return FromList(vals)
}

// Create a new set from a list. If l is nil, the set is empty.
func FromList[K comparable](l []K) Set[K] {
	out := make(map[K]struct{})
	for _, v := range l {
		out[v] = struct{}{}
	}
	return out
}

// Convert the set to a list.
func (s Set[K]) AsList() []K {
	out := make([]K, len(s))
	i := 0
	for v := range s {
		out[i] = v
		i++
	}
	return out
}

// Get the size of the set.
func (s Set[K]) Size() int {
	return len(s)
}

// Return whether the set has the given value.
func (s Set[K]) Has(v K) bool {
	_, ok := s[v]
	return ok
}

func (s Set[K]) String() string {
	out := "("
	i := 0
	for v := range s {
		out += fmt.Sprint(v)
		if i != len(s)-1 {
			out += ", "
		}
		i++
	}
	return out + ")"
}

// Get the set that's the union of all provided.
func Union[K comparable](sets ...Set[K]) Set[K] {
	out := make(map[K]struct{})
	for _, set := range sets {
		for v := range set {
			out[v] = struct{}{}
		}
	}
	return out
}

// Get the set that's the intersection of all provided.
func Intersection[K comparable](sets ...Set[K]) Set[K] {
	if len(sets) == 0 {
		return make(map[K]struct{})
	}
	out := FromList(sets[0].AsList())
	for v := range out {
		for i := 1; i < len(sets); i++ {
			if !sets[i].Has(v) {
				delete(out, v)
				break
			}
		}
	}
	return out
}

// Subtract set b from set a.
func Sub[K comparable](a, b Set[K]) Set[K] {
	out := FromList(a.AsList())
	for v := range b {
		delete(out, v)
	}
	return out
}

// Return whether b is a non-strict subset of a.
func IsSubset[K comparable](a, b Set[K]) bool {
	for v := range b {
		if !a.Has(v) {
			return false
		}
	}
	return true
}

// Whether two sets are equal.
func IsEqual[K comparable](a, b Set[K]) bool {
	return IsSubset(a, b) && IsSubset(b, a)
}

// Return whether b is a strict subset of a.
func IsSubsetStrict[K comparable](a, b Set[K]) bool {
	return IsSubset(a, b) && !IsSubset(b, a)
}
