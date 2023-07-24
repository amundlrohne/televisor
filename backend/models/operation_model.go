package models

type OperationEdge struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Count int    `json:"count"`
}

type Operations map[string]Operation

type Operation map[string]OperationEdge

func (edges Operation) IsConnected(root string, target string, visited ...string) bool {

	if len(visited) == 0 {
		visited = []string{}
		visited = append(visited, root)
	}

	// Perf improvement, check if target exists, if not return false

	for _, v := range edges {
		if v.From == root {
			visited = append(visited, v.To)
			if v.To == target {
				return true
			}
			return edges.IsConnected(v.To, target, visited...)
		}
	}

	return false
}

func (operations Operations) ClearReflexiveEdges() Operations {
	result := make(Operations)

	for operationKey, operation := range operations {
		result[operationKey] = make(Operation)
		for edgeKey, edge := range operation {
			if edge.From != edge.To {
				result[operationKey][edgeKey] = edge
			}
		}
	}

	return result
}

func (operations Operations) CombineEdges() []OperationEdge {
	ops := operations.ClearReflexiveEdges()
	edgeMap := make(map[string]OperationEdge)

	for _, op := range ops {
		for _, e := range op {
			key := e.From + "-" + e.To
			if edge, ok := edgeMap[key]; ok {
				edge.Count += e.Count
				edgeMap[key] = edge
			} else {
				edgeMap[key] = e
			}
		}
	}

	edges := []OperationEdge{}

	for _, v := range edgeMap {
		edges = append(edges, v)
	}

	return edges
}

func (operation Operation) AddEdge(opName string, from string, to string) {
	if e, ok := operation[opName]; ok {
		e.Count++
		operation[opName] = e
	} else {
		operation[opName] = OperationEdge{
			From:  from,
			To:    to,
			Count: 1,
		}
	}
}
