package main

import (
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
			funcGraph, err := collector.CreateGraph(factStorage, sourceContext)
			panicIfErr(err)

			// line offset
			startPoints := make([]*collector.FuncVertex, 0)
			for path, lines := range lineMap {
				functionFile := factStorage.GetByFile(path)
				if functionFile != nil {
					for _, eachUnit := range functionFile.Units {
						// append these def lines
						if eachUnit.GetSpan().ContainAnyLine(lines...) {
							cur := &collector.FuncVertex{
								Function: eachUnit,
								FuncPos: &collector.FuncPos{
									Path:  path,
									Lang:  string(functionFile.Language),
									Start: int(eachUnit.GetSpan().Start.Row),
									End:   int(eachUnit.GetSpan().End.Row),
								},
							}
							startPoints = append(startPoints, cur)
						}
					}
				}
			}

			// start scan
			for _, eachPtr := range startPoints {
				log.Infof("start point: %v", eachPtr.Id())
				counts := funcGraph.InfluenceCount(eachPtr)
				log.Infof("counts: %d", counts)
			}

			log.Infof("diff finished.")
			return nil
		},
	}
	app.Commands = append(app.Commands, diffCmd)
}
