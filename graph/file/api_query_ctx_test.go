package file

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/williamfzc/srctx/graph/common"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))

	opts := common.DefaultGraphOptions()
	opts.Src = src
	opts.LsifFile = filepath.Join(src, "dump.lsif")
	fg, err := CreateFileGraphFromDirWithLSIF(opts)
	assert.Nil(t, err)
	assert.NotEmpty(t, fg)

	t.Run("List", func(t *testing.T) {
		fs := fg.ListFiles()
		assert.NotEmpty(t, fs)
	})

	t.Run("GetHot", func(t *testing.T) {
		fileVertex := fg.GetById("graph/function/api_query_test.go")
		assert.NotEmpty(t, fileVertex)

		// test file should not be referenced
		shouldEmpty := fg.DirectReferencedIds(fileVertex)
		assert.Empty(t, shouldEmpty)

		shouldNotEmpty := fg.DirectReferenceIds(fileVertex)
		assert.NotEmpty(t, shouldNotEmpty)
		log.Infof("outv: %d", len(shouldNotEmpty))
	})

	t.Run("Entries", func(t *testing.T) {
		entries := fg.EntryIds(fg.GetById("graph/function/api_query_test.go"))
		// test case usually is an entry point
		assert.Len(t, entries, 1)
	})
}
