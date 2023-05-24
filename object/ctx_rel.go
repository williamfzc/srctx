package object

import (
	"fmt"

	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

func (sc *SourceContext) RefsByFileName(fileName string) ([]*RelVertex, error) {
	// get all the reference points in this file
	fileId := sc.FileId(fileName)
	if fileId == 0 {
		return nil, fmt.Errorf("no file named: %s", fileName)
	}

	// collect all the nodes starting from this file
	startPoints := make([]*RelVertex, 0)
	err := graph.BFS(sc.FactGraph, fileId, func(i int) bool {
		// exclude itself
		if i == fileId {
			return false
		}
		if _, err := sc.FactGraph.Edge(fileId, i); err != nil {
			return true
		}

		v, err := sc.FactGraph.Vertex(i)
		if err != nil {
			log.Warnf("unknown vertex: %d", i)
			return false
		}
		startPoints = append(startPoints, v.ToRelVertex())
		return false
	})
	if err != nil {
		return nil, err
	}
	return startPoints, nil
}

func (sc *SourceContext) RefsByLine(fileName string, lineNum int) ([]*RelVertex, error) {
	allVertexes, err := sc.RefsByFileName(fileName)
	if err != nil {
		return nil, err
	}
	log.Debugf("file %s refs: %d", fileName, len(allVertexes))
	ret := make([]*RelVertex, 0)
	for _, each := range allVertexes {
		if each.LineNumber() == lineNum {
			ret = append(ret, each)
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("no ref found in %s %d", fileName, lineNum)
	}
	return ret, nil
}

func (sc *SourceContext) RefsByLineAndChar(fileName string, lineNum int, charNum int) ([]*RelVertex, error) {
	allVertexes, err := sc.RefsByLine(fileName, lineNum)
	if err != nil {
		return nil, err
	}

	ret := make([]*RelVertex, 0)
	for _, each := range allVertexes {
		if int(each.Range.Character) == charNum {
			ret = append(ret, each)
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("no ref found in %s %d", fileName, lineNum)
	}
	return ret, nil
}

func (sc *SourceContext) RefsFromDefId(defId int) ([]*FactVertex, error) {
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
			return true
		}

		vertex, err := sc.FactGraph.Vertex(i)
		if err != nil {
			return false
		}

		ret = append(ret, vertex)
		return false
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// RefsFromLine todo: need some rename ...
func (sc *SourceContext) RefsFromLine(fileName string, lineNum int, charLength int) ([]*FactVertex, error) {
	startPoints, err := sc.RefsByLine(fileName, lineNum)
	if err != nil {
		return nil, err
	}

	// search all the related points
	ret := make(map[int]*FactVertex, 0)
	for _, each := range startPoints {
		// optimize
		if int(each.Range.Length) != charLength {
			continue
		}

		curRet, err := sc.RefsFromDefId(each.Id())
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
