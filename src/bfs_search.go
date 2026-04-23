package src

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func BFSSearch(root *Node, query string, topN int) ([]*Node, []LogEntry, int, error) {
	if root == nil {
		return nil, nil, 0, fmt.Errorf("root node is nil")
	}
	selector, err := ParseSelector(query)
	if err != nil {
		return nil, nil, 0, err
	}

	type levelResult struct {
		matched  *Node
		logEntry *LogEntry
		children []*Node
	}

	var (
		results      []*Node
		logs         []LogEntry
		nodesVisited int32
	)

	currentLevel := []*Node{root}

	for len(currentLevel) > 0 {
		levelResults := make([]levelResult, len(currentLevel))
		var wg sync.WaitGroup

		for i, node := range currentLevel {
			wg.Add(1)
			go func(idx int, n *Node) {
				defer wg.Done()
				r := levelResult{children: n.Children}
				if n.Type == ElementNode || n.Type == DocumentNode {
					atomic.AddInt32(&nodesVisited, 1)
					if n.Type == ElementNode {
						match := selector.Match(n)
						if match {
							r.matched = n
						}
						status := "visited"
						if match {
							status = "matched"
						}
						entry := LogEntry{NodeID: n.ID, Tag: n.Tag, Status: status}
						r.logEntry = &entry
					}
				}
				levelResults[idx] = r
			}(i, node)
		}

		wg.Wait()

		var nextLevel []*Node
		stop := false
		for _, r := range levelResults {
			if r.logEntry != nil {
				logs = append(logs, *r.logEntry)
			}
			if r.matched != nil {
				results = append(results, r.matched)
			}
			nextLevel = append(nextLevel, r.children...)
			if topN > 0 && len(results) >= topN {
				stop = true
				break
			}
		}

		if stop {
			break
		}
		currentLevel = nextLevel
	}

	return results, logs, int(atomic.LoadInt32(&nodesVisited)), nil
}
