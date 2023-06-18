package function

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestModify(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.Cache)

	t.Run("RemoveNode", func(t *testing.T) {
		before, err := fg.g.Order()
		assert.Nil(t, err)
		err = fg.RemoveNodeById("graph/function/api_modify_test.go:#12-#28:function||TestModify|*testing.T|")
		assert.Nil(t, err)
		after, err := fg.g.Order()
		assert.Nil(t, err)
		assert.Equal(t, before, after+1)
	})
}
