// Minimax algorithm implementation in Go
package minimax

// SCORE is the default score for the terminal state
const SCORE = 100

// Node represents a node in the minimax tree
// T is the type of the state and must be comparable
type node[T comparable] struct {
	val      int        // Value (score) of the node
	alpha    int        // Alpha value for alpha-beta pruning
	beta     int        // Beta value for alpha-beta pruning
	depth    int        // Depth of the node in the tree
	isMax    bool       // True if the node is a max node
	elem     *T         // Stores game state (pointer)
	children []*node[T] // Children of the node
	bestMove *node[T]   // Best move to make (pointer)
}

// Minimax is the main struct that holds the move map (cache)
type Minimax[T comparable] struct {
	moveMap map[T]*T //Cache
	config  struct {
		isTerminal func(*T) bool
		utility    func(*T) int
		successors func(*T) []*T
		isMax      bool
	}
}

// Solve returns the best possible move for the given state
func (m Minimax[T]) Solve(state T) *T {
	bestMove := m.moveMap[state]
	if bestMove != nil {
		return bestMove
	}

	// No best move found, possibly pruned tree (from suboptimal move)
	// Rerun algorithm to find best move
	newMM := Make(&state, m.config.isTerminal, m.config.utility, m.config.successors, m.config.isMax)
	m.moveMap = newMM.moveMap
	return m.Solve(state)
}

// Make creates a new Minimax struct. You must provide:
// - state: the initial gamestate
// - isTerminal: a function that returns true if the state is terminal
// - utility: a function that should return -1 if the state is a loss for the AI,
// 1 if it's a win and 0 if it's a draw
// - successors: a function that returns the possible moves from the state
// - isMax: true if the initial state is a max node (AI's turn)
func Make[T comparable](state *T, isTerminal func(*T) bool, utility func(*T) int, successors func(*T) []*T, isMax bool) Minimax[T] {

	var makeNode func(*T, int, bool) node[T]
	makeNode = func(state *T, depth int, isMax bool) node[T] {

		successors := successors(state)

		var children []*node[T]
		for _, succ := range successors {
			n := makeNode(succ, depth+1, !isMax)
			children = append(children, &n)
		}

		n := node[T]{
			val:      0,
			alpha:    -SCORE,
			beta:     SCORE,
			depth:    depth,
			isMax:    isMax,
			elem:     state,
			children: children,
			bestMove: nil,
		}

		return n
	}

	mp := make(map[T]*T)
	root := makeNode(state, 0, isMax)
	minimax(&root, isTerminal, utility, mp)

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

// Minimax algorithm implementation
func minimax[T comparable](n *node[T], isTerminal func(*T) bool, utility func(*T) int, mp map[T]*T) {

	// Initialize alpha and beta for root node
	n.alpha = -SCORE
	n.beta = SCORE

	// Best move already calculated, skipping
	if n.bestMove != nil {
		return
	}

	// Terminal move found, return score
	if isTerminal(n.elem) || len(n.children) == 0 {
		switch u := utility(n.elem); {
		case u > 0:
			n.val = SCORE - n.depth
		case u < 0:
			n.val = n.depth - SCORE
		default:
			n.val = 0
		}
		return
	}

	var bestMove *node[T]
	if n.isMax {
		max_eval := -SCORE
		for _, child := range n.children {
			// Pass down alpha and beta values
			child.alpha = n.alpha
			child.beta = n.beta

			// Recursive minimax call
			minimax(child, isTerminal, utility, mp)
			eval := child.val
			if eval > max_eval {
				max_eval = eval
				bestMove = child
			}
			n.alpha = max(n.alpha, max_eval)

			// Prune if possible
			if n.beta <= n.alpha {
				break // Beta cutoff
			}
		}
		n.val = max_eval
	} else {
		min_eval := SCORE
		for _, child := range n.children {
			// Pass down alpha and beta values
			child.alpha = n.alpha
			child.beta = n.beta

			minimax(child, isTerminal, utility, mp)
			eval := child.val
			if eval < min_eval {
				min_eval = eval
				bestMove = child
			}
			n.beta = min(n.beta, min_eval)

			if n.beta <= n.alpha {
				break // Alpha cutoff
			}

		}
		n.val = min_eval
	}

	// Store best move and update cache
	n.bestMove = bestMove
	mp[*n.elem] = n.bestMove.elem
	return
}
