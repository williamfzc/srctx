package main

import (
	"os"

	"github.com/urfave/cli/v2"
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

	err := app.Run(args)
	panicIfErr(err)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
