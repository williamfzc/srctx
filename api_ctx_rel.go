package srctx

import (
	"github.com/dominikbraun/graph"
)

func (sc *SourceContext) RefsByDefId(DefId int) ([]*FactVertex, error) {
	// check
	ret := make([]*FactVertex, 0)
	_, err := sc.RelGraph.Vertex(DefId)
	if err != nil {
		// no ref info, it's ok
		return ret, nil
	}

	err = graph.DFS(sc.RelGraph, DefId, func(i int) bool {
		// exclude itself
		if DefId == i {
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
