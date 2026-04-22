package src

import "fmt"

func BFSSearch(root *Node, query string, topN int) ([]*Node, []LogEntry, int, error) {
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

	queue := []*Node{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.Type == ElementNode || node.Type == DocumentNode {
			nodesVisited++
			match := false
			if node.Type == ElementNode && selector.Match(node) {
				match = true
				results = append(results, node)
			}
			if node.Type == ElementNode {
				status := "visited"
				if match {
					status = "matched"
				}
				logs = append(logs, LogEntry{
					NodeID: node.ID,
					Tag:    node.Tag,
					Status: status,
				})
			}
			if topN > 0 && len(results) >= topN {
				break
			}
		}

		queue = append(queue, node.Children...)
	}

	return results, logs, nodesVisited, nil
}