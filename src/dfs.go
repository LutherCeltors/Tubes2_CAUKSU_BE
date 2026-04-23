package src

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type LogEntry struct {
	NodeID int    `json:"nodeId"`
	Tag    string `json:"tag"`
	Status string `json:"status"`
}

func SearchDFS(root *Node, query string, topN int) ([]*Node, []LogEntry, int, error) {
	if root == nil {
		return nil, nil, 0, fmt.Errorf("root node is nil")
	}
	selector, err := ParseSelector(query)
	if err != nil {
		return nil, nil, 0, err
	}

	type subtreeResult struct {
		results []*Node
		logs    []LogEntry
		visited int
	}

	var totalFound int32

	var dfs func(n *Node) subtreeResult
	dfs = func(n *Node) subtreeResult {
		if n == nil {
			return subtreeResult{}
		}
		if topN > 0 && atomic.LoadInt32(&totalFound) >= int32(topN) {
			return subtreeResult{}
		}

		var r subtreeResult

		if n.Type == ElementNode || n.Type == DocumentNode {
			r.visited = 1
			if n.Type == ElementNode {
				match := selector.Match(n)
				if match {
					atomic.AddInt32(&totalFound, 1)
					r.results = []*Node{n}
				}
				status := "visited"
				if match {
					status = "matched"
				}
				r.logs = []LogEntry{{NodeID: n.ID, Tag: n.Tag, Status: status}}
			}
		}

		if len(n.Children) == 0 {
			return r
		}

		childResults := make([]subtreeResult, len(n.Children))
		var wg sync.WaitGroup
		for i, child := range n.Children {
			wg.Add(1)
			go func(idx int, c *Node) {
				defer wg.Done()
				childResults[idx] = dfs(c)
			}(i, child)
		}
		wg.Wait()

		for _, cr := range childResults {
			r.results = append(r.results, cr.results...)
			r.logs = append(r.logs, cr.logs...)
			r.visited += cr.visited
		}

		return r
	}

	res := dfs(root)

	results := res.results
	if topN > 0 && len(results) > topN {
		results = results[:topN]
	}

	return results, res.logs, res.visited, nil
}

type JSONNode struct {
	ID         int               `json:"id"`
	Tag        string            `json:"tag"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Children   []*JSONNode       `json:"children,omitempty"`
}

func ConvertToJSONNode(n *Node) *JSONNode {
	if n == nil {
		return nil
	}
	if n.Type != ElementNode && n.Type != DocumentNode {
		return nil
	}
	attrs := make(map[string]string)
	for _, a := range n.Attrs {
		attrs[a.Name] = a.Value
	}
	tag := n.Tag
	if n.Type == DocumentNode {
		tag = "document"
	}
	res := &JSONNode{
		ID:         n.ID,
		Tag:        tag,
		Attributes: attrs,
	}
	for _, child := range n.Children {
		if cJson := ConvertToJSONNode(child); cJson != nil {
			res.Children = append(res.Children, cJson)
		}
	}
	return res
}
