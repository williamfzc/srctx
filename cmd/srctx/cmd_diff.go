package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/gocarina/gocsv"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/collector"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/parser"
)

func AddDiffCmd(app *cli.App) {
	var src string
	var before string
	var after string
	var lsifZip string
	var outputJson string
	var outputCsv string
	var outputDot string
	var funcEnhance bool

	diffCmd := &cli.Command{
		Name:  "diff",
		Usage: "diff with lsif",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "src",
				Value:       ".",
				Usage:       "repo path",
				Destination: &src,
			},
			&cli.StringFlag{
				Name:        "before",
				Value:       "HEAD~1",
				Usage:       "before rev",
				Destination: &before,
			},
			&cli.StringFlag{
				Name:        "after",
				Value:       "HEAD",
				Usage:       "after rev",
				Destination: &after,
			},
			&cli.StringFlag{
				Name:        "lsif",
				Value:       "./dump.lsif",
				Usage:       "lsif path, can be zip or origin file",
				Destination: &lsifZip,
			},
			&cli.StringFlag{
				Name:        "outputJson",
				Value:       "",
				Usage:       "json output",
				Destination: &outputJson,
			},
			&cli.StringFlag{
				Name:        "outputCsv",
				Value:       "srctx-diff.csv",
				Usage:       "csv output",
				Destination: &outputCsv,
			},
			&cli.StringFlag{
				Name:        "outputDot",
				Value:       "",
				Usage:       "reference dot file output",
				Destination: &outputDot,
			},
			// experimental
			&cli.BoolFlag{
				Name:        "funcEnhance",
				Value:       true,
				Usage:       "function level diff json",
				Destination: &funcEnhance,
			},
		},
		Action: func(cCtx *cli.Context) error {
			// prepare
			lineMap, err := diff.GitDiff(src, before, after)
			panicIfErr(err)
			sourceContext, err := parser.FromLsifFile(lsifZip)
			panicIfErr(err)

			// calc
			lineStats := make([]*LineStat, 0)
			log.Infof("diff / total (files): %d / %d", len(lineMap), len(sourceContext.Files()))

			fileRefMap := make(map[string]*fileVertex)
			// directly
			for path := range lineMap {
				fileRefMap[path] = &fileVertex{
					Name:     path,
					Refs:     nil,
					Directly: true,
				}
			}

			for path, lines := range lineMap {
				curFileLines := make([]*LineStat, 0, len(lines))
				for _, eachLine := range lines {
					lineStat := NewLineStat(path, eachLine)
					vertices, _ := sourceContext.RefsByLine(path, eachLine)
					log.Debugf("path %s line %d affected %d vertexes", path, eachLine, len(vertices))

					lineStat.RefScope.TotalRefCount = len(vertices)
					for _, eachVertex := range vertices {
						refFileName := sourceContext.FileName(eachVertex.FileId)

						// indirectly
						if v, ok := fileRefMap[refFileName]; ok {
							v.Refs = append(v.Refs, refFileName)
						} else {
							fileRefMap[refFileName] = &fileVertex{
								Name:     refFileName,
								Refs:     []string{path},
								Directly: false,
							}
						}

						if refFileName != path {
							lineStat.RefScope.CrossFileRefCount++
						}
						if filepath.Dir(refFileName) != filepath.Dir(path) {
							lineStat.RefScope.CrossDirRefCount++
						}
					}
					curFileLines = append(curFileLines, lineStat)
				}

				if funcEnhance {
					functionFile, err := collector.GetFunctionMetadataFromFile(filepath.Join(src, path))
					if _, ok := err.(*collector.NotSupportLangError); ok {
						log.Warnf("file %v not supported", err)
						goto eachFileEnd
					}

					if err != nil {
						return err
					}
					// what happened in these lines
					influences := make([]*object.Function, 0)
					for _, eachUnit := range functionFile.Units {
						if eachUnit.GetSpan().ContainAnyLine(lines...) {
							influences = append(influences, eachUnit)
						}
					}
					for _, eachFunc := range influences {
						// def line (not precise
						funcDefLine := eachFunc.GetSpan().Start.Row + 1
						vertices, _ := sourceContext.RefsByLine(path, int(funcDefLine))
						log.Debugf("[func scope] path %s line %d affected %d vertexes", path, funcDefLine, len(vertices))
						// update status
						for _, eachLine := range curFileLines {
							eachLine.FuncRefScope.TotalFuncRefCount += len(vertices)
							for _, eachVertex := range vertices {
								fileName := sourceContext.FileName(eachVertex.FileId)
								if fileName != path {
									eachLine.FuncRefScope.CrossFuncFileRefCount++
								}
								if filepath.Dir(fileName) != filepath.Dir(path) {
									eachLine.FuncRefScope.CrossFuncDirRefCount++
								}
							}
						}
					}
				}
			eachFileEnd:
				lineStats = append(lineStats, curFileLines...)
			}
			log.Infof("diff finished.")

			if outputJson != "" {
				data, err := json.Marshal(lineStats)
				panicIfErr(err)
				err = os.WriteFile(outputJson, data, 0644)
				panicIfErr(err)
				log.Infof("dump json to %s", outputJson)
			}

			if outputCsv != "" {
				csvFile, err := os.OpenFile(outputCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
				panicIfErr(err)
				defer csvFile.Close()

				unsafeLines := make([]*LineStat, 0)
				for _, each := range lineStats {
					if !each.IsSafe() {
						unsafeLines = append(unsafeLines, each)
					}
				}

				if err := gocsv.MarshalFile(&unsafeLines, csvFile); err != nil { // Load clients from file
					panic(err)
				}
				log.Infof("dump csv to %s", outputCsv)
			}

			if outputDot != "" {
				// only create a file level graph
				fileGraph := graph.New((*fileVertex).Id, graph.Directed())
				for _, vertex := range fileRefMap {
					if vertex.Directly {
						_ = fileGraph.AddVertex(vertex, func(vertexProperties *graph.VertexProperties) {
							vertexProperties.Attributes["style"] = "filled"
							vertexProperties.Attributes["fillcolor"] = "yellow"
						})
					} else {
						_ = fileGraph.AddVertex(vertex)
					}
				}
				for _, vertex := range fileRefMap {
					for _, eachRef := range vertex.Refs {
						// ignore self ref
						if eachRef != vertex.Id() {
							_ = fileGraph.AddEdge(eachRef, vertex.Id())
						}
					}
				}
				f, _ := os.Create(outputDot)
				_ = draw.DOT(fileGraph, f)
				log.Infof("dump dot to %s", outputDot)
			}
			return nil
		},
	}
	app.Commands = append(app.Commands, diffCmd)
}

type fileVertex struct {
	Name     string
	Refs     []string
	Directly bool
}

func (vertex *fileVertex) Id() string {
	return vertex.Name
}
