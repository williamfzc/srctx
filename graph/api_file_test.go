package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncGraph_ToFileGraph(t *testing.T) {
	fg, err := CreateFuncGraphFromDirWithLSIF("../", "../dump.lsif")
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.cache)

	fileGraph, err := fg.ToFileGraph()
	assert.Nil(t, err)
	size, err := fileGraph.g.Size()
	assert.Nil(t, err)
	assert.NotEqual(t, size, 0)
}
