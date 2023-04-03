package srctx

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/lexer"
	"github.com/williamfzc/srctx/parser/lsif"
)

func FromLsifZip(lsifZip string) (*SourceContext, error) {
	file, err := os.Open(lsifZip)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	newParser, err := lsif.NewParser(context.Background(), file)
	if err != nil {
		return nil, err
	}
	defer newParser.Close()
	return FromParser(newParser)
}

func FromParser(readyParser *lsif.Parser) (*SourceContext, error) {
	ret := NewSourceContext()
	factGraph := ret.FactGraph
	relGraph := ret.RelGraph

	// file level
	for eachFileId, eachFile := range readyParser.Docs.Entries {
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

				tokens, err := lexer.File2Tokens(eachFile)
				if err == nil {
					// it's ok
					_ = tokens[rawRange.Line]
					// todo: switch handler here
				}

				eachRangeVertex := &FactVertex{
					DocId:  int(eachRangeId),
					FileId: int(eachFileId),
					Kind:   FactDef,
					Range:  rawRange,
					Extras: nil,
				}
				log.Debugf("file %s range %v", eachFile, rawRange)

				err = factGraph.AddVertex(eachRangeVertex)
				if err != nil {
					return nil, err
				}

				// and edge
				err = factGraph.AddEdge(int(eachFileId), eachRangeVertex.Id(), EdgeAttrContains)
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
		refRanges := make(map[lsif.Id]interface{}, 0)
		for _, eachRef := range refs {
			log.Debugf("def %d refed in file %s, line %d",
				eachReferenceResultId,
				readyParser.Docs.Entries[eachRef.DocId],
				eachRef.Line)

			refRanges[eachRef.RangeId] = nil
		}

		// connect between ranges to ranges
		// range - next -> resultSet - text/references -> referenceResult - item -> range
		for eachRefRange := range refRanges {
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

			_ = relGraph.AddVertex(&RelVertex{
				DocId: int(eachRefRange),
				Kind:  FactRef,
			})
			_ = relGraph.AddVertex(&RelVertex{
				DocId: int(foundRange),
				Kind:  FactDef,
			})

			err := relGraph.AddEdge(int(foundRange), int(eachRefRange), EdgeAttrReference)
			if err != nil {
				return nil, err
			}
		}
	}
	return &ret, nil
}
