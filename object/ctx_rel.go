package object

import (
	"fmt"

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
	for each := range sc.FactAdjMap[fileId] {
		factVertex, err := sc.FactGraph.Vertex(each)
		if err != nil {
			return nil, err
		}
		startPoints = append(startPoints, factVertex.ToRelVertex())
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

	for each := range sc.RelAdjMap[defId] {
		vertex, err := sc.FactGraph.Vertex(each)
		if err != nil {
			return nil, err
		}
		ret = append(ret, vertex)
	}

	return ret, nil
}

func (sc *SourceContext) RefsFromLineWithLimit(fileName string, lineNum int, charLength int) ([]*FactVertex, error) {
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
