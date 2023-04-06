package object

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/dominikbraun/graph"
	"github.com/williamfzc/srctx/parser/lsif"
)

type FactKind = string
type RelKind = string

const (
	EdgeTypeName = "label"

	FactFile FactKind = "file"
	FactDef  FactKind = "def"
	FactRef  FactKind = "ref"

	RelContains  RelKind = "contains"
	RelReference RelKind = "reference"
)

var EdgeAttrContains = graph.EdgeAttribute(EdgeTypeName, RelContains)
var EdgeAttrReference = graph.EdgeAttribute(EdgeTypeName, RelReference)

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
	DefType   DefType
	RawTokens []chroma.Token
}

func (v *FactVertex) Id() int {
	return v.DocId
}

func (v *FactVertex) LineNumber() int {
	return int(v.Range.Line + 1)
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
	return int(v.Range.Line + 1)
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

type DefType = string

const (
	DefFunction  DefType = "function"
	DefClass     DefType = "class"
	DefNamespace DefType = "namespace"
	DefUnknown   DefType = ""
)
