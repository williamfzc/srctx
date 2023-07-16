package function

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/williamfzc/srctx/graph/common"

	"github.com/dominikbraun/graph"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	object2 "github.com/opensibyl/sibyl2/pkg/extractor/object"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser"
)

const (
	TagEntry = "entry"
)

type FuncPos struct {
	Path  string `json:"path,omitempty"`
	Lang  string `json:"lang,omitempty"`
	Start int    `json:"start,omitempty"`
	End   int    `json:"end,omitempty"`
}

func (f *FuncPos) Repr() string {
	return fmt.Sprintf("%s#%d-%d", f.Path, f.Start, f.End)
}

type Vertex struct {
	*object2.Function
	*FuncPos

	// https://github.com/williamfzc/srctx/issues/41
	Tags map[string]struct{} `json:"tags,omitempty"`
}

func (fv *Vertex) Id() string {
	return fmt.Sprintf("%v:#%d-#%d:%s", fv.Path, fv.Start, fv.End, fv.GetSignature())
}

func (fv *Vertex) PosKey() string {
	return fmt.Sprintf("%s#%d", fv.Path, fv.Start)
}

func (fv *Vertex) ListTags() []string {
	ret := make([]string, 0, len(fv.Tags))
	for each := range fv.Tags {
		ret = append(ret, each)
	}
	return ret
}

func (fv *Vertex) ContainTag(tag string) bool {
	if _, ok := fv.Tags[tag]; ok {
		return true
	}
	return false
}

func (fv *Vertex) AddTag(tag string) {
	fv.Tags[tag] = struct{}{}
}

func (fv *Vertex) RemoveTag(tag string) {
	delete(fv.Tags, tag)
}

func CreateFuncVertex(f *object2.Function, fr *extractor.FunctionFileResult) *Vertex {
	cur := &Vertex{
		Function: f,
		FuncPos: &FuncPos{
			Path: fr.Path,
			Lang: string(fr.Language),
			// sync with real lines
			Start: int(f.GetSpan().Start.Row + 1),
			End:   int(f.GetSpan().End.Row + 1),
		},
		Tags: make(map[string]struct{}),
	}
	return cur
}

type Graph struct {
	// reference graph (called graph), ref -> def
	g graph.Graph[string, *Vertex]
	// reverse reference graph (call graph), def -> ref
	rg graph.Graph[string, *Vertex]

	// k: file, v: function
	Cache map[string][]*Vertex
	// k: id, v: function
	IdCache map[string]*Vertex

	// source context ptr
	sc *object.SourceContext
}

func NewEmptyFuncGraph() *Graph {
	return &Graph{
		g:       graph.New((*Vertex).Id, graph.Directed()),
		rg:      graph.New((*Vertex).Id, graph.Directed()),
		Cache:   make(map[string][]*Vertex),
		IdCache: make(map[string]*Vertex),
		sc:      nil,
	}
}

func CreateFuncGraph(fact *FactStorage, relationship *object.SourceContext) (*Graph, error) {
	fg := NewEmptyFuncGraph()

	// add all the nodes
	for path, file := range fact.cache {
		for _, eachFunc := range file.Units {
			cur := CreateFuncVertex(eachFunc, file)
			_ = fg.g.AddVertex(cur)
			fg.Cache[path] = append(fg.Cache[path], cur)
			fg.IdCache[cur.Id()] = cur
		}
	}

	// also reverse graph
	rg, err := fg.g.Clone()
	if err != nil {
		return nil, err
	}
	fg.rg = rg

	// building edges
	log.Infof("edges building")
	for path, funcs := range fg.Cache {
		for _, eachFunc := range funcs {
			// there are multi defs happened in this line
			refs, err := relationship.RefsFromLineWithLimit(path, eachFunc.DefLine, len(eachFunc.Name))
			log.Debugf("search from %s#%d, ref: %d", path, eachFunc.DefLine, len(refs))
			if err != nil {
				// no refs
				continue
			}
			for _, eachRef := range refs {
				refFile := relationship.FileName(eachRef.FileId)
				refFile = filepath.ToSlash(filepath.Clean(refFile))

				isFuncRef := false
				symbols := fact.GetSymbolsByFileAndLine(refFile, eachRef.IndexLineNumber())
				for _, eachSymbol := range symbols {
					if eachSymbol.Unit.Content == eachFunc.Name {
						isFuncRef = true
						break
					}
				}
				if !isFuncRef {
					continue
				}

				for _, eachPossibleFunc := range fg.Cache[refFile] {
					// eachPossibleFunc 's range contains eachFunc 's ref
					// so eachPossibleFunc calls eachFunc
					if eachPossibleFunc.GetSpan().ContainLine(eachRef.IndexLineNumber()) {
						refLineNumber := eachRef.LineNumber()
						log.Debugf("%v refed in %s#%v", eachFunc.Id(), refFile, refLineNumber)

						if edge, err := fg.g.Edge(eachFunc.Id(), eachPossibleFunc.Id()); err == nil {
							storage := edge.Properties.Data.(*common.EdgeStorage)
							storage.RefLines[refLineNumber] = struct{}{}
						} else {
							_ = fg.g.AddEdge(eachFunc.Id(), eachPossibleFunc.Id(), graph.EdgeData(common.NewEdgeStorage()))
						}

						if edge, err := fg.rg.Edge(eachPossibleFunc.Id(), eachFunc.Id()); err == nil {
							storage := edge.Properties.Data.(*common.EdgeStorage)
							storage.RefLines[refLineNumber] = struct{}{}
						} else {
							_ = fg.rg.AddEdge(eachPossibleFunc.Id(), eachFunc.Id(), graph.EdgeData(common.NewEdgeStorage()))
						}
					}
				}
			}
		}
	}

	// entries tag
	entries := fg.FilterFunctions(func(funcVertex *Vertex) bool {
		return len(fg.DirectReferencedIds(funcVertex)) == 0
	})
	log.Infof("detect entries: %d", len(entries))
	for _, entry := range entries {
		entry.AddTag(TagEntry)
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

func CreateFuncGraphFromGolangDir(src string) (*Graph, error) {
	sourceContext, err := parser.FromGolangSrc(src)
	if err != nil {
		return nil, err
	}
	return srcctx2graph(src, sourceContext, core.LangGo)
}

func CreateFuncGraphFromDirWithLSIF(src string, lsifFile string, lang core.LangType) (*Graph, error) {
	sourceContext, err := parser.FromLsifFile(lsifFile, src)
	if err != nil {
		return nil, err
	}
	return srcctx2graph(src, sourceContext, lang)
}

func CreateFuncGraphFromDirWithSCIP(src string, scipFile string, lang core.LangType) (*Graph, error) {
	sourceContext, err := parser.FromScipFile(scipFile, src)
	if err != nil {
		return nil, err
	}
	return srcctx2graph(src, sourceContext, lang)
}

func srcctx2graph(src string, sourceContext *object.SourceContext, lang core.LangType) (*Graph, error) {
	log.Infof("createing fact with sibyl2")

	// change workdir because sibyl2 needs to access the files
	originWorkdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = os.Chdir(src)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Chdir(originWorkdir)
	}()

	factStorage, err := CreateFact(src, lang)
	if err != nil {
		return nil, err
	}
	log.Infof("fact ready. creating func graph ...")
	funcGraph, err := CreateFuncGraph(factStorage, sourceContext)
	if err != nil {
		return nil, err
	}
	return funcGraph, nil
}
