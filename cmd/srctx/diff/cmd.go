package diff

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/parser/lsif"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  srcFlagName,
		Value: ".",
		Usage: "project path",
	},
	&cli.StringFlag{
		Name:  repoRootFlagName,
		Value: "",
		Usage: "root path of your repo",
	},
	&cli.StringFlag{
		Name:  beforeFlagName,
		Value: "HEAD~1",
		Usage: "before rev",
	},
	&cli.StringFlag{
		Name:  afterFlagName,
		Value: "HEAD",
		Usage: "after rev",
	},
	&cli.StringFlag{
		Name:  lsifFlagName,
		Value: "./dump.lsif",
		Usage: "lsif path, can be zip or origin file",
	},
	&cli.StringFlag{
		Name:  scipFlagName,
		Value: "",
		Usage: "scip file",
	},
	&cli.StringFlag{
		Name:  nodeLevelFlagName,
		Value: nodeLevelFunc,
		Usage: "graph level (file or func)",
	},
	&cli.StringFlag{
		Name:  outputJsonFlagName,
		Value: "",
		Usage: "json output",
	},
	&cli.StringFlag{
		Name:  outputCsvFlagName,
		Value: "",
		Usage: "csv output",
	},
	&cli.StringFlag{
		Name:  outputDotFlagName,
		Value: "",
		Usage: "reference dot file output",
	},
	&cli.StringFlag{
		Name:  outputHtmlFlagName,
		Value: "",
		Usage: "render html report with g6",
	},
	&cli.BoolFlag{
		Name:  withIndexFlagName,
		Value: false,
		Usage: "create indexes first if enabled, currently support golang only",
	},
	&cli.StringFlag{
		Name:  cacheTypeFlagName,
		Value: lsif.CacheTypeFile,
		Usage: "mem or file",
	},
	&cli.StringFlag{
		Name:  langFlagName,
		Value: string(core.LangUnknown),
		Usage: "language of repo",
	},
	&cli.BoolFlag{
		Name:  noDiffFlagName,
		Value: false,
		Usage: "will not calc git diff if enabled",
	},
	&cli.StringFlag{
		Name:  indexCmdFlagName,
		Value: "",
		Usage: "specific scip or lsif cmd",
	},
}

func AddDiffCmd(app *cli.App) {
	diffCmd := &cli.Command{
		Name:  "diff",
		Usage: "diff with lsif",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			opts := NewOptionsFromCliFlags(cCtx)
			// standardize the path
			src, err := filepath.Abs(opts.Src)
			if err != nil {
				return err
			}
			optsFromSrc, err := NewOptionsFromSrc(src)
			if err != nil {
				// ok
				log.Infof("no config: %v", err)
			} else {
				log.Infof("config file found")
				opts = optsFromSrc
			}
			opts.Src = src

			err = MainDiff(opts)
			if err != nil {
				return err
			}
			return nil
		},
	}
	app.Commands = append(app.Commands, diffCmd)
}

func AddConfigCmd(app *cli.App) {
	configCmd := &cli.Command{
		Name:  "diffcfg",
		Usage: "create config file for diff",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			opts := NewOptionsFromCliFlags(cCtx)
			jsonContent, err := json.Marshal(opts)
			if err != nil {
				return fmt.Errorf("failed to marshal options: %w", err)
			}

			configFile := filepath.Join(opts.Src, DefaultConfigFile)
			if err := os.WriteFile(configFile, jsonContent, 0o644); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			log.Infof("create config file finished: %v", configFile)
			return nil
		},
	}
	app.Commands = append(app.Commands, configCmd)
}
