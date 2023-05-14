package diff

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
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
	var withIndex bool

	flags := []cli.Flag{
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
		&cli.BoolFlag{
			Name:        "withIndex",
			Value:       false,
			Usage:       "create indexes first if enabled, currently support golang only",
			Destination: &withIndex,
		},
	}

	diffCmd := &cli.Command{
		Name:  "diff",
		Usage: "diff with lsif",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			// standardize the path
			src, err := filepath.Abs(src)
			if err != nil {
				return err
			}
			log.Infof("start diffing: %v", src)

			// prepare
			lineMap, err := diff.GitDiff(src, before, after)
			if err != nil {
				return err
			}

			// metadata
			var funcGraph *graph.FuncGraph
			if withIndex {
				funcGraph, err = graph.CreateFuncGraphFromGolangDir(src)
			} else {
				funcGraph, err = graph.CreateFuncGraphFromDir(src, lsifZip)
			}
			if err != nil {
				return err
			}

			// look up start points
			startPoints := make([]*graph.FuncVertex, 0)
			for path, lines := range lineMap {
				curPoints := funcGraph.GetFunctionsByFileLines(path, lines)
				startPoints = append(startPoints, curPoints...)
			}

			// start scan
			stats := make([]*graph.VertexStat, 0)
			for _, eachPtr := range startPoints {
				eachStat := funcGraph.Stat(eachPtr)
				stats = append(stats, eachStat)
				log.Infof("start point: %v, refed: %d, ref: %d", eachPtr.Id(), eachStat.Referenced, eachStat.Reference)
			}
			log.Infof("diff finished.")

			// output
			if outputDot != "" {
				log.Infof("creating dot file: %v", outputDot)
				// colorful
				for _, eachStat := range stats {
					for _, eachVisited := range eachStat.VisitedIds() {
						err := funcGraph.Highlight(eachVisited)
						if err != nil {
							return err
						}
					}
				}
				for _, eachStat := range stats {
					err := funcGraph.FillWithRed(eachStat.Root.Id())
					if err != nil {
						return err
					}
				}

				err := funcGraph.DrawDot(outputDot)
				if err != nil {
					return err
				}
			}

			if outputCsv != "" || outputJson != "" {
				// need to access files
				originWorkdir, err := os.Getwd()
				if err != nil {
					return err
				}
				err = os.Chdir(src)
				if err != nil {
					return err
				}
				defer func() {
					_ = os.Chdir(originWorkdir)
				}()

				fileMap := make(map[string]*FileVertex)
				for _, eachStat := range stats {
					path := eachStat.Root.FuncPos.Path

					if cur, ok := fileMap[path]; ok {
						cur.AffectedReferenceIds = append(cur.AffectedReferenceIds, eachStat.VisitedIds()...)
					} else {
						totalLine, err := lineCounter(path)
						if err != nil {
							return err
						}
						fileMap[path] = &FileVertex{
							FileName:             path,
							AffectedLines:        len(lineMap[path]),
							TotalLines:           totalLine,
							AffectedFunctions:    len(funcGraph.GetFunctionsByFileLines(path, lineMap[path])),
							TotalFunctions:       len(funcGraph.GetFunctionsByFile(path)),
							AffectedReferenceIds: eachStat.VisitedIds(),
							TotalReferences:      funcGraph.FuncCount(),
						}
					}
				}
				fileList := make([]*FileVertex, 0, len(fileMap))
				for _, v := range fileMap {
					// calc
					v.AffectedLinePercent = float32(v.AffectedLines) / float32(v.TotalLines)

					m := make(map[string]struct{})
					for _, each := range v.AffectedReferenceIds {
						m[each] = struct{}{}
					}
					v.AffectedReferences = len(m)
					v.AffectedReferencePercent = float32(v.AffectedReferences) / float32(v.TotalReferences)
					v.AffectedFunctionPercent = float32(v.AffectedFunctions) / float32(v.TotalFunctions)

					fileList = append(fileList, v)
				}

				if outputCsv != "" {
					log.Infof("creating output csv: %v", outputCsv)
					csvFile, err := os.OpenFile(outputCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
					if err != nil {
						return err
					}
					defer csvFile.Close()
					if err := gocsv.MarshalFile(&fileList, csvFile); err != nil {
						return err
					}
				}
				if outputJson != "" {
					log.Infof("creating output json: %s", outputJson)
					contentBytes, err := json.Marshal(&fileList)
					if err != nil {
						return err
					}
					err = os.WriteFile(outputJson, contentBytes, os.ModePerm)
					if err != nil {
						return err
					}
				}
			}

			log.Infof("everything done.")
			return nil
		},
	}
	app.Commands = append(app.Commands, diffCmd)
}

type FileVertex struct {
	FileName                 string  `csv:"fileName" json:"fileName"`
	AffectedLinePercent      float32 `csv:"affectedLinePercent" json:"affectedLinePercent"`
	AffectedFunctionPercent  float32 `csv:"affectedFunctionPercent" json:"affectedFunctionPercent"`
	AffectedReferencePercent float32 `csv:"affectedReferencePercent" json:"affectedReferencePercent"`

	AffectedLines int `csv:"affectedLines" json:"affectedLines"`
	TotalLines    int `csv:"totalLines" json:"totalLines"`

	AffectedFunctions int `csv:"affectedFunctions" json:"affectedFunctions"`
	TotalFunctions    int `csv:"totalFunctions" json:"totalFunctions"`

	AffectedReferences   int      `csv:"affectedReferences" json:"affectedReferences"`
	AffectedReferenceIds []string `csv:"-" json:"-"`
	TotalReferences      int      `csv:"totalReferences" json:"totalReferences"`
}

// https://stackoverflow.com/a/24563853
func lineCounter(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount, nil
}
