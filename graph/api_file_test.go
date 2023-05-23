package graph

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncGraph_ToFileGraph(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(curFile))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"))
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.cache)

	fileGraph, err := fg.ToFileGraph()
	assert.Nil(t, err)
	size, err := fileGraph.g.Size()
	assert.Nil(t, err)
	assert.NotEqual(t, size, 0)
}
