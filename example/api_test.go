package example

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/graph/function"
)

func TestFunc(t *testing.T) {
	_, curFile, _, _ := runtime.Caller(0)
	src := filepath.Dir(filepath.Dir(curFile))
	lsif := "../dump.lsif"
	lang := core.LangGo

	funcGraph, err := function.CreateFuncGraphFromDirWithLSIF(src, lsif, lang)
	if err != nil {
		panic(err)
	}

	t.Run("file", func(t *testing.T) {
		files := funcGraph.ListFiles()
		for _, each := range files {
			log.Debugf("file: %v", each)
		}
	})

	t.Run("func", func(t *testing.T) {
		functions := funcGraph.GetFunctionsByFile("graph/function/api_query_func.go")
		for _, each := range functions {
			// about this function
			log.Infof("func: %v", each.Id())
			log.Infof("decl location: %v", each.FuncPos.Repr())
			log.Infof("func name: %v", each.Name)

			// context of this function
			outVs := funcGraph.DirectReferencedIds(each)
			log.Infof("this function referenced by %v other functions", len(outVs))
			for _, eachOutV := range outVs {
				outV, _ := funcGraph.GetById(eachOutV)
				log.Infof("%v directly referenced by %v", each.Name, outV.Name)
			}
			transOutVs := funcGraph.TransitiveReferencedIds(each)
			log.Infof("this function transitively referenced by %d other functions", len(transOutVs))

			allEntries := funcGraph.ListEntries()
			entries := funcGraph.EntryIds(each)
			log.Infof("this function affects %d/%d entries", len(entries), len(allEntries))
		}
	})
}
