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
	// random edit
	app.Name = "srctx1"
	app.Usage = "source context tool2"

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
