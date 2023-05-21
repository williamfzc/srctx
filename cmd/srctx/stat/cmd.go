package stat

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/graph"
)

func AddStatCmd(app *cli.App) {
	var src string
	var lsifZip string
	var scipFile string
	var outputJson string
	var outputCsv string
	var outputDot string
	var nodeLevel string
	var withIndex bool

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "src",
			Value:       ".",
			Usage:       "repo path",
			Destination: &src,
		},
		&cli.StringFlag{
			Name:        "lsif",
			Value:       "./dump.lsif",
			Usage:       "lsif path, can be zip or origin file",
			Destination: &lsifZip,
		},
		&cli.StringFlag{
			Name:        "scip",
			Value:       "",
			Usage:       "scip file",
			Destination: &scipFile,
		},
		&cli.StringFlag{
			Name:        "outputJson",
			Value:       "",
			Usage:       "json output",
			Destination: &outputJson,
		},
		&cli.StringFlag{
			Name:        "outputCsv",
			Value:       "",
			Usage:       "csv output",
			Destination: &outputCsv,
		},
		&cli.StringFlag{
			Name:        "outputDot",
			Value:       "",
			Usage:       "reference dot file output",
			Destination: &outputDot,
		},
		&cli.StringFlag{
			Name:        "nodeLevel",
			Value:       "func",
			Usage:       "graph level (file or func)",
			Destination: &nodeLevel,
		},
		&cli.BoolFlag{
			Name:        "withIndex",
			Value:       false,
			Usage:       "create indexes first if enabled, currently support golang only",
			Destination: &withIndex,
		},
	}

	statCmd := &cli.Command{
		Name:  "stat",
		Usage: "",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			// standardize the path
			src, err := filepath.Abs(src)
			if err != nil {
				return err
			}

			// metadata
			var funcGraph *graph.FuncGraph

			if scipFile != "" {
				// using SCIP
				log.Infof("using SCIP as index")
				funcGraph, err = graph.CreateFuncGraphFromDirWithSCIP(src, scipFile)
			} else {
				// using LSIF
				log.Infof("using LSIF as index")
				if withIndex {
					funcGraph, err = graph.CreateFuncGraphFromGolangDir(src)
				} else {
					funcGraph, err = graph.CreateFuncGraphFromDirWithLSIF(src, lsifZip)
				}
			}

			if err != nil {
				return err
			}
			// output
			if outputDot != "" {
				log.Infof("creating dot file: %v", outputDot)

				var err error
				switch nodeLevel {
				case "func":
					err = funcGraph.DrawDot(outputDot)
				case "file":
					fileGraph, err := funcGraph.ToFileGraph()
					if err != nil {
						return err
					}
					err = fileGraph.DrawDot(outputDot)
				}

				if err != nil {
					return err
				}
			}

			log.Infof("everything done.")
			return nil
		},
	}
	app.Commands = append(app.Commands, statCmd)
}
