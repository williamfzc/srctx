package main

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/collector"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/parser"
	"golang.org/x/exp/slices"
)

func AddDiffCmd(app *cli.App) {
	var src string
	var before string
	var after string
	var lsifZip string
	var outputJson string
	var outputCsv string
	var outputDot string

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
		},
		Action: func(cCtx *cli.Context) error {
			// prepare
			lineMap, err := diff.GitDiff(src, before, after)
			panicIfErr(err)
			sourceContext, err := parser.FromLsifFile(lsifZip)
			panicIfErr(err)

			// metadata
			factStorage, err := collector.CreateFact(src)
			panicIfErr(err)

			// line offset
			funcDefLineMap := make(map[string][]int)
			for path, lines := range lineMap {
				functionFile := factStorage.GetByFile(path)
				if functionFile != nil {
					for _, eachUnit := range functionFile.Units {
						// append these def lines
						if eachUnit.GetSpan().ContainAnyLine(lines...) {
							defLine := int(eachUnit.GetSpan().Start.Row + 1)
							if !slices.Contains(funcDefLineMap[path], defLine) {
								funcDefLineMap[path] = append(funcDefLineMap[path], defLine)
							}
						}
					}
				}
			}
			for k, v := range funcDefLineMap {
				log.Infof("file %v append %d lines", k, len(v))
				for _, eachLine := range v {
					if slices.Contains(lineMap[k], eachLine) {
						lineMap[k] = append(lineMap[k], eachLine)
					}
				}
			}

			// calc file ref counts
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
					// which lines will reference this line
					vertices, _ := sourceContext.RefsByLine(path, eachLine)
					// todo: and this line will reference what
					// todo: and a BFS search for specific depth
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

				lineStats = append(lineStats, curFileLines...)
			}
			log.Infof("diff finished.")

			if outputJson != "" {
				exportJson(outputJson, lineStats)
			}
			if outputCsv != "" {
				exportCsv(outputCsv, lineStats)
			}
			if outputDot != "" {
				exportDot(outputDot, fileRefMap)
			}
			return nil
		},
	}
	app.Commands = append(app.Commands, diffCmd)
}
