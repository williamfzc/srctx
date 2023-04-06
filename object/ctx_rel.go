package object

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

func (sc *SourceContext) RefsByFileName(fileName string) ([]*RelVertex, error) {
	// get all the reference points in this file
	startId := sc.FileId(fileName)
	if startId == 0 {
		return nil, fmt.Errorf("no file named: %s", fileName)
	}

	ret := make([]*RelVertex, 0)
	err := graph.DFS(sc.RelGraph, startId, func(i int) bool {
		// exclude itself
		if i == startId {
			return false
		}
		vertex, err := sc.RelGraph.Vertex(i)
		if err != nil {
			return true
		}

		if vertex.Kind == FactRef {
			if _, err := sc.RelGraph.Edge(startId, i); err == nil {
				ret = append(ret, vertex)
			}
		}

		return false
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

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
	allVertexes, err := sc.RefsByFileName(fileName)
	if err != nil {
		return nil, err
	}
	startPoints := make([]*RelVertex, 0)
	for _, each := range allVertexes {
		if each.LineNumber() == lineNum {
			startPoints = append(startPoints, each)
		}
	}
	if len(startPoints) == 0 {
		return nil, fmt.Errorf("no ref found in %s %d", fileName, lineNum)
	}

	// search all the related points
	ret := make(map[int]*FactVertex, 0)
	for _, each := range startPoints {
		curRet, err := sc.RefsByDefId(each.Id())
		if err != nil {
			return nil, err
		}
		for _, eachRef := range curRet {
			ret[eachRef.Id()] = eachRef
		}
	}

	final := make([]*FactVertex, 0, len(ret))
	for _, v := range ret {
		final = append(final, v)
	}
	return final, nil
}
