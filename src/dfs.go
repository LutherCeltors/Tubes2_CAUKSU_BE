package src

import "fmt"

type LogEntry struct {
	NodeID int    `json:"nodeId"`
	Tag    string `json:"tag"`
	Status string `json:"status"`
}

func CalculateMaxDepth(root *Node) int {
	if root == nil {
		return 0
	}
	maxChildDepth := 0
	for _, child := range root.Children {
		depth := CalculateMaxDepth(child)
		if depth > maxChildDepth {
			maxChildDepth = depth
		}
	}
	return 1 + maxChildDepth
}

// SearchDFS performs Depth-First Search on the DOM tree to match CSS selectors
func SearchDFS(root *Node, query string, topN int) ([]*Node, []LogEntry, int, error) {
	if root == nil {
		return nil, nil, 0, fmt.Errorf("root node is nil")
	}

	selector, err := ParseSelector(query)
	if err != nil {
		return nil, nil, 0, err
	}

	var results []*Node
	var logs []LogEntry
	nodesVisited := 0

	var dfs func(n *Node) bool
	dfs = func(n *Node) bool {
		if n == nil {
			return false
		}

		// Only ElementNodes are truly visible/traversable for logic logs usually,
		// but DocumentNode is the root. We only log ElementNode and DocumentNode.
		if n.Type == ElementNode || n.Type == DocumentNode {
			nodesVisited++
			match := false
			if n.Type == ElementNode && selector.Match(n) {
				match = true
				results = append(results, n)
			}

			status := "visited"
			if match {
				status = "matched"
			}

			// Add to traversal log
			if n.Type == ElementNode {
				logs = append(logs, LogEntry{
					NodeID: n.ID,
					Tag:    n.Tag,
					Status: status,
				})
			}

			// If we reached TopN, stop traversing
			if topN > 0 && len(results) >= topN {
				return true
			}
		}

		// Traverse children depth-first
		for _, child := range n.Children {
			if dfs(child) {
				return true // stop signal early exit
			}
		}
		return false
	}

	dfs(root)

	return results, logs, nodesVisited, nil
}
