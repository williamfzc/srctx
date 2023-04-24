package collector

import (
	"fmt"

	"github.com/dominikbraun/graph"
	object2 "github.com/opensibyl/sibyl2/pkg/extractor/object"
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
	return fv.Start
}

func (fv *FuncVertex) Id() string {
	return fmt.Sprintf("%v:#%d-#%d:%s", fv.Path, fv.Start, fv.End, fv.GetSignature())
}

func (fv *FuncVertex) PosKey() string {
	return fmt.Sprintf("%s#%d", fv.Path, fv.Start)
}

type FuncGraph struct {
	g     graph.Graph[string, *FuncVertex]
	cache map[string][]*FuncVertex
}

func (fg *FuncGraph) InfluenceCount(f *FuncVertex) int {
	ret := 0
	_ = graph.BFS(fg.g, f.Id(), func(s string) bool {
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
			cur := &FuncVertex{
				Function: eachFunc,
				FuncPos: &FuncPos{
					Path:  file.Path,
					Lang:  string(file.Language),
					Start: int(eachFunc.GetSpan().Start.Row),
					End:   int(eachFunc.GetSpan().End.Row),
				},
			}
			_ = fg.g.AddVertex(cur)
			fg.cache[path] = append(fg.cache[path], cur)
		}
	}

	// building edges
	for _, funcs := range fg.cache {
		for _, eachFunc := range funcs {
			refs, _ := relationship.RefsByLine(eachFunc.Path, eachFunc.DefLine())
			for _, eachRef := range refs {
				refFile := relationship.FileName(eachRef.FileId)
				for _, eachPossibleFunc := range fg.cache[refFile] {
					if eachPossibleFunc.GetSpan().ContainLine(eachRef.LineNumber()) {
						// build edge
						_ = fg.g.AddEdge(eachFunc.Id(), eachPossibleFunc.Id())
					}
				}
			}
		}
	}

	return fg, nil
}
