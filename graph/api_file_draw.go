package graph

import (
	"os"
	"strconv"

	"github.com/dominikbraun/graph/draw"
)

func (fg *FileGraph) DrawDot(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// draw the call graph
	err = draw.DOT(fg.g, file, draw.GraphAttribute("rankdir", "LR"))
	if err != nil {
		return err
	}
	return nil
}

func (fg *FileGraph) FillWithYellow(vertexHash string) error {
	err := fg.setProperty(vertexHash, "style", "filled")
	if err != nil {
		return err
	}
	err = fg.setProperty(vertexHash, "color", "yellow")
	if err != nil {
		return err
	}
	return nil
}

func (fg *FileGraph) FillWithRed(vertexHash string) error {
	err := fg.setProperty(vertexHash, "style", "filled")
	if err != nil {
		return err
	}
	err = fg.setProperty(vertexHash, "color", "red")
	if err != nil {
		return err
	}
	return nil
}

func (fg *FileGraph) setProperty(vertexHash string, propertyK string, propertyV string) error {
	_, properties, err := fg.g.VertexWithProperties(vertexHash)
	if err != nil {
		return err
	}
	properties.Attributes[propertyK] = propertyV
	return nil
}

func (fg *FileGraph) ToG6Data() (*G6Data, error) {
	data := &G6Data{
		Nodes: make([]*G6Node, 0),
		Edges: make([]*G6Edge, 0),
	}

	adjacencyMap, err := fg.g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	// mapping
	mapping := make(map[string]int)
	curId := 0

	// Nodes
	for nodeId := range adjacencyMap {
		node, err := fg.g.Vertex(nodeId)
		if err != nil {
			return nil, err
		}
		mapping[node.Id()] = curId
		curNode := &G6Node{
			Id:    strconv.Itoa(curId),
			Label: node.Path,
			Style: &G6NodeStyle{},
		}
		curId++
		data.Nodes = append(data.Nodes, curNode)
	}
	// Edges
	for src, targets := range adjacencyMap {
		for target := range targets {
			srcId := mapping[src]
			targetId := mapping[target]

			curEdge := &G6Edge{
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
