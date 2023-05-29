package parser

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	log "github.com/sirupsen/logrus"
	lsifgo "github.com/sourcegraph/lsif-go/cmd/lsif-go/api"
	"github.com/sourcegraph/scip/bindings/go/scip"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser/lsif"
	"google.golang.org/protobuf/proto"
)

func UseTempFileCache() {
	lsif.CacheType = lsif.CacheTypeFile
}

func UseMemCache() {
	lsif.CacheType = lsif.CacheTypeMem
}

func FromGolangSrc(srcDir string) (*object.SourceContext, error) {
	// change workdir because lsif needs to access the files
	originWorkdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = os.Chdir(srcDir)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Chdir(originWorkdir)
	}()
	// create indexes
	log.Infof("creating index for %v", srcDir)
	err = lsifgo.MainArgs([]string{
		"-v",
		"--project-root", srcDir,
		"--repository-root", srcDir,
		"--module-root", srcDir,
	})
	if err != nil {
		return nil, err
	}
	return FromLsifFile("./dump.lsif", srcDir)
}

func FromScipFile(scipFile string, srcDir string) (*object.SourceContext, error) {
	scipIndex, err := readFromOption(scipFile)
	if err != nil {
		return nil, err
	}

	lsifIndex, err := scip.ConvertSCIPToLSIF(scipIndex)
	if err != nil {
		return nil, err
	}

	// still save this file for debug
	lsifFile := filepath.Join(srcDir, "dump.lsif")
	lsifWriter, err := os.OpenFile(lsifFile, os.O_WRONLY|os.O_CREATE, 0666)
	defer lsifWriter.Close()

	err = scip.WriteNDJSON(scip.ElementsToJsonElements(lsifIndex), lsifWriter)
	if err != nil {
		return nil, err
	}

	log.Infof("scip -> lsif converted")
	return FromLsifFile(lsifFile, srcDir)
}

// https://github.com/williamfzc/scip/blob/main/cmd/option_from.go
func readFromOption(fromPath string) (*scip.Index, error) {
	var scipReader io.Reader
	if fromPath == "-" {
		scipReader = os.Stdin
	} else if !strings.HasSuffix(fromPath, ".scip") && !strings.HasSuffix(fromPath, ".lsif-typed") {
		return nil, errors.Newf("expected file with .scip extension but found %s", fromPath)
	} else {
		scipFile, err := os.Open(fromPath)
		defer scipFile.Close()
		if err != nil {
			return nil, err
		}
		scipReader = scipFile
	}

	scipBytes, err := io.ReadAll(scipReader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read SCIP index at path %s", fromPath)
	}

	scipIndex := scip.Index{}
	err = proto.Unmarshal(scipBytes, &scipIndex)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse SCIP index at path %s", fromPath)
	}
	return &scipIndex, nil
}

