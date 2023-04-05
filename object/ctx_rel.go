package object

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

func (sc *SourceContext) RefsByDefId(defId int) ([]*FactVertex, error) {
	// check
	ret := make([]*FactVertex, 0)
	_, err := sc.RelGraph.Vertex(defId)
	if err != nil {
		// no ref info, it's ok
		return ret, nil
	}

	err = graph.BFS(sc.RelGraph, defId, func(i int) bool {
		// exclude itself
		if defId == i {
			return false
		}
		// connected to current?
		if _, err := sc.RelGraph.Edge(defId, i); err != nil {
			return false
		}

		vertex, err := sc.FactGraph.Vertex(i)
		if err != nil {
			return true
		}

		ret = append(ret, vertex)
		return false
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (sc *SourceContext) RefsByLine(fileName string, lineNum int) ([]*FactVertex, error) {
	allVertexes, err := sc.DefVertexesByFileName(fileName)
	if err != nil {
		return nil, err
	}
	var startPoint *FactVertex
	for _, each := range allVertexes {
		if each.LineNumber() == lineNum {
			startPoint = each
			break
		}
	}
	if startPoint == nil {
		return nil, fmt.Errorf("no def found in %s %d", fileName, lineNum)
	}
	ret, err := sc.RefsByDefId(startPoint.Id())
	if err != nil {
		return nil, err
	}
	return ret, nil
}
