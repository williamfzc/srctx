package graph

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

type FileGraph struct {
	// reference graph (called graph)
	g graph.Graph[string, *FileVertex]
	// reverse reference graph (call graph)
	rg graph.Graph[string, *FileVertex]
}

type FileVertex struct {
	Path       string
	Referenced int
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
	err := funcGraph2FileGraph(fg.g, fileGraph.g)
	if err != nil {
		return nil, err
	}
	err = funcGraph2FileGraph(fg.rg, fileGraph.rg)
	if err != nil {
		return nil, err
	}

	nodeCount, err := fileGraph.g.Order()
	if err != nil {
		return nil, err
	}
	edgeCount, err := fileGraph.g.Size()
	if err != nil {
		return nil, err
	}
	log.Infof("file graph ready. nodes: %d, edges: %d", nodeCount, edgeCount)

	return fileGraph, nil
}

func funcGraph2FileGraph(f graph.Graph[string, *FuncVertex], g graph.Graph[string, *FileVertex]) error {
	edges, err := f.Edges()
	if err != nil {
		return err
	}
	for _, eachEdge := range edges {
		source, err := f.Vertex(eachEdge.Source)
		if err != nil {
			log.Warnf("vertex not found: %v", eachEdge.Source)
			continue
		}
		target, err := f.Vertex(eachEdge.Target)
		if err != nil {
			log.Warnf("vertex not found: %v", eachEdge.Target)
			continue
		}
		sourceFile := &FileVertex{Path: source.Path}
		targetFile := &FileVertex{Path: target.Path}
		if sv, err := g.Vertex(sourceFile.Id()); err == nil {
			sv.Referenced++
		} else {
			_ = g.AddVertex(sourceFile)
		}
		if tv, err := g.Vertex(targetFile.Id()); err == nil {
			tv.Referenced++
		} else {
			_ = g.AddVertex(targetFile)
		}
		_ = g.AddEdge(sourceFile.Id(), targetFile.Id())
	}

	return nil
}
