package dump

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/parser"
	"os"
	"slices"
	"strconv"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "src",
		Value: ".",
		Usage: "project path",
	},
	&cli.StringFlag{
		Name:  "csv",
		Value: "output.csv",
		Usage: "output csv file",
	},
}

type Options struct {
	Src string `json:"src"`
	Csv string `json:"csv"`
}

func AddDumpCmd(app *cli.App) {
	dumpCmd := &cli.Command{
		Name:  "dump",
		Usage: "dump file relations",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			src := cCtx.String("src")
			csvPath := cCtx.String("csv")

			sourceContext, err := parser.FromGolangSrc(src)
			if err != nil {
				panic(err)
			}

			files := sourceContext.Files()
			slices.Sort(files)
			log.Infof("files in lsif: %d", len(files))

			fileCount := len(files)
			relationMatrix := make([][]int, fileCount)
			for i := range relationMatrix {
				relationMatrix[i] = make([]int, fileCount)
			}

			fileIndexMap := make(map[string]int)
			for idx, file := range files {
				fileIndexMap[file] = idx
			}

			for _, file := range files {
				defs, err := sourceContext.DefsByFileName(file)
				if err != nil {
					panic(err)
				}

				for _, def := range defs {
					refs, err := sourceContext.RefsFromDefId(def.Id())
					if err != nil {
						panic(err)
					}

					for _, ref := range refs {
						refFile := sourceContext.FileName(ref.FileId)
						if refFile == file {
							continue
						}
						if refFile != "" {
							relationMatrix[fileIndexMap[file]][fileIndexMap[refFile]]++
						}
					}
				}
			}

			csvFile, err := os.Create(csvPath)
			if err != nil {
				panic(err)
			}
			defer csvFile.Close()

			writer := csv.NewWriter(csvFile)
			defer writer.Flush()

			header := append([]string{""}, files...)
			if err := writer.Write(header); err != nil {
				panic(err)
			}

			for i, row := range relationMatrix {
				csvRow := make([]string, len(row)+1)
				csvRow[0] = files[i]
				for j, val := range row {
					if val == 0 {
						csvRow[j+1] = ""
					} else {
						csvRow[j+1] = strconv.Itoa(val)
					}
				}
				if err := writer.Write(csvRow); err != nil {
					panic(err)
				}
			}

			log.Infof("CSV file generated successfully.")
			return nil
		},
	}
	app.Commands = append(app.Commands, dumpCmd)
}
