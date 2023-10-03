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
		opts := DefaultGraphOptions()
		opts.Src = src
		opts.LsifFile = filepath.Join(src, "dump.lsif")

		fg, err := CreateFileGraphFromDirWithLSIF(opts)
		assert.Nil(t, err)
		assert.NotEmpty(t, fg)
	})

	t.Run("create index", func(t *testing.T) {
		t.Skip("this case did not work in github actions")
		opts := DefaultGraphOptions()
		opts.Src = src

		fg, err := CreateFileGraphFromGolangDir(opts)
		assert.Nil(t, err)
		assert.NotEmpty(t, fg)
	})
}
