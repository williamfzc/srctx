package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/parser"
)

func main() {
	src := flag.String("src", ".", "repo path")
	before := flag.String("before", "HEAD~1", "before rev")
	after := flag.String("after", "HEAD", "after rev")
	lsifZip := flag.String("lsif", "./dump.lsif.zip", "lsif zip path")
	outputJson := flag.String("outputJson", "./srctx.diff.json", "srctx json output")
	flag.Parse()

	// prepare
	lineMap, err := diff.GitDiff(*src, *before, *after)
	panicIfErr(err)
	sourceContext, err := parser.FromLsifZip(*lsifZip)
	panicIfErr(err)

	// calc
	lineStats := make([]*LineStat, 0)
	log.Infof("diff / total (files): %d / %d", len(lineMap), len(sourceContext.Files()))
	for path, lines := range lineMap {
		for _, eachLine := range lines {
			lineStat := NewLineStat(path, eachLine)
			vertices, _ := sourceContext.RefsByLine(path, eachLine)
			log.Debugf("path %s line %d affected %d vertexes", path, eachLine, len(vertices))

			lineStat.RefScope.TotalRefCount = len(vertices)
			for _, eachVertex := range vertices {
				refFileName := sourceContext.FileName(eachVertex.FileId)
				if refFileName != path {
					lineStat.RefScope.CrossFileRefCount++
				}
				if filepath.Dir(refFileName) != filepath.Dir(path) {
					lineStat.RefScope.CrossDirRefCount++
				}
			}
			lineStats = append(lineStats, lineStat)
		}
	}
	log.Infof("diff finished.")

	data, err := json.Marshal(lineStats)
	panicIfErr(err)
	err = os.WriteFile(*outputJson, data, 0644)
	panicIfErr(err)
	log.Infof("dump finished.")
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
