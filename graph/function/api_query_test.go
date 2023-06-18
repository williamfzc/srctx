package function

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.Cache)

	t.Run("GetFunctionsByFile", func(t *testing.T) {
		testFuncs := fg.GetFunctionsByFile("graph/function/api_query_test.go")
		assert.NotEmpty(t, testFuncs)

		for _, eachFunc := range testFuncs {
			if eachFunc.Name == "TestQuery" {
				beingRefs := fg.DirectReferencedIds(eachFunc)
				refOut := fg.DirectReferenceIds(eachFunc)
				assert.Len(t, beingRefs, 0)
				assert.Len(t, refOut, 4)
			}
		}
	})
}
