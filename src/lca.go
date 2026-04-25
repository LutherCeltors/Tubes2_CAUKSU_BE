package src

import "fmt"

type LCABinaryLifting struct {
	Root  *Node
	LOG   int
	Up    map[*Node][]*Node
	Depth map[*Node]int
	ByID  map[int]*Node
}

func PreproccessLCABinaryLifting(root *Node) (*LCABinaryLifting, error) {
	if root == nil {
		return nil, fmt.Errorf("root node is nil")
	}

	nodeCount := 0

	TraverseDFS(root, func(n *Node, parent *Node, depth int) bool {
		nodeCount++
		return true
	})

	log := 1
	for (1 << log) <= nodeCount {
		log++
	}

	lca := &LCABinaryLifting{
		Root:  root,
		LOG:   log,
		Up:    make(map[*Node][]*Node),
		Depth: make(map[*Node]int),
		ByID:  make(map[int]*Node),
	}

	TraverseDFS(root, func(n *Node, parent *Node, depth int) bool {
		lca.Depth[n] = depth
		lca.ByID[n.ID] = n

		lca.Up[n] = make([]*Node, log)

		if parent == nil {
			for j := 0; j < log; j++ {
				lca.Up[n][j] = n
			}
		} else {
			lca.Up[n][0] = parent

			for j := 1; j < log; j++ {
				midAncestor := lca.Up[n][j-1]
				lca.Up[n][j] = lca.Up[midAncestor][j-1]
			}
		}

		return true
	})

	return lca, nil
}

func (lca *LCABinaryLifting) Lift(n *Node, k int, logs *[]LogEntry, visited map[*Node]bool, nodesVisited *int, batchIndex *int,
) *Node {
	if n == nil {
		return nil
	}

	for j := 0; j < lca.LOG; j++ {
		if k&(1<<j) != 0 {
			n = lca.Up[n][j]

			appendLCALog(logs, visited, nodesVisited, n, "visited", *batchIndex)

			(*batchIndex)++
		}
	}

	return n
}

func (lca *LCABinaryLifting) SearchLCA(a *Node, b *Node) ([]*Node, []LogEntry, int, error) {
	var results []*Node
	var logs []LogEntry

	batchIndex := 0
	nodesVisited := 0
	visited := make(map[*Node]bool)

	if a == nil || b == nil {
		return nil, nil, 0, fmt.Errorf("node is nil")
	}

	if _, ok := lca.Depth[a]; !ok {
		return nil, nil, 0, fmt.Errorf("first node is not registered in LCA table")
	}

	if _, ok := lca.Depth[b]; !ok {
		return nil, nil, 0, fmt.Errorf("second node is not registered in LCA table")
	}

	if lca.Depth[a] < lca.Depth[b] {
		a, b = b, a
	}

	appendLCALog(&logs, visited, &nodesVisited, a, "visited", batchIndex)
	appendLCALog(&logs, visited, &nodesVisited, b, "visited", batchIndex)
	batchIndex++

	if a == b {
		appendLCALog(&logs, visited, &nodesVisited, a, "matched", batchIndex)
		results = append(results, a)
		return results, logs, nodesVisited, nil
	}

	depthDiff := lca.Depth[a] - lca.Depth[b]

	if depthDiff > 0 {
		a = lca.Lift(a, depthDiff, &logs, visited, &nodesVisited, &batchIndex)

		if a == b {
			appendLCALog(&logs, visited, &nodesVisited, a, "matched", batchIndex)
			results = append(results, a)
			return results, logs, nodesVisited, nil
		}
	}

	for j := lca.LOG - 1; j >= 0; j-- {
		if lca.Up[a][j] != lca.Up[b][j] {
			a = lca.Up[a][j]
			b = lca.Up[b][j]

			appendLCALog(&logs, visited, &nodesVisited, a, "visited", batchIndex)
			appendLCALog(&logs, visited, &nodesVisited, b, "visited", batchIndex)

			batchIndex++
		}
	}

	ancestor := lca.Up[a][0]

	appendLCALog(&logs, visited, &nodesVisited, ancestor, "matched", batchIndex)
	results = append(results, ancestor)

	return results, logs, nodesVisited, nil
}

func (lca *LCABinaryLifting) SearchLCAByID(aID int, bID int) ([]*Node, []LogEntry, int, error) {
	a, ok := lca.ByID[aID]
	if !ok {
		return nil, nil, 0, fmt.Errorf("node with ID %d not found", aID)
	}

	b, ok := lca.ByID[bID]
	if !ok {
		return nil, nil, 0, fmt.Errorf("node with ID %d not found", bID)
	}

	return lca.SearchLCA(a, b)
}

func isLoggableLCANode(n *Node) bool {
	return n != nil && (n.Type == ElementNode || n.Type == DocumentNode)
}

func appendLCALog(
	logs *[]LogEntry,
	visited map[*Node]bool,
	nodesVisited *int,
	n *Node,
	status string,
	batchIndex int,
) {
	if !isLoggableLCANode(n) {
		return
	}

	if !visited[n] {
		visited[n] = true
		(*nodesVisited)++
	}

	*logs = append(*logs, LogEntry{
		NodeID: n.ID,
		Tag:    n.Tag,
		Status: status,
		Batch:  batchIndex,
	})
}