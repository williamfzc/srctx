package main

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/graph"
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
			// standardize the path
			src, err := filepath.Abs(src)
			if err != nil {
				return err
			}

			// prepare
			lineMap, err := diff.GitDiff(src, before, after)
			if err != nil {
				return err
			}

			// metadata
			funcGraph, err := graph.CreateFuncGraphFromDir(src, lsifZip)
			panicIfErr(err)

			// line offset
			startPoints := make([]*graph.FuncVertex, 0)
			for path, lines := range lineMap {
				functions := funcGraph.GetFunctionsByFile(path)
				if functions != nil {
					for _, eachFunc := range functions {
						// append these def lines
						if eachFunc.GetSpan().ContainAnyLine(lines...) {
							startPoints = append(startPoints, eachFunc)
						}
					}
				}
			}

			// start scan
			for _, eachPtr := range startPoints {
				log.Infof("start point: %v", eachPtr.Id())

				ids := funcGraph.TransitiveReferencedIds(eachPtr)
				log.Infof("counts: %d", len(ids))
				for _, each := range ids {
					err := funcGraph.Highlight(each)
					if err != nil {
						return err
					}
				}
			}
			for _, eachPtr := range startPoints {
				err := funcGraph.FillWithRed(eachPtr.Id())
				if err != nil {
					return err
				}
			}

			log.Infof("diff finished.")

			// output
			if outputDot != "" {
				err := funcGraph.DrawDot(outputDot)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	app.Commands = append(app.Commands, diffCmd)
}
