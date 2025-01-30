
# Minimax Algorithm in Go

This repository contains a Go implementation of the Minimax algorithm with alpha-beta pruning. The Minimax algorithm is commonly used in decision-making and game theory, particularly in two-player games like chess or tic-tac-toe, to determine the optimal move for a player.

## Features

- **Alpha-Beta Pruning**: The algorithm includes the [alpha-beta pruning](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning) optimization.
- **Lazy Expansion**: Nodes are expanded only when necessary, improving memory usage.

## Usage

To use the Minimax algorithm in your project, follow these steps:

1. **Define the Game State**: Create a type that represents the state of your game. This type must be comparable.

2. **Implement Required Functions**:
   - `isTerminal`: A function that checks if a given state is a terminal state (i.e., the game is over).
   - `utility`: A function that returns the utility value of a terminal state (e.g., 1 for a win, -1 for a loss, 0 for a draw).
   - `successors`: A function that returns the possible successor states from a given state.

3. **Create a Minimax Instance**: Use the `Make` function to create a Minimax instance with the initial state and the functions defined above.

4. **Solve for the Best Move**: Call the `Solve` method on the Minimax instance to get the best move for the current state.

### Example

You can check [a tictactoe implementation in Go in my Github](https://github.com/abtsousa/tictacgo).

```go
package main

import (
	"fmt"
	"github.com/abtsousa/minimax-go"
)

func main() {
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
	mm := minimax.Make(&state, isTerminal, utility, successors, true)

	// Solve for the best move
	bestMove := mm.Solve(state)

	fmt.Println("Best move:", *bestMove)
}
```
## License

This project is licensed under the GPLv3 License. See the LICENSE file for details.

## Contributing

Contributions are welcome. Please open an issue or submit a pull request for any improvements or bug fixes.
