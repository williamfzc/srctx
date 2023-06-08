package file

import (
	"strconv"

	"github.com/williamfzc/srctx/graph/visual/g6"
)

func (fg *FileGraph) ToG6Data() (*g6.G6Data, error) {
	data := g6.EmptyG6Data()

	adjacencyMap, err := fg.G.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	// mapping
	mapping := make(map[string]int)
	curId := 0

	// Nodes
	for nodeId := range adjacencyMap {
		node, err := fg.G.Vertex(nodeId)
		if err != nil {
			return nil, err
		}
		mapping[node.Id()] = curId
		curNode := &g6.G6Node{
			Id:    strconv.Itoa(curId),
			Label: node.Path,
			Style: &g6.G6NodeStyle{},
		}
		curId++
		data.Nodes = append(data.Nodes, curNode)
	}
	// Edges
	for src, targets := range adjacencyMap {
		for target := range targets {
			srcId := mapping[src]
			targetId := mapping[target]

			curEdge := &g6.G6Edge{
				Source: strconv.Itoa(srcId),
				Target: strconv.Itoa(targetId),
			}
			data.Edges = append(data.Edges, curEdge)
		}
	}
	return data, nil
}

func (fg *FileGraph) DrawG6Html(filename string) error {
	data, err := fg.ToG6Data()
	if err != nil {
		return err
	}
	err = data.RenderHtml(filename)
	if err != nil {
		return err
	}

	return nil
}
