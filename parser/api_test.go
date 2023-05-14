package parser

import (
	"path/filepath"
	"testing"

	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	srcctxResult, err := FromLsifFile("./lsif/testdata/dump.lsif.zip", ".")
	assert.Nil(t, err)

	factGraph := srcctxResult.FactGraph
	relGraph := srcctxResult.RelGraph

	t.Run("test_internal", func(t *testing.T) {
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
			err = graph.BFS(relGraph, relVertex.Id(), func(j int) bool {
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
	})

	t.Run("test_ctx", func(t *testing.T) {
		fileName := "morestrings/reverse.go"
		allDefVertexes, err := srcctxResult.DefsByFileName(fileName)
		assert.Nil(t, err)
		for _, each := range allDefVertexes {
			vertices, err := srcctxResult.RefsByDefId(each.Id())
			assert.Nil(t, err)
			for _, eachV := range vertices {
				log.Infof("def in file %s %d:%d, ref in: %s %d:%d",
					fileName, each.LineNumber(),
					each.Range.Character+1,
					srcctxResult.FileName(eachV.FileId),
					eachV.LineNumber(),
					eachV.Range.Character+1)
			}
		}
	})
}

func TestGen(t *testing.T) {
	root, err := filepath.Abs("../")
	assert.Nil(t, err)
	_, err = FromGolangSrc(root)
	assert.Nil(t, err)
}
