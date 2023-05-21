package graph

import (
	"github.com/dominikbraun/graph"
)

type FileGraph struct {
	// reference graph (called graph)
	g graph.Graph[string, *FileVertex]
	// reverse reference graph (call graph)
	rg graph.Graph[string, *FileVertex]
}

type FileVertex struct {
	Path string
}

func (fv *FileVertex) Id() string {
	return fv.Path
}

func (fg *FuncGraph) ToFileGraph() (*FileGraph, error) {
	// create graph
	fileGraph := &FileGraph{
		g:  graph.New((*FileVertex).Id, graph.Directed()),
		rg: graph.New((*FileVertex).Id, graph.Directed()),
	}
	// building edges
	edges, err := fg.rg.Edges()
	if err != nil {
		return nil, err
	}
	for _, eachEdge := range edges {
		source, err := fg.rg.Vertex(eachEdge.Source)
		if err != nil {
			return nil, err
		}
		target, err := fg.rg.Vertex(eachEdge.Target)
		if err != nil {
			return nil, err
		}
		sourceFile := &FileVertex{Path: source.Path}
		targetFile := &FileVertex{Path: target.Path}
		_ = fileGraph.rg.AddVertex(sourceFile)
		_ = fileGraph.rg.AddVertex(targetFile)
		_ = fileGraph.rg.AddEdge(sourceFile.Path, targetFile.Path)
	}

	edges, err = fg.g.Edges()
	if err != nil {
		return nil, err
	}
	for _, eachEdge := range edges {
		source, err := fg.g.Vertex(eachEdge.Source)
		if err != nil {
			return nil, err
		}
		target, err := fg.g.Vertex(eachEdge.Target)
		if err != nil {
			return nil, err
		}
		sourceFile := &FileVertex{Path: source.Path}
		targetFile := &FileVertex{Path: target.Path}
		_ = fileGraph.g.AddVertex(sourceFile)
		_ = fileGraph.g.AddVertex(targetFile)
		_ = fileGraph.g.AddEdge(sourceFile.Path, targetFile.Path)
	}

	return fileGraph, nil
}
