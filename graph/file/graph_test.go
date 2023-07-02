package file

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))

	t.Run("from lsif", func(t *testing.T) {
		fg, err := CreateFileGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"))
		assert.Nil(t, err)
		assert.NotEmpty(t, fg)
	})

	t.Run("create index", func(t *testing.T) {
		fg, err := CreateFileGraphFromGolangDir(src)
		assert.Nil(t, err)
		assert.NotEmpty(t, fg)
	})
}
