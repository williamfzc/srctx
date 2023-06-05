package diff

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/opensibyl/sibyl2/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/graph"
	"github.com/williamfzc/srctx/parser"
	"github.com/williamfzc/srctx/parser/lsif"
)

const (
	nodeLevelFunc = "func"
	nodeLevelFile = "file"
	nodeLevelDir  = "dir"
)

func AddDiffCmd(app *cli.App) {
	// required
	var src string
	var repoRoot string
	var before string
	var after string
	var lsifZip string
	var scipFile string

	// output
	var outputJson string
	var outputCsv string
	var outputDot string
	var outputHtml string
	var nodeLevel string

	// options
	var withIndex bool
	var cacheType string
	var lang string
	var noDiff bool

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "src",
			Value:       ".",
			Usage:       "project path",
			Destination: &src,
		},
		&cli.StringFlag{
			Name:        "repoRoot",
			Value:       "",
			Usage:       "root path of your repo",
			Destination: &repoRoot,
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
			Name:        "scip",
			Value:       "",
			Usage:       "scip file",
			Destination: &scipFile,
		},
		&cli.StringFlag{
			Name:        "nodeLevel",
			Value:       nodeLevelFunc,
			Usage:       "graph level (file or func or dir)",
			Destination: &nodeLevel,
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
			Name:        "outputHtml",
			Value:       "",
			Usage:       "render html report with g6",
			Destination: &outputHtml,
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
		&cli.BoolFlag{
			Name:        "noDiff",
			Value:       false,
			Usage:       "will not calc git diff if enabled",
			Destination: &noDiff,
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

			if cacheType != lsif.CacheTypeFile {
				parser.UseMemCache()
			}

			lang := core.LangType(lang)

			var lineMap diff.AffectedLineMap
			if !noDiff {
				// prepare
				lineMap, err := diff.GitDiff(src, before, after)
				if err != nil {
					return err
				}

				// git diff will always start from the root of repo
				// src will always start from the root of project
				if repoRoot != "" {
					repoRoot, err := filepath.Abs(repoRoot)
					if err != nil {
						return err
					}

					log.Infof("path sync from %s to %s", repoRoot, src)
					modifiedLineMap := make(map[string][]int)
					for file, lines := range lineMap {
						absFile := filepath.Join(repoRoot, file)
						relPath, err := filepath.Rel(src, absFile)
						if err != nil {
							return err
						}
						modifiedLineMap[relPath] = lines
					}
					lineMap = modifiedLineMap
				}
			} else {
				log.Infof("noDiff enabled")
				lineMap = make(diff.AffectedLineMap)
			}

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
					switch lang {
					case core.LangGo:
						funcGraph, err = graph.CreateFuncGraphFromGolangDir(src, lang)
					default:
						return errors.New("did not specify `--lang`")
					}
				} else {
					funcGraph, err = graph.CreateFuncGraphFromDirWithLSIF(src, lsifZip, lang)
				}
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

				switch nodeLevel {
				case nodeLevelFunc:
					// colorful
					for _, eachStat := range stats {
						for _, eachVisited := range eachStat.VisitedIds() {
							err := funcGraph.FillWithYellow(eachVisited)
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
				case nodeLevelFile:
					fileGraph, err := funcGraph.ToFileGraph()
					if err != nil {
						return err
					}
					err = fileGraph.DrawDot(outputDot)
					if err != nil {
						return err
					}
				case nodeLevelDir:
					dirGraph, err := funcGraph.ToDirGraph()
					if err != nil {
						return err
					}
					err = dirGraph.DrawDot(outputDot)
					if err != nil {
						return err
					}
				}
			}

			if outputCsv != "" || outputJson != "" {
				// need to access files
				// todo: should be removed
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

			if outputHtml != "" {
				log.Infof("createing output html: %s", outputHtml)

				var g6data *graph.G6Data
				if nodeLevel != nodeLevelFunc {
					fileGraph, err := funcGraph.ToFileGraph()
					if err != nil {
						return err
					}
					g6data, err = fileGraph.ToG6Data()
				} else {
					g6data, err = funcGraph.ToG6Data()
				}
				if err != nil {
					return err
				}

				// todo: colorful
				err = g6data.RenderHtml(outputHtml)
				if err != nil {
					return err
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
