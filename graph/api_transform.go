package graph

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

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

func (fg *FuncGraph) ToDirGraph() (*FileGraph, error) {
	fileGraph, err := fg.ToFileGraph()
	if err != nil {
		return nil, err
	}
	return fileGraph.ToDirGraph()
}
