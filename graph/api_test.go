package graph

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	fg, err := CreateFuncGraphFromDirWithLSIF("../", "../dump.lsif")
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.cache)

	testFuncs := fg.GetFunctionsByFile("graph/api_test.go")
	assert.NotEmpty(t, testFuncs)

	for _, eachFunc := range testFuncs {
		if eachFunc.Name == "TestApi" {
			beingRefs := fg.ReferencedIds(eachFunc)
			refOut := fg.ReferenceIds(eachFunc)
			assert.Len(t, beingRefs, 0)
			assert.Len(t, refOut, 4)
		}
	}
}

func TestFuncGraph_DrawDot(t *testing.T) {
	fg, err := CreateFuncGraphFromDirWithLSIF("../", "../dump.lsif")
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.cache)

	dotFile := "a.gv"
	defer os.Remove(dotFile)
	err = fg.DrawDot(dotFile)
	assert.Nil(t, err)
	assert.FileExists(t, dotFile)
}
