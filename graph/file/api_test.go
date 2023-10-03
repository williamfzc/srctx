package file_test

import (
	"path/filepath"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/williamfzc/srctx/graph/file"
)

func TestFileGraph(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	opts := file.DefaultGraphOptions()
	opts.Src = src
	opts.LsifFile = filepath.Join(src, "dump.lsif")
	fileGraph, err := file.CreateFileGraphFromDirWithLSIF(opts)
	assert.Nil(t, err)

	t.Run("Transform", func(t *testing.T) {
		size, err := fileGraph.G.Size()
		assert.Nil(t, err)
		assert.NotEqual(t, size, 0)

		// dir level
		dirGraph, err := fileGraph.ToDirGraph()
		assert.Nil(t, err)
		assert.NotEqual(t, dirGraph, 0)
	})

	t.Run("RemoveNode", func(t *testing.T) {
		before, err := fileGraph.G.Order()
		assert.Nil(t, err)
		err = fileGraph.RemoveNodeById("graph/file/api_test.go")
		assert.Nil(t, err)
		after, err := fileGraph.G.Order()
		assert.Nil(t, err)
		assert.Equal(t, before, after+1)
	})

	t.Run("DrawG6", func(t *testing.T) {
		err = fileGraph.DrawG6Html("b.html")
		assert.Nil(t, err)
	})

	t.Run("Relation", func(t *testing.T) {
		edgeStorage, err := fileGraph.RelationBetween("graph/function/api_query.go", "graph/function/api_query_test.go")
		assert.Nil(t, err)
		log.Debugf("ref lines: %v", edgeStorage.RefLines)
		assert.NotEmpty(t, edgeStorage.RefLines)
	})

	t.Run("stat", func(t *testing.T) {
		ptr := fileGraph.GetById("graph/function/api_query.go")
		stat := fileGraph.GlobalStat([]*file.Vertex{ptr})
		assert.NotEmpty(t, stat)
		assert.NotEmpty(t, stat.ImpactUnitsMap)
	})
}
