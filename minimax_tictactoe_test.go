
package minimax

import (
	"testing"
	ttt "github.com/abtsousa/tictacgo/tictactoe"
)

func TestMinimaxTicTacToe(t *testing.T) {
	tests := []struct {
		name           string
		initialState   ttt.State
		expectedResult ttt.State
	}{
		{
			name: "AI can win in one move",
			initialState: ttt.State{
				// X O X
				// O O X
				// - - -
				XBoard: 0b101_001_000,
				OBoard: 0b010_110_000,
				XPlays: false, // O's turn (AI)
			},
			expectedResult: ttt.State{ // AI should win
        // X O X
        // O O X
        // - O -
        XBoard: 0b101_001_000,
        OBoard: 0b010_110_010,
        XPlays: true, // X's turn
      },
		},
		{
			name: "AI can win or block",
			initialState: ttt.State{
				// X X -
				// O O -
				// - - -
				XBoard: 0b110_000_000,
				OBoard: 0b000_110_000,
				XPlays: false, // O's turn (AI)
			},
			expectedResult: ttt.State{
				// X X -
				// O O O
				// - - -
				XBoard: 0b110_000_000,
				OBoard: 0b000_111_000,
				XPlays: true, // X's turn
			},
		},
		{
			name: "AI can block human win",
			initialState: ttt.State{
				// X X -
				// O - -
				// - - -
				XBoard: 0b110_000_000,
				OBoard: 0b000_100_000,
				XPlays: false, // O's turn (AI)
			},
			expectedResult: ttt.State{
				// X X O
				// O - -
				// - - -
				XBoard: 0b110_000_000,
				OBoard: 0b001_100_000,
				XPlays: true, // X's turn
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize minimax with AI playing as O
			mm := Make(
				&tt.initialState,
				ttt.IsTerminal,
				ttt.Utility,
				ttt.GetSuccessors,
				true, // AI is maximizing player
			)

			// Get AI's move
			nextState := mm.Solve(tt.initialState)
			if nextState == nil {
				t.Fatal("Expected a valid move, got nil")
			}

			// Check if the move is valid
			if nextState.OBoard&tt.initialState.FreeStates() == 0 {
				t.Error("AI made an invalid move")
			}

			// Verify game result matches expected
			if *nextState != tt.expectedResult {
				t.Errorf("Expected game state %v, got %v", tt.expectedResult, nextState)
			}

		})
	}
}

func TestMinimaxTicTacToeCompleteness(t *testing.T) {
	// Test that AI can handle all possible game states
	initialState := ttt.State{
		XBoard: 0,
		OBoard: 0,
		XPlays: false,
	}

	mm := Make(
		&initialState,
		ttt.IsTerminal,
		ttt.Utility,
		ttt.GetSuccessors,
		true,
	)

	// Play through a few moves to test game progression
	currentState := initialState
	moveCount := 0
	maxMoves := 9

	for !ttt.IsTerminal(&currentState) && moveCount < maxMoves {
		nextState := mm.Solve(currentState)
		if nextState == nil {
			t.Fatal("Got nil move in non-terminal position")
		}

		// Verify move is valid
		if (nextState.XBoard&currentState.FreeStates() == 0 && nextState.OBoard == currentState.OBoard) ||
    (nextState.OBoard&currentState.FreeStates() == 0 && nextState.XBoard == currentState.XBoard) {
			t.Error("AI made invalid move")
      t.Errorf("Current state: %09b %09b %09b %v", currentState.XBoard, currentState.OBoard, currentState.FreeStates(), currentState.XPlays)
      t.Errorf("Next state: %09b %09b %09b %v", nextState.XBoard, nextState.OBoard, nextState.FreeStates(), nextState.XPlays)
		}

		currentState = *nextState
		moveCount++
	}

	if !ttt.IsTerminal(&currentState) {
		t.Error("Game did not reach terminal state")
	}
}
