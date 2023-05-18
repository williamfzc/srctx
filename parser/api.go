package parser

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	log "github.com/sirupsen/logrus"
	lsifgo "github.com/sourcegraph/lsif-go/cmd/lsif-go/api"
	"github.com/sourcegraph/scip/bindings/go/scip"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser/lexer"
	"github.com/williamfzc/srctx/parser/lsif"
	"google.golang.org/protobuf/proto"
)

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
	lsifFile := "./dump.lsif"
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

	// change workdir because srctx needs to access the files
	// lsif uses relative paths
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

	log.Infof("parser ready")
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

				// access filesystem
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
	// then rel graph

	// refs
	log.Infof("def ref map size: %d", len(readyParser.Docs.Ranges.DefRefs))
	for eachReferenceResultId, eachDef := range readyParser.Docs.Ranges.DefRefs {
		defFileId := eachDef.DocId
		log.Debugf("def %d in file %s line %d",
			eachReferenceResultId,
			readyParser.Docs.Entries[defFileId],
			eachDef.Line)
		defVertex, err := factGraph.Vertex(int(eachDef.RangeId))
		if err != nil {
			return nil, err
		}
		defVertexInRel := defVertex.ToRelVertex()
		_ = relGraph.AddVertex(defVertexInRel)

		refs := readyParser.Docs.Ranges.References.GetItems(eachReferenceResultId)
		for _, eachRef := range refs {
			log.Debugf("def %v refed in file %s, line %d",
				defVertexInRel,
				readyParser.Docs.Entries[eachRef.DocId],
				eachRef.Line)

			relVertex, err := factGraph.Vertex(int(eachRef.RangeId))
			if err != nil {
				return nil, err
			}
			relVertexInRel := relVertex.ToRelVertex()
			_ = relGraph.AddVertex(relVertexInRel)
			_ = relGraph.AddEdge(defVertexInRel.Id(), relVertexInRel.Id(), object.EdgeAttrReference)
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
	log.Infof("graph ready. fact: %d, rel %d", factSize, relSize)
	log.Infof("file count: %d", len(ret.Files()))

	return &ret, nil
}
