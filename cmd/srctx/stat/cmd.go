package stat

import (
	"path/filepath"

	"github.com/opensibyl/sibyl2/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/graph"
	"github.com/williamfzc/srctx/parser"
	"github.com/williamfzc/srctx/parser/lsif"
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
	var cacheType string
	var lang string

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
		&cli.StringFlag{
			Name:        "cacheType",
			Value:       lsif.CacheTypeFile,
			Usage:       "mem or file",
			Destination: &cacheType,
		},
		&cli.StringFlag{
			Name:        "lang",
			Value:       string(core.LangUnknown),
			Usage:       "language of repo",
			Destination: &lang,
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

			if cacheType != lsif.CacheTypeFile {
				parser.UseMemCache()
			}

			lang := core.LangType(lang)

			// metadata
			var funcGraph *graph.FuncGraph

			if scipFile != "" {
				// using SCIP
				log.Infof("using SCIP as index")
				funcGraph, err = graph.CreateFuncGraphFromDirWithSCIP(src, scipFile, lang)
			} else {
				// using LSIF
				log.Infof("using LSIF as index")
				if withIndex {
					funcGraph, err = graph.CreateFuncGraphFromGolangDir(src, lang)
				} else {
					funcGraph, err = graph.CreateFuncGraphFromDirWithLSIF(src, lsifZip, lang)
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
				case "dir":
					dirGraph, err := funcGraph.ToDirGraph()
					if err != nil {
						return err
					}
					err = dirGraph.DrawDot(outputDot)
					if err != nil {
						return err
					}
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
