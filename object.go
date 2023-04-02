package srctx

import (
	"github.com/dominikbraun/graph"
	"github.com/williamfzc/srctx/parser"
)

type FactKind = string
type RelKind = string

const (
	EdgeTypeName = "relType"

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
	Range  *parser.Range
	Extras interface{}
}

type FileExtras struct {
	Path string
}

func (v *FactVertex) GetId() int {
	return v.DocId
}

type RelVertex struct {
	DocId int
	Kind  RelKind
}

func (v *RelVertex) GetId() int {
	return v.DocId
}

type SourceContext struct {
	FactGraph graph.Graph[int, *FactVertex]
	RelGraph  graph.Graph[int, *RelVertex]
}
