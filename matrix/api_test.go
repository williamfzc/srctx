package matrix

import (
	"path/filepath"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/williamfzc/srctx/graph/common"
	"github.com/williamfzc/srctx/graph/file"
)

func TestMatrix(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(curFile))
	opts := common.DefaultGraphOptions()
	opts.Src = src
	opts.LsifFile = filepath.Join(src, "dump.lsif")
	fileGraph, err := file.CreateFileGraphFromDirWithLSIF(opts)
	assert.Nil(t, err)

	t.Run("test_a", func(t *testing.T) {
		matrixFromGraph, err := CreateMatrixFromGraph(fileGraph.G)
		assert.Nil(t, err)
		assert.NotEmpty(t, matrixFromGraph)

		targetFile := "graph/common/edge.go"
		matrixFromGraph.ForEach(targetFile, func(i int, v float64) {
			log.Infof("%s value: %v", matrixFromGraph.ById(i), v)
		})
	})
}
