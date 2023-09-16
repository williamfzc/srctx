package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx"
	"github.com/williamfzc/srctx/cmd/srctx/diff"
)

func main() {
	mainFunc(os.Args)
}

func mainFunc(args []string) {
	app := cli.NewApp()
	app.Name = "srctx"
	app.Usage = "source context tool"

	diff.AddDiffCmd(app)
	diff.AddConfigCmd(app)

	log.Infof("srctx version %v (%s)", srctx.Version, srctx.RepoUrl)
	err := app.Run(args)
	panicIfErr(err)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	environment := os.Getenv("SRCTX_ENV")

	if environment == "production" || environment == "prod" {
		// can be saved to log file
		log.SetOutput(os.Stdout)
		log.SetLevel(log.WarnLevel)
	}
}