func FromLsifFile(lsifFile string, srcDir string) (*object.SourceContext, error) {
	log.Infof("read lsif file: %v", lsifFile)
	file, err := os.Open(lsifFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var newParser *lsif.Parser
	ctx := context.Background()
	if strings.HasSuffix(lsifFile, ".zip") {
		newParser, err = lsif.NewParser(ctx, file)
	} else {
		newParser, err = lsif.NewParserRaw(ctx, file)
	}
	if err != nil {
		return nil, err
	}
	defer newParser.Close()

	log.Infof("index parser ready")
	return FromParser(newParser)
}

func FromParser(readyParser *lsif.Parser) (*object.SourceContext, error) {
	log.Infof("cache type: %v", lsif.CacheType)
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

		// this file path can be relative
		eachFile = filepath.ToSlash(filepath.Clean(eachFile))
		ret.FileMapping[eachFile] = int(eachFileId)
	}

	log.Infof("creating fact graph")
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

				// access filesystem
				eachRangeVertex := &object.FactVertex{
					DocId:  int(eachRangeId),
					FileId: int(eachFileId),
					Kind:   object.FactDef,
					Range:  rawRange,
					Extras: &object.DefExtras{},
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
	// then rel graph
	log.Infof("creating rel graph, reference result map size: %d", len(readyParser.Docs.Ranges.DefRefs))
	for eachReferenceResultId, eachRef := range readyParser.Docs.Ranges.DefRefs {
		refFileId := eachRef.DocId
		log.Debugf("reference result %d in file %s line %d",
			eachReferenceResultId,
			readyParser.Docs.Entries[refFileId],
			eachRef.Line)

		refs := readyParser.Docs.Ranges.References.GetItems(eachReferenceResultId)
		if refs == nil {
			// eachReferenceResultId is definitionResult
			// try to search referenceResult
			// definitionResult -> resultSet -> referenceResult
			resultSetId, ok := readyParser.Docs.Ranges.NextMap[eachRef.RangeId]
			if !ok {
				// should not happen
				return nil, fmt.Errorf("no result set id found for range: %d", eachRef.RangeId)
			}
			// reset and update
			eachReferenceResultId = readyParser.Docs.Ranges.TextReferenceMap[resultSetId]
			refs = readyParser.Docs.Ranges.References.GetItems(eachReferenceResultId)
		}

		refRanges := make(map[lsif.Id]lsif.Id, 0)
		for _, eachRef := range refs {
			log.Debugf("reference result %d refed in file %s, line %d",
				eachReferenceResultId,
				readyParser.Docs.Entries[eachRef.DocId],
				eachRef.Line)

			refRanges[eachRef.RangeId] = eachRef.DocId
		}

		// search the definition:
		// ref range - next -> resultSet
		// resultSet - text/definition -> definitionResult
		// definitionResult - item/edge -> def ranges
		for eachRefRange := range refRanges {
			// starts with the ref point
			resultSetId, ok := readyParser.Docs.Ranges.NextMap[eachRefRange]
			if !ok {
				log.Warnf("failed to jump with nextMap: %v", eachRefRange)
				continue
			}
			foundDefinitionResult, ok := readyParser.Docs.Ranges.TextDefinitionMap[resultSetId]
			if !ok {
				log.Warnf("failed to jump with reference map: %v", resultSetId)
				continue
			}

			rawRange := &lsif.Range{}
			err := readyParser.Docs.Ranges.Cache.Entry(eachRefRange, rawRange)
			if err != nil {
				return nil, err
			}
			eachRefVertex := &object.RelVertex{
				DocId: int(eachRefRange),
				Kind:  object.RelReference,
				Range: rawRange,
			}
			log.Debugf("ref vertex: %v %v", eachRefVertex, rawRange)
			_ = relGraph.AddVertex(eachRefVertex)

			edgeToDefRanges, ok := readyParser.Docs.Ranges.RawEdgeMap[foundDefinitionResult]
			if !ok {
				log.Warnf("failed to jump with raw edge map: %v", resultSetId)
				continue
			}
			// only one range, actually
			for _, edgeToDefRange := range edgeToDefRanges {
				defRanges := edgeToDefRange.RangeIds
				for _, defRangeId := range defRanges {
					defRange := &lsif.Range{}
					err = readyParser.Docs.Ranges.Cache.Entry(defRangeId, defRange)
					if err != nil {
						return nil, err
					}
					defVertex := &object.RelVertex{
						DocId:  int(defRangeId),
						FileId: 0,
						Kind:   object.RelReference,
						Range:  defRange,
					}
					_ = relGraph.AddVertex(defVertex)
					// definition -> reference
					_ = relGraph.AddEdge(defVertex.Id(), eachRefVertex.Id())
				}
			}
		}
	}

	factSize, err := factGraph.Size()
	if err != nil {
		return nil, err
	}
	relSize, err := relGraph.Size()
	if err != nil {
		return nil, err
	}
	log.Infof("base graph ready. fact: %d, rel %d", factSize, relSize)
	log.Infof("file count: %d", len(ret.Files()))

	return &ret, nil
}
