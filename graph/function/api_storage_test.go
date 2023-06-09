package function

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fg, err := CreateFuncGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"), core.LangGo)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg.Cache)

	temp := "./temp.msgpack"
	err = fg.DumpFile(temp)
	assert.Nil(t, err)
	assert.FileExists(t, temp)

	newFg, err := LoadFile(temp)
	assert.Nil(t, err)

	oldOrd, _ := fg.g.Order()
	oldSize, _ := fg.g.Size()
	newOrd, _ := newFg.g.Order()
	newSize, _ := newFg.g.Size()
	assert.Equal(t, oldOrd, newOrd)
	assert.Equal(t, oldSize, newSize)
}
