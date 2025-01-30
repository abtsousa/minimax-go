package minimax

import (
	"testing"
)

// TestMinimaxTerminalState tests the Minimax algorithm with a terminal state.
func TestMinimaxTerminalState(t *testing.T) {
	// Define a simple terminal state
	state := 0

	// Define the terminal function
	isTerminal := func(s *int) bool {
		return *s == 0
	}

	// Define the utility function
	utility := func(s *int) int {
		if *s == 0 {
			return 1 // Win for the AI
		}
		return 0
	}

	// Define the successors function (no successors for terminal state)
	successors := func(s *int) []*int {
		return []*int{}
	}

	// Create the Minimax struct
	mm := Make(&state, isTerminal, utility, successors, true)

	// Solve for the best move
	bestMove := mm.Solve(state)

	// Since it's a terminal state, the best move should be nil
	if bestMove != nil {
		t.Errorf("Expected best move to be nil for terminal state, got %v", bestMove)
	}
}

// TestMinimaxSimpleGame tests the Minimax algorithm with a simple game tree.
func TestMinimaxSimpleGame(t *testing.T) {
	// Define a simple game state
	state := 1

	// Define the terminal function
	isTerminal := func(s *int) bool {
		return *s%5 == 0 || *s > 100
	}

	// Define the utility function
	utility := func(s *int) int {
		if *s%5 == 0 {
			return 1 // Win for the AI
		}
		if *s > 100 {
			return -1 // Loss for the AI
		}
		return 0
	}

	// Define the successors function
	successors := func(s *int) []*int {
		if isTerminal(s) {
			return []*int{}
		}
		t, u := 2*(*s), 2*(*s)+1
		return []*int{&t, &u}
	}

	// Create the Minimax struct
	mm := Make(&state, isTerminal, utility, successors, true)

	// Solve for the best move
	bestMove := mm.Solve(state)

	// Check if the best move is correct
	if bestMove == nil {
		t.Error("Expected a best move, got nil")
	}
}
