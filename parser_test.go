package srctx

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/dominikbraun/graph"
	"github.com/stretchr/testify/assert"
	"github.com/williamfzc/srctx/parser"
)

type FactKind = string
type RelKind = string

const (
	EdgeTypeName = "relType"

	FactFile     FactKind = "file"
	FactDef      FactKind = "def"
	FactRef      FactKind = "ref"
	RelContains  RelKind  = "contains"
	RelReference RelKind  = "reference"
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

func TestParser(t *testing.T) {
	file, err := os.Open("./dump.lsif.zip")
	assert.Nil(t, err)
	defer file.Close()
	newParser, err := parser.NewParser(context.Background(), file)
	assert.Nil(t, err)
	assert.NotEmpty(t, newParser.Docs)

	factGraph := graph.New((*FactVertex).GetId, graph.Directed())
	relGraph := graph.New((*RelVertex).GetId, graph.Directed())

	// file level
	for eachFileId, eachFile := range newParser.Docs.Entries {
		eachFileVertex := &FactVertex{
			DocId:  int(eachFileId),
			FileId: int(eachFileId),
			Range:  nil,
			Kind:   FactFile,
			Extras: &FileExtras{
				Path: eachFile,
			},
		}
		err := factGraph.AddVertex(eachFileVertex)
		assert.Nil(t, err)
	}

	// contains / defs
	for eachFileId, eachFile := range newParser.Docs.Entries {
		// def ranges in this file
		if inVs, ok := newParser.Docs.DocRanges[eachFileId]; ok {
			for _, eachRangeId := range inVs {
				rawRange := &parser.Range{}
				err := newParser.Docs.Ranges.Cache.Entry(eachRangeId, rawRange)
				assert.Nil(t, err)

				eachRangeVertex := &FactVertex{
					DocId:  int(eachRangeId),
					FileId: int(eachFileId),
					Kind:   FactDef,
					Range:  rawRange,
					Extras: nil,
				}
				log.Printf("f %s range %v", eachFile, rawRange)

				err = factGraph.AddVertex(eachRangeVertex)
				assert.Nil(t, err)

				// and edge
				err = factGraph.AddEdge(int(eachFileId), eachRangeVertex.GetId(), EdgeAttrContains)
				assert.Nil(t, err)
			}
		} else {
			log.Printf("no any links to %d", eachFileId)
		}
	}

	// the fact graph is ready
	// then real graph

	// refs
	reverseNextMap := reverseMap(newParser.Docs.Ranges.NextMap)
	reverseRefMap := reverseMap(newParser.Docs.Ranges.TextReferenceMap)
	for eachReferenceResultId, eachDef := range newParser.Docs.Ranges.DefRefs {
		refFileId := eachDef.DocId
		log.Printf("def %d in file %s line %d",
			eachReferenceResultId,
			newParser.Docs.Entries[refFileId],
			eachDef.Line)

		refs := newParser.Docs.Ranges.References.GetItems(eachReferenceResultId)
		refRanges := make(map[parser.Id]interface{}, 0)
		for _, eachRef := range refs {
			assert.Nil(t, err)
			log.Printf("def %d refed in file %s, line %d",
				eachReferenceResultId,
				newParser.Docs.Entries[eachRef.DocId],
				eachRef.Line)

			refRanges[eachRef.RangeId] = nil
		}

		// connect between ranges to ranges
		// range - next -> resultSet - text/references -> referenceResult - item -> range
		for eachRefRange := range refRanges {
			// starts with the ref point
			resultSetId, ok := newParser.Docs.Ranges.NextMap[eachRefRange]
			if !ok {
				log.Printf("failed to jump with nextMap: %v", eachRefRange)
				continue
			}
			foundReferenceResultId, ok := newParser.Docs.Ranges.TextReferenceMap[resultSetId]
			if !ok {
				log.Printf("failed to jump with reference map: %v", resultSetId)
				continue
			}
			assert.Equal(t, eachReferenceResultId, foundReferenceResultId)

			foundItem, ok := reverseRefMap[foundReferenceResultId]
			if !ok {
				log.Printf("failed to jump with rev ref map: %v", resultSetId)
				continue
			}

			foundRange, ok := reverseNextMap[foundItem]
			if !ok {
				log.Printf("failed to jump with rev next map: %v", resultSetId)
				continue
			}

			_ = relGraph.AddVertex(&RelVertex{
				DocId: int(eachRefRange),
				Kind:  FactRef,
			})
			_ = relGraph.AddVertex(&RelVertex{
				DocId: int(foundRange),
				Kind:  FactDef,
			})

			err = relGraph.AddEdge(int(foundRange), int(eachRefRange), EdgeAttrReference)
			assert.Nil(t, err)
		}
	}

	// test these graphs
	_ = graph.DFS(factGraph, 4, func(i int) bool {
		vertex, err := factGraph.Vertex(i)
		if err != nil {
			return true
		}
		log.Printf("def in file %d range: %v", vertex.FileId, vertex.Range)

		// any links?
		relVertex, err := relGraph.Vertex(i)
		if err != nil {
			return false
		}
		err = graph.BFS(relGraph, relVertex.GetId(), func(j int) bool {
			cur, err := factGraph.Vertex(j)
			if err != nil {
				return true
			}
			log.Printf("refered by file %d range: %v", cur.FileId, cur.Range)
			return false
		})
		if err != nil {
			return false
		}

		return false
	})
}

func reverseMap(m map[parser.Id]parser.Id) map[parser.Id]parser.Id {
	n := make(map[parser.Id]parser.Id, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
