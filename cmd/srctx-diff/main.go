package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/parser"
)

func main() {
	src := flag.String("src", ".", "repo path")
	before := flag.String("before", "HEAD~1", "before rev")
	after := flag.String("after", "HEAD", "after rev")
	lsifZip := flag.String("lsif", "./dump.lsif.zip", "lsif zip path")
	flag.Parse()

	// prepare
	lineMap, err := diff.GitDiff(*src, *before, *after)
	panicIfErr(err)
	sourceContext, err := parser.FromLsifZip(*lsifZip)
	panicIfErr(err)

	// calc
	log.Infof("diff / total (files): %d / %d", len(lineMap), len(sourceContext.Files()))
	for path, lines := range lineMap {
		for _, eachLine := range lines {
			vertices, _ := sourceContext.RefsByLine(path, eachLine)
			log.Infof("path %s line %d affected %d vertexes", path, eachLine, len(vertices))
		}
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
