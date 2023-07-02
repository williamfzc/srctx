package file

import (
	"path/filepath"
	"strconv"

	"github.com/williamfzc/srctx/graph/visual/g6"
)

const TagRed = "red"

func (fg *Graph) ToG6Data() (*g6.Data, error) {
	data := g6.EmptyG6Data()

	adjacencyMap, err := fg.G.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	// cache
	cache := make(map[string]*Vertex)
	// dir combos
	dirCombos := make(map[string]struct{})
	for nodeId := range adjacencyMap {
		node, err := fg.G.Vertex(nodeId)
		if err != nil {
			return nil, err
		}
		cache[node.Id()] = node

		eachDir := filepath.Dir(node.Path)
		dirCombos[eachDir] = struct{}{}
		data.Combos = append(data.Combos, &g6.Combo{
			Id:        eachDir,
			Label:     eachDir,
			Collapsed: false,
		})
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
		curNode := &g6.Node{
			Id:      strconv.Itoa(curId),
			Label:   node.Id(),
			Style:   &g6.NodeStyle{},
			ComboId: filepath.Dir(node.Path),
		}
		if node.ContainTag(TagRed) {
			curNode.Style.Fill = "red"
		}

		curId++
		data.Nodes = append(data.Nodes, curNode)
	}

	// Edges
	for src, targets := range adjacencyMap {
		for target := range targets {
			srcNode := cache[src]
			targetNode := cache[target]

			srcId := mapping[srcNode.Id()]
			targetId := mapping[targetNode.Id()]

			curEdge := &g6.Edge{
				Source: strconv.Itoa(srcId),
				Target: strconv.Itoa(targetId),
			}
			data.Edges = append(data.Edges, curEdge)
		}
	}

	return data, nil
}

func (fg *Graph) DrawG6Html(filename string) error {
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
