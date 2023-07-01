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
	fg, err := CreateFileGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"))
	assert.Nil(t, err)
	assert.NotEmpty(t, fg)
}
