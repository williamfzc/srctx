package graph

import (
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

func TestStorage(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(curFile))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.cache)

	temp := "./a.json"
	err = fg.DumpJsonFile(temp)
	assert.Nil(t, err)
	assert.FileExists(t, temp)
}
