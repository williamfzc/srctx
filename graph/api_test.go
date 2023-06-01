package graph

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestFuncGraph(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(curFile))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.cache)

	t.Run("GetFunctionsByFile", func(t *testing.T) {
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
	})

	t.Run("DrawDot", func(t *testing.T) {
		dotFile := "a.gv"
		defer os.Remove(dotFile)
		err = fg.DrawDot(dotFile)
		assert.Nil(t, err)
		assert.FileExists(t, dotFile)
	})

	t.Run("DrawHtml", func(t *testing.T) {
		htmlFile := "a.html"
		defer os.Remove(htmlFile)
		err = fg.DrawG6Html(htmlFile)
		assert.Nil(t, err)
	})

	t.Run("RemoveNode", func(t *testing.T) {
		before, err := fg.g.Order()
		assert.Nil(t, err)
		err = fg.RemoveNodeById("graph/api_test.go:#13-#58:graph||TestFuncGraph|*testing.T|")
		assert.Nil(t, err)
		after, err := fg.g.Order()
		assert.Nil(t, err)
		assert.Equal(t, before, after+1)
	})
}
