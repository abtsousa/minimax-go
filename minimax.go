// Package minimax provides an implementation of the Minimax algorithm with alpha-beta pruning.
// The Minimax algorithm is used in decision-making and game theory to determine the optimal move for a player.
//
// The package includes features such as alpha-beta pruning and lazy node expansion to optimize performance and memory usage.
//
// Usage:
//
// 1. Define the game state by creating a type that represents the state of your game.
//
// 2. Implement the required functions:
// - state: the initial gamestate
// - isTerminal: a function that returns true if the state is terminal
// - utility: a function that should return -1 if the state is a loss for the AI, 1 if it's a win and 0 if it's a draw
// - successors: a function that returns the possible moves from the state
// - isMax: true if the initial state is a max node (AI's turn)
//
// 3. Create a Minimax instance using the `Make` function.
//
// 4. Solve for the best move using the `Solve` method.
//
// Example:
//
//	mm := minimax.Make(&state, isTerminal, utility, successors, true)
//	bestMove := mm.Solve(state)
package minimax

// score is the default score for the terminal state
const score = 100

// Node represents a node in the minimax tree
// T is the type of the state and must be comparable
type node[T comparable] struct {
	val      int        // Value (score) of the node
	alpha    int        // Alpha value for alpha-beta pruning
	beta     int        // Beta value for alpha-beta pruning
	depth    int        // Depth of the node in the tree
	elem     *T         // Stores game state (pointer)
	children []*node[T] // Children of the node (generated lazily)
	bestMove *node[T]   // Best move to make (pointer)
	isMax    bool       // True if the node is a max node
	expanded bool       // Whether children have been generated
}

// Minimax is the main struct that holds the move map (cache)
type Minimax[T comparable] struct {
	moveMap map[T]*T // Cache
	config  struct {
		isTerminal func(*T) bool
		utility    func(*T) int
		successors func(*T) []*T
		isMax      bool
	}
}

// Solve returns the best possible move for the given state
func (m Minimax[T]) Solve(state T) *T {
	if m.config.isTerminal(&state) {
		return nil
	}

	bestMove := m.moveMap[state]
	if bestMove != nil {
		return bestMove
	}

	// No best move found, possibly pruned tree (from suboptimal move)
	// Rerun algorithm to find best move
	cf := m.config
	newMM := Make(&state, cf.isTerminal, cf.utility, cf.successors, cf.isMax)
	m.moveMap = newMM.moveMap
	return m.Solve(state)
}

// Make creates a new Minimax struct. You must provide:
// - state: the initial gamestate
// - isTerminal: a function that returns true if the state is terminal
// - utility: a function that should return -1 if the state is a loss for the AI, 1 if it's a win and 0 if it's a draw
// - successors: a function that returns the possible moves from the state
// - isMax: true if the initial state is a max node (AI's turn)
func Make[T comparable](state *T, isTerminal func(*T) bool,
	utility func(*T) int, successors func(*T) []*T, isMax bool,
) Minimax[T] {
	root := &node[T]{
		val:      0,
		alpha:    -score,
		beta:     score,
		depth:    0,
		isMax:    isMax,
		elem:     state,
		expanded: false,
	}

	mp := make(map[T]*T)
	minimax(root, isTerminal, utility, successors, mp)

	return Minimax[T]{
		moveMap: mp,
		config: struct {
			isTerminal func(*T) bool
			utility    func(*T) int
			successors func(*T) []*T
			isMax      bool
		}{
			isTerminal: isTerminal,
			utility:    utility,
			successors: successors,
			isMax:      isMax,
		},
	}
}

// expandNode generates children nodes only when needed
func expandNode[T comparable](n *node[T], successors func(*T) []*T) {
	if n.expanded {
		return
	}

	successorStates := successors(n.elem)
	n.children = make([]*node[T], 0, len(successorStates))

	for _, succ := range successorStates {
		child := &node[T]{
			val:      0,
			alpha:    -score,
			beta:     score,
			depth:    n.depth + 1,
			isMax:    !n.isMax,
			elem:     succ,
			expanded: false,
		}
		n.children = append(n.children, child)
	}

	n.expanded = true
}

func minimax[T comparable](n *node[T], isTerminal func(*T) bool,
	utility func(*T) int, successors func(*T) []*T, mp map[T]*T,
) {
	// Best move already calculated, skipping
	if n.bestMove != nil {
		return
	}

	// Terminal move found, return score
	if isTerminal(n.elem) {
		switch u := utility(n.elem); {
		case u > 0:
			n.val = score - n.depth
		case u < 0:
			n.val = n.depth - score
		default:
			n.val = 0
		}
		return
	}

	// Lazily expand node
	expandNode(n, successors)

	// If no children after expansion, treat as terminal
	if len(n.children) == 0 {
		n.val = utility(n.elem)
		return
	}

	var bestMove *node[T]
	if n.isMax {
		maxEval := -score
		for _, child := range n.children {
			child.alpha = n.alpha
			child.beta = n.beta

			minimax(child, isTerminal, utility, successors, mp)
			eval := child.val
			if eval > maxEval {
				maxEval = eval
				bestMove = child
			}
			n.alpha = max(n.alpha, maxEval)

			if n.beta <= n.alpha {
				break // Beta cutoff
			}
		}
		n.val = maxEval
	} else {
		minEval := score
		for _, child := range n.children {
			child.alpha = n.alpha
			child.beta = n.beta

			minimax(child, isTerminal, utility, successors, mp)
			eval := child.val
			if eval < minEval {
				minEval = eval
				bestMove = child
			}
			n.beta = min(n.beta, minEval)

			if n.beta <= n.alpha {
				break // Alpha cutoff
			}
		}
		n.val = minEval
	}

	n.bestMove = bestMove
	mp[*n.elem] = n.bestMove.elem
}
