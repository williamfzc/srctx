package srctx

import (
	"testing"

	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMainApi(t *testing.T) {
	srcctxResult, err := FromLsifZip("./parser/testdata/dump.lsif.zip")
	assert.Nil(t, err)

	factGraph := srcctxResult.FactGraph
	relGraph := srcctxResult.RelGraph

	// test these graphs
	_ = graph.DFS(factGraph, 4, func(i int) bool {
		vertex, err := factGraph.Vertex(i)
		if err != nil {
			return true
		}
		log.Infof("def in file %d range: %v", vertex.FileId, vertex.Range)

		// any links?
		relVertex, err := relGraph.Vertex(i)
		if err != nil {
			return false
		}
		err = graph.BFS(relGraph, relVertex.GetId(), func(j int) bool {
			cur, err := factGraph.Vertex(j)
			if err != nil {
				return true
			}
			log.Infof("refered by file %d range: %v", cur.FileId, cur.Range)
			return false
		})
		if err != nil {
			return false
		}

		return false
	})
}
