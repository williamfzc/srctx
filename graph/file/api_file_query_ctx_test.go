package file

import (
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(filepath.Dir(curFile)))
	fg, err := CreateFileGraphFromDirWithLSIF(src, filepath.Join(src, "dump.lsif"))
	assert.Nil(t, err)
	assert.NotEmpty(t, fg)

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
}
