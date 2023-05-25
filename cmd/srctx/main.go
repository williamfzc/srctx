package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx"
	"github.com/williamfzc/srctx/cmd/srctx/diff"
	"github.com/williamfzc/srctx/cmd/srctx/stat"
)

func main() {
	mainFunc(os.Args)
}

func mainFunc(args []string) {
	app := cli.NewApp()
	app.Name = "srctx"
	app.Usage = "source context tool"

	diff.AddDiffCmd(app)
	stat.AddStatCmd(app)

	log.Infof("srctx version %v", srctx.Version)
	err := app.Run(args)
	panicIfErr(err)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
