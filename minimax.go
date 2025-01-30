// Minimax algorithm implementation in Go
package minimax

// SCORE is the default score for the terminal state
const SCORE = 100

// Node represents a node in the minimax tree
// T is the type of the state and must be comparable
type node[T comparable] struct {
	val      int        // Value (score) of the node
	depth    int        // Depth of the node in the tree
	isMax    bool       // True if the node is a max node
	elem     *T         // Stores game state (pointer)
	children []*node[T] // Children of the node
	bestMove *node[T]   // Best move to make (pointer)
}

// Minimax is the main struct that holds the move map (cache)
type Minimax[T comparable] struct {
	moveMap map[T]*T //Cache
}

// Solve returns the best possible move for the given state
func (m Minimax[T]) Solve(state T) *T {
	return m.moveMap[state]
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
	minimax[T](&root, isTerminal, utility, mp)

	return Minimax[T]{moveMap: mp}
}

// Minimax algorithm implementation
func minimax[T comparable](n *node[T], isTerminal func(*T) bool, utility func(*T) int, mp map[T]*T) {

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

	// Recursive minimax call
	var bestMove *node[T]
	if n.isMax {
		max := -SCORE
		for _, child := range n.children {
			minimax(child, isTerminal, utility, mp)
			eval := child.val
			if eval > max {
				max = eval
				bestMove = child
			}
		}
		n.val = max
	} else {
		min := SCORE
		for _, child := range n.children {
			minimax(child, isTerminal, utility, mp)
			eval := child.val
			if eval < min {
				min = eval
				bestMove = child
			}
		}
		n.val = min
	}

	// Store best move and update cache
	n.bestMove = bestMove
	mp[*n.elem] = n.bestMove.elem
	return
}
