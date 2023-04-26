package collector

import (
	"fmt"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	object2 "github.com/opensibyl/sibyl2/pkg/extractor/object"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/object"
)

type FuncPos struct {
	Path  string
	Lang  string
	Start int
	End   int
}

type FuncVertex struct {
	*object2.Function
	*FuncPos
}

func (fv *FuncVertex) DefLine() int {
	// not always correct
	// todo: need sibyl2 improvement
	return fv.Start
}

func (fv *FuncVertex) Id() string {
	return fmt.Sprintf("%v:#%d-#%d:%s", fv.Path, fv.Start, fv.End, fv.GetSignature())
}

func (fv *FuncVertex) PosKey() string {
	return fmt.Sprintf("%s#%d", fv.Path, fv.Start)
}

func CreateFuncVertex(f *object2.Function, fr *extractor.FunctionFileResult) *FuncVertex {
	cur := &FuncVertex{
		Function: f,
		FuncPos: &FuncPos{
			Path: fr.Path,
			Lang: string(fr.Language),
			// sync with real lines
			Start: int(f.GetSpan().Start.Row + 1),
			End:   int(f.GetSpan().End.Row + 1),
		},
	}
	return cur
}

type FuncGraph struct {
	g     graph.Graph[string, *FuncVertex]
	cache map[string][]*FuncVertex
}

func (fg *FuncGraph) InfluenceCount(f *FuncVertex) int {
	ret := 0
	startPoint := f.Id()
	_ = graph.BFS(fg.g, startPoint, func(s string) bool {
		if startPoint == s {
			return false
		}
		if _, err := fg.g.Edge(startPoint, s); err != nil {
			return true
		}

		log.Infof("direct ref: %s", s)
		ret++
		return false
	})
	return ret
}

func CreateGraph(fact *FactStorage, relationship *object.SourceContext) (*FuncGraph, error) {
	fg := &FuncGraph{
		g:     graph.New((*FuncVertex).Id, graph.Directed()),
		cache: make(map[string][]*FuncVertex),
	}

	// add all the nodes
	for path, file := range fact.cache {
		for _, eachFunc := range file.Units {
			cur := CreateFuncVertex(eachFunc, file)
			_ = fg.g.AddVertex(cur)
			fg.cache[path] = append(fg.cache[path], cur)
		}
	}

	// building edges
	for path, funcs := range fg.cache {
		for _, eachFunc := range funcs {
			refs, err := relationship.RefsByLine(path, eachFunc.DefLine())
			log.Infof("search from %s#%d, ref: %d", path, eachFunc.DefLine(), len(refs))
			if err != nil {
				// no refs
				continue
			}
			for _, eachRef := range refs {
				refFile := relationship.FileName(eachRef.FileId)
				for _, eachPossibleFunc := range fg.cache[refFile] {
					if eachPossibleFunc.GetSpan().ContainLine(eachRef.LineNumber()) {
						// build `referenced by` edge
						// double check from file
						lineContent := fileCache.GetLine(refFile, eachRef.LineNumber())
						if !strings.Contains(lineContent, eachFunc.Name) {
							log.Infof("%s not refed in %s", lineContent, eachFunc.Name)
							continue
						}
						log.Infof("%v refed in %s#%v", eachFunc.Id(), refFile, eachRef.LineNumber())
						_ = fg.g.AddEdge(eachFunc.Id(), eachPossibleFunc.Id())
					}
				}
			}
		}
	}

	nodeCount, err := fg.g.Order()
	if err != nil {
		return nil, err
	}
	edgeCount, err := fg.g.Size()
	if err != nil {
		return nil, err
	}
	log.Infof("func graph ready. nodes: %d, edges: %d", nodeCount, edgeCount)

	return fg, nil
}
