package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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
				Usage:       "srctx json output",
				Destination: &outputJson,
			},
			&cli.StringFlag{
				Name:        "outputCsv",
				Value:       "srctx-diff.csv",
				Usage:       "srctx csv output",
				Destination: &outputCsv,
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
			for path, lines := range lineMap {
				for _, eachLine := range lines {
					lineStat := NewLineStat(path, eachLine)
					vertices, _ := sourceContext.RefsByLine(path, eachLine)
					log.Debugf("path %s line %d affected %d vertexes", path, eachLine, len(vertices))

					lineStat.RefScope.TotalRefCount = len(vertices)
					for _, eachVertex := range vertices {
						refFileName := sourceContext.FileName(eachVertex.FileId)
						if refFileName != path {
							lineStat.RefScope.CrossFileRefCount++
						}
						if filepath.Dir(refFileName) != filepath.Dir(path) {
							lineStat.RefScope.CrossDirRefCount++
						}
					}
					lineStats = append(lineStats, lineStat)
				}
			}
			log.Infof("diff finished.")

			if outputJson != "" {
				data, err := json.Marshal(lineStats)
				panicIfErr(err)
				err = os.WriteFile(outputJson, data, 0644)
				panicIfErr(err)
			}

			if outputCsv != "" {
				csvFile, err := os.OpenFile(outputCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
				panicIfErr(err)
				defer csvFile.Close()
				if err := gocsv.MarshalFile(&lineStats, csvFile); err != nil { // Load clients from file
					panic(err)
				}
			}

			log.Infof("dump finished.")
			return nil
		},
	}
	app.Commands = append(app.Commands, diffCmd)
}
