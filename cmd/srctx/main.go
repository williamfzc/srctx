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
	// random edit
	app.Name = "srctx1"
	app.Usage = "source context tool2"

	diff.AddDiffCmd(app)

	err := app.Run(args)
	panicIfErr(err)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
