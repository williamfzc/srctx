package graph

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dominikbraun/graph/draw"
	"github.com/goccy/go-json"
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

func (fg *FileGraph) DrawG6Html(filename string) error {
	data := &g6data{
		Nodes: make([]*g6node, 0),
		Edges: make([]*g6edge, 0),
	}

	adjacencyMap, err := fg.g.AdjacencyMap()
	if err != nil {
		return err
	}
	// mapping
	mapping := make(map[string]int)
	curId := 0

	// Nodes
	for nodeId := range adjacencyMap {
		node, err := fg.g.Vertex(nodeId)
		if err != nil {
			return err
		}
		mapping[node.Id()] = curId
		curNode := &g6node{
			Id:    strconv.Itoa(curId),
			Label: node.Path,
		}
		curId++
		data.Nodes = append(data.Nodes, curNode)
	}
	// Edges
	for src, targets := range adjacencyMap {
		for target := range targets {
			srcId := mapping[src]
			targetId := mapping[target]

			curEdge := &g6edge{
				Source: strconv.Itoa(srcId),
				Target: strconv.Itoa(targetId),
			}
			data.Edges = append(data.Edges, curEdge)
		}
	}
	// render
	dataRaw, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	htmlContent := fmt.Sprintf(g6template, dataRaw)
	err = os.WriteFile(filename, []byte(htmlContent), 0o666)
	if err != nil {
		return err
	}

	return nil
}
