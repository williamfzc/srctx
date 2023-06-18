package function

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dominikbraun/graph"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	object2 "github.com/opensibyl/sibyl2/pkg/extractor/object"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser"
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

type FuncVertex struct {
	*object2.Function
	*FuncPos

	// https://github.com/williamfzc/srctx/issues/41
	Tags map[string]struct{} `json:"tags,omitempty"`
}

func (fv *FuncVertex) Id() string {
	return fmt.Sprintf("%v:#%d-#%d:%s", fv.Path, fv.Start, fv.End, fv.GetSignature())
}

func (fv *FuncVertex) PosKey() string {
	return fmt.Sprintf("%s#%d", fv.Path, fv.Start)
}

func (fv *FuncVertex) ListTags() []string {
	ret := make([]string, 0, len(fv.Tags))
	for each := range fv.Tags {
		ret = append(ret, each)
	}
	return ret
}

func (fv *FuncVertex) ContainTag(tag string) bool {
	if _, ok := fv.Tags[tag]; ok {
		return true
	}
	return false
}

func (fv *FuncVertex) AddTag(tag string) {
	fv.Tags[tag] = struct{}{}
}

func (fv *FuncVertex) RemoveTag(tag string) {
	delete(fv.Tags, tag)
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
		Tags: make(map[string]struct{}),
	}
	return cur
}

type FuncGraph struct {
	// reference graph (called graph), ref -> def
	g graph.Graph[string, *FuncVertex]
	// reverse reference graph (call graph), def -> ref
	rg graph.Graph[string, *FuncVertex]

	// k: file, v: function
	Cache map[string][]*FuncVertex
	// k: id, v: function
	IdCache map[string]*FuncVertex

	// source context ptr
	sc *object.SourceContext
}

func NewEmptyFuncGraph() *FuncGraph {
	return &FuncGraph{
		g:       graph.New((*FuncVertex).Id, graph.Directed()),
		rg:      graph.New((*FuncVertex).Id, graph.Directed()),
		Cache:   make(map[string][]*FuncVertex),
		IdCache: make(map[string]*FuncVertex),
		sc:      nil,
	}
}

func CreateFuncGraph(fact *FactStorage, relationship *object.SourceContext) (*FuncGraph, error) {
	fg := NewEmptyFuncGraph()

	// add all the nodes
	for path, file := range fact.cache {
		for _, eachFunc := range file.Units {
			cur := CreateFuncVertex(eachFunc, file)
			_ = fg.g.AddVertex(cur)
			fg.Cache[path] = append(fg.Cache[path], cur)
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
					if eachPossibleFunc.GetSpan().ContainLine(eachRef.IndexLineNumber()) {
						// build `referenced by` edge
						log.Debugf("%v refed in %s#%v", eachFunc.Id(), refFile, eachRef.LineNumber())
						// eachFunc def, eachPossibleFunc ref
						_ = fg.rg.AddEdge(eachFunc.Id(), eachPossibleFunc.Id())
						_ = fg.g.AddEdge(eachPossibleFunc.Id(), eachFunc.Id())
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

func CreateFuncGraphFromGolangDir(src string) (*FuncGraph, error) {
	sourceContext, err := parser.FromGolangSrc(src)
	if err != nil {
		return nil, err
	}
	return srcctx2graph(src, sourceContext, core.LangGo)
}

func CreateFuncGraphFromDirWithLSIF(src string, lsifFile string, lang core.LangType) (*FuncGraph, error) {
	sourceContext, err := parser.FromLsifFile(lsifFile, src)
	if err != nil {
		return nil, err
	}
	return srcctx2graph(src, sourceContext, lang)
}

func CreateFuncGraphFromDirWithSCIP(src string, scipFile string, lang core.LangType) (*FuncGraph, error) {
	sourceContext, err := parser.FromScipFile(scipFile, src)
	if err != nil {
		return nil, err
	}
	return srcctx2graph(src, sourceContext, lang)
}

func srcctx2graph(src string, sourceContext *object.SourceContext, lang core.LangType) (*FuncGraph, error) {
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
