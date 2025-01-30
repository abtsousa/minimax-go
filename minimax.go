package minimax

const SCORE = 100

type node[T comparable] struct {
	val      int
	depth    int
	isMax    bool
	elem     *T
	children []*node[T]
	bestMove *node[T]
}

type Minimax[T comparable] struct {
	moveMap map[T]*T //Cache
}

func (m Minimax[T]) Solve(state T) *T {
	return m.moveMap[state]
}

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

func minimax[T comparable](n *node[T], isTerminal func(*T) bool, utility func(*T) int, mp map[T]*T) {

	// Already calculated, skipping
	if n.bestMove != nil {
		return
	}

	// Terminal move found, return score
	if isTerminal(n.elem) || len(n.children) == 0 {
		n.val = utility(n.elem) * (SCORE - n.depth)
		return
	}

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
	n.bestMove = bestMove
	mp[*n.elem] = n.bestMove.elem
	return
}
