package srctx

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

func (sc *SourceContext) Files() []string {
	ret := make([]string, 0, len(sc.FileMapping))
	for k := range sc.FileMapping {
		ret = append(ret, k)
	}
	return ret
}

func (sc *SourceContext) FileId(fileName string) int {
	return sc.FileMapping[fileName]
}

func (sc *SourceContext) FileName(fileId int) string {
	for curName, curFileId := range sc.FileMapping {
		if curFileId == fileId {
			return curName
		}
	}
	return ""
}

func (sc *SourceContext) FileVertexByName(fileName string) *FactVertex {
	factVertex, err := sc.FactGraph.Vertex(sc.FileId(fileName))
	if err != nil {
		// if not found
		return nil
	}
	return factVertex
}

func (sc *SourceContext) DefVertexesByFileName(fileName string) ([]*FactVertex, error) {
	startId := sc.FileId(fileName)
	if startId == 0 {
		return nil, fmt.Errorf("no file named: %s", fileName)
	}

	ret := make([]*FactVertex, 0)
	err := graph.DFS(sc.FactGraph, startId, func(i int) bool {
		// exclude itself
		if i == startId {
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
