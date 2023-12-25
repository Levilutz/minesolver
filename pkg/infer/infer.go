package infer

// Something to provide basic logic functions related to a specific fact type.
type Logic[T any] interface {
	// Return whether two facts are equivalent.
	Eq(a, b T) bool

	// Return any facts derivable from a pair of given.
	DeduceDual(a, b T) []T

	// Return whether two facts should be compared.
	Relevant(a, b T) bool

	// Return whether the given fact is a terminal conclusion.
	IsConclusion(f T) bool
}

// An inference engine for a specific type of fact.
type Engine[T any] struct {
	// The logic rules used to run the engine.
	logic Logic[T]

	// The set of known facts.
	facts []T

	// The queue of deductions to run.
	deduceQ [][]T

	// Whether the engine reached a conclusion.
	hasConclusion bool

	// The conclusion of inference, if hasConclusion = true.
	conclusion T
}

// Create a new inference engine, given a set of logic functions.
func NewEngine[T any](logic Logic[T]) *Engine[T] {
	return &Engine[T]{
		logic:   logic,
		facts:   []T{},
		deduceQ: [][]T{},
	}
}

// Check if the engine contains a given fact.
func (e *Engine[T]) HasFact(f T) bool {
	for _, other := range e.facts {
		if e.logic.Eq(f, other) {
			return true
		}
	}
	return false
}

// Inform the inference engine of a new fact.
func (e *Engine[T]) AddFact(f T) {
	if e.HasFact(f) {
		return
	} else if e.logic.IsConclusion(f) {
		e.hasConclusion = true
		e.conclusion = f
		return
	}
	for _, other := range e.facts {
		if e.logic.Relevant(f, other) {
			e.deduceQ = append(e.deduceQ, []T{f, other})
		}
	}
	e.facts = append(e.facts, f)
}

// Run the given number of deductive steps, or until a final conclusion is reached.
func (e *Engine[T]) Deduce(maxSteps int) {
	for r := 0; r < maxSteps; r++ {
		if e.hasConclusion || len(e.deduceQ) == 0 {
			return
		}
		next := e.deduceQ[0]
		e.deduceQ = e.deduceQ[1:]
		var out []T = nil
		if len(next) == 2 {
			out = e.logic.DeduceDual(next[0], next[1])
		}
		for _, f := range out {
			e.AddFact(f)
			if e.hasConclusion {
				return
			}
		}
	}
}

// Check whether a final conclusion was found.
func (e *Engine[T]) HasConclusion() bool {
	return e.hasConclusion
}

// Get the conclusion fact reached by the engine.
func (e *Engine[T]) Conclusion() T {
	return e.conclusion
}
