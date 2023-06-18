package function

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestFuncGraph_Draw(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.Cache)

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
}
