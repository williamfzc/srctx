package file_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/williamfzc/srctx/graph"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestFuncGraph_ToFileGraph(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fg, err := graph.CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.Cache)
	fileGraph, err := fg.ToFileGraph()
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
		fileGraph, err := fg.ToFileGraph()
		assert.Nil(t, err)
		before, err := fileGraph.G.Order()
		assert.Nil(t, err)
		err = fileGraph.RemoveNodeById("graph/file/api_file_test.go")
		assert.Nil(t, err)
		after, err := fileGraph.G.Order()
		assert.Nil(t, err)
		assert.Equal(t, before, after+1)
	})

	t.Run("DrawG6", func(t *testing.T) {
		fileGraph, err := fg.ToFileGraph()
		assert.Nil(t, err)
		err = fileGraph.DrawG6Html("b.html")
		assert.Nil(t, err)
	})
}
