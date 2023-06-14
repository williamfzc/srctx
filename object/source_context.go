package object

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/parser/lsif"
)

type (
	FactKind = string
	RelKind  = string
)

const (
	EdgeTypeName = "label"

	FactFile FactKind = "file"
	FactDef  FactKind = "def"

	RelContains  RelKind = "contains"
	RelReference RelKind = "reference"
)

var (
	EdgeAttrContains  = graph.EdgeAttribute(EdgeTypeName, RelContains)
	EdgeAttrReference = graph.EdgeAttribute(EdgeTypeName, RelReference)
)

type FactVertex struct {
	DocId  int
	FileId int
	Kind   FactKind
	Range  *lsif.Range
	Extras interface{}
}

type FileExtras struct {
	Path string
}

type DefExtras struct {
	DefType   string
	RawTokens []chroma.Token
}

func (v *FactVertex) Id() int {
	return v.DocId
}

func (v *FactVertex) LineNumber() int {
	return v.IndexLineNumber() + 1
}

func (v *FactVertex) IndexLineNumber() int {
	return int(v.Range.Line)
}

func (v *FactVertex) ToRelVertex() *RelVertex {
	return &RelVertex{
		DocId:  v.DocId,
		FileId: v.FileId,
		Kind:   v.Kind,
		Range:  v.Range,
	}
}

type RelVertex struct {
	DocId  int
	FileId int
	Kind   FactKind
	Range  *lsif.Range
}

func (v *RelVertex) Id() int {
	return v.DocId
}

func (v *RelVertex) LineNumber() int {
	rangeObj := v.Range
	if rangeObj == nil {
		log.Warnf("range is nil: %v", v)
		return -1
	}
	return int(rangeObj.Line + 1)
}

func (v *RelVertex) CharNumber() int {
	return int(v.Range.Character + 1)
}

type SourceContext struct {
	FileMapping map[string]int
	FactGraph   graph.Graph[int, *FactVertex]
	RelGraph    graph.Graph[int, *RelVertex]
}

func NewSourceContext() SourceContext {
	factGraph := graph.New((*FactVertex).Id, graph.Directed())
	relGraph := graph.New((*RelVertex).Id, graph.Directed())

	return SourceContext{
		FileMapping: make(map[string]int),
		FactGraph:   factGraph,
		RelGraph:    relGraph,
	}
}
