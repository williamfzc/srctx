package graph

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/graph/file"
)

func (fg *FuncGraph) ToFileGraph() (*file.FileGraph, error) {
	// create graph
	fileGraph := &file.FileGraph{
		G:  graph.New((*file.FileVertex).Id, graph.Directed()),
		Rg: graph.New((*file.FileVertex).Id, graph.Directed()),
	}
	// building edges
	err := FuncGraph2FileGraph(fg.g, fileGraph.G)
	if err != nil {
		return nil, err
	}
	err = FuncGraph2FileGraph(fg.rg, fileGraph.Rg)
	if err != nil {
		return nil, err
	}

	nodeCount, err := fileGraph.G.Order()
	if err != nil {
		return nil, err
	}
	edgeCount, err := fileGraph.G.Size()
	if err != nil {
		return nil, err
	}
	log.Infof("file graph ready. nodes: %d, edges: %d", nodeCount, edgeCount)

	return fileGraph, nil
}

func (fg *FuncGraph) ToDirGraph() (*file.FileGraph, error) {
	fileGraph, err := fg.ToFileGraph()
	if err != nil {
		return nil, err
	}
	return fileGraph.ToDirGraph()
}

func FuncGraph2FileGraph(f graph.Graph[string, *FuncVertex], g graph.Graph[string, *file.FileVertex]) error {
	m, err := f.AdjacencyMap()
	if err != nil {
		return err
	}
	// add all the vertices
	for k := range m {
		v, err := f.Vertex(k)
		if err != nil {
			return err
		}
		_ = g.AddVertex(file.Path2vertex(v.Path))
	}

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

		// ignore self ptr
		if source.Path == target.Path {
			continue
		}

		sourceFile, err := g.Vertex(source.Path)
		if err != nil {
			return err
		}
		targetFile, err := g.Vertex(target.Path)
		if err != nil {
			return err
		}
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
