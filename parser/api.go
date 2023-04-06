package parser

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser/lexer"
	"github.com/williamfzc/srctx/parser/lsif"
)

func FromLsifFile(lsifFile string) (*object.SourceContext, error) {
	file, err := os.Open(lsifFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var newParser *lsif.Parser
	if strings.HasSuffix(lsifFile, ".zip") {
		newParser, err = lsif.NewParser(context.Background(), file)
	} else {
		newParser, err = lsif.NewParserRaw(context.Background(), file)
	}
	if err != nil {
		return nil, err
	}
	defer newParser.Close()
	return FromParser(newParser)
}

func FromParser(readyParser *lsif.Parser) (*object.SourceContext, error) {
	ret := object.NewSourceContext()
	factGraph := ret.FactGraph
	relGraph := ret.RelGraph

	// file level
	for eachFileId, eachFile := range readyParser.Docs.Entries {
		eachFileVertex := &object.FactVertex{
			DocId:  int(eachFileId),
			FileId: int(eachFileId),
			Range:  nil,
			Kind:   object.FactFile,
			Extras: &object.FileExtras{
				Path: eachFile,
			},
		}
		err := factGraph.AddVertex(eachFileVertex)
		if err != nil {
			return nil, err
		}
		ret.FileMapping[eachFile] = int(eachFileId)
	}

	// contains / defs
	for eachFileId, eachFile := range readyParser.Docs.Entries {
		// def ranges in this file
		if inVs, ok := readyParser.Docs.DocRanges[eachFileId]; ok {
			for _, eachRangeId := range inVs {
				rawRange := &lsif.Range{}
				err := readyParser.Docs.Ranges.Cache.Entry(eachRangeId, rawRange)
				if err != nil {
					return nil, err
				}

				defExtras := &object.DefExtras{}
				tokens, err := lexer.File2Tokens(eachFile)
				if err == nil {
					log.Debugf("file %s, tokens: %d, cur line: %d", eachFile, len(tokens), rawRange.Line)
					if int(rawRange.Line) >= len(tokens) {
						log.Warnf("access out of range: %d", rawRange.Line)
						continue
					}
					curLineTokens := tokens[rawRange.Line]
					defType := lexer.TypeFromTokens(curLineTokens)
					defExtras.DefType = defType
					defExtras.RawTokens = curLineTokens
				}

				eachRangeVertex := &object.FactVertex{
					DocId:  int(eachRangeId),
					FileId: int(eachFileId),
					Kind:   object.FactDef,
					Range:  rawRange,
					Extras: defExtras,
				}
				log.Debugf("file %s range %v", eachFile, rawRange)

				err = factGraph.AddVertex(eachRangeVertex)
				if err != nil {
					return nil, err
				}

				// and edge
				err = factGraph.AddEdge(int(eachFileId), eachRangeVertex.Id(), object.EdgeAttrContains)
				if err != nil {
					return nil, err
				}
			}
		} else {
			log.Warnf("no any links to %d", eachFileId)
		}
	}

	// the fact graph is ready
	// then real graph

	// refs
	reverseNextMap := reverseMap(readyParser.Docs.Ranges.NextMap)
	reverseRefMap := reverseMap(readyParser.Docs.Ranges.TextReferenceMap)
	for eachReferenceResultId, eachDef := range readyParser.Docs.Ranges.DefRefs {
		refFileId := eachDef.DocId
		log.Debugf("def %d in file %s line %d",
			eachReferenceResultId,
			readyParser.Docs.Entries[refFileId],
			eachDef.Line)

		refs := readyParser.Docs.Ranges.References.GetItems(eachReferenceResultId)
		refRanges := make(map[lsif.Id]lsif.Id, 0)
		for _, eachRef := range refs {
			log.Debugf("def %d refed in file %s, line %d",
				eachReferenceResultId,
				readyParser.Docs.Entries[eachRef.DocId],
				eachRef.Line)

			refRanges[eachRef.RangeId] = eachRef.DocId
		}

		// connect between ranges to ranges
		// range - next -> resultSet - text/references -> referenceResult - item -> range
		for eachRefRange, eachFileId := range refRanges {
			// starts with the ref point
			resultSetId, ok := readyParser.Docs.Ranges.NextMap[eachRefRange]
			if !ok {
				log.Warnf("failed to jump with nextMap: %v", eachRefRange)
				continue
			}
			foundReferenceResultId, ok := readyParser.Docs.Ranges.TextReferenceMap[resultSetId]
			if !ok {
				log.Warnf("failed to jump with reference map: %v", resultSetId)
				continue
			}

			foundItem, ok := reverseRefMap[foundReferenceResultId]
			if !ok {
				log.Warnf("failed to jump with rev ref map: %v", resultSetId)
				continue
			}

			foundRange, ok := reverseNextMap[foundItem]
			if !ok {
				log.Warnf("failed to jump with rev next map: %v", resultSetId)
				continue
			}

			// ref file
			eachFileVertex := &object.RelVertex{
				DocId: int(eachFileId),
				Kind:  object.FactFile,
			}
			rawRange := &lsif.Range{}
			err := readyParser.Docs.Ranges.Cache.Entry(eachRefRange, rawRange)

			eachRefVertex := &object.RelVertex{
				DocId:  int(eachRefRange),
				FileId: eachFileVertex.DocId,
				Kind:   object.FactRef,
				Range:  rawRange,
			}

			_ = relGraph.AddVertex(eachFileVertex, func(vertexProperties *graph.VertexProperties) {
				vertexProperties.Attributes["label"] = ret.FileName(eachFileVertex.DocId)
			})
			_ = relGraph.AddVertex(eachRefVertex, func(vertexProperties *graph.VertexProperties) {
				vertexProperties.Attributes["label"] = fmt.Sprintf("%d:%d", eachRefVertex.LineNumber(), eachRefVertex.Range.Character)
			})
			// edge between file and ref
			_ = relGraph.AddEdge(eachFileVertex.Id(), eachRefVertex.Id(), object.EdgeAttrContains)

			// edge between ref and def
			defRange := &lsif.Range{}
			err = readyParser.Docs.Ranges.Cache.Entry(foundRange, defRange)
			if err != nil {
				return nil, err
			}

			eachDefVertex := &object.RelVertex{
				DocId: int(foundRange),
				Kind:  object.FactDef,
				Range: defRange,
			}
			_ = relGraph.AddVertex(eachDefVertex, func(vertexProperties *graph.VertexProperties) {
				vertexProperties.Attributes["label"] = fmt.Sprintf("%d:%d", eachDefVertex.LineNumber(), eachDefVertex.Range.Character)
			})
			err = relGraph.AddEdge(eachRefVertex.Id(), eachDefVertex.Id(), object.EdgeAttrReference)

			// edge between file and def
			// find its belonging
			eachDefVertexInFact, err := factGraph.Vertex(eachDefVertex.Id())
			if err != nil {
				return nil, err
			}
			defFileVertexInFact, err := factGraph.Vertex(eachDefVertexInFact.FileId)
			if err != nil {
				return nil, err
			}
			eachDefFileVertexInRef := defFileVertexInFact.ToRelVertex()

			if err != nil {
				return nil, err
			}
			_ = relGraph.AddVertex(eachDefFileVertexInRef, func(vertexProperties *graph.VertexProperties) {
				vertexProperties.Attributes["label"] = ret.FileName(eachDefFileVertexInRef.DocId)
			})
			_ = relGraph.AddEdge(eachDefFileVertexInRef.Id(), eachDefVertex.Id(), object.EdgeAttrContains)

			if err != nil {
				return nil, err
			}
		}
	}
	return &ret, nil
}
