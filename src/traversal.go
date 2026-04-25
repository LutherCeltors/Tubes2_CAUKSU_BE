package src

type DFSVisitFunc func(n *Node, parent *Node, depth int) bool

func TraverseDFS(root *Node, visit DFSVisitFunc) {
	traverseDFS(root, nil, 0, visit)
}

func traverseDFS(n *Node, parent *Node, depth int, visit DFSVisitFunc) bool {
	if n == nil {
		return true
	}

	if visit != nil {
		shouldContinue := visit(n, parent, depth)
		if !shouldContinue {
			return false
		}
	}

	for _, child := range n.Children {
		if !traverseDFS(child, n, depth+1, visit) {
			return false
		}
	}

	return true
}