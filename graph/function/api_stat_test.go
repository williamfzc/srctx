package function

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestStat(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.Cache)

	t.Run("stat", func(t *testing.T) {
		ptr, err := fg.GetById("graph/function/graph.go:#216-#222:function||CreateFuncGraphFromDirWithLSIF|string,string,core.LangType|*Graph,error")
		assert.Nil(t, err)
		stat := fg.GlobalStat([]*Vertex{ptr})
		assert.NotEmpty(t, stat)
		assert.NotEmpty(t, stat.ImpactUnitsMap)
	})
}
