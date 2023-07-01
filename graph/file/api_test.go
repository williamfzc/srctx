package file_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/williamfzc/srctx/graph/file"

	"github.com/stretchr/testify/assert"
)

func TestFileGraph(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fileGraph, err := file.CreateFileGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"))
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
}
