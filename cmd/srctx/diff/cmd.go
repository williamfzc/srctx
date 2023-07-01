package diff

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/williamfzc/srctx/graph/function"

	"github.com/williamfzc/srctx/graph/visual/g6"

	"github.com/gocarina/gocsv"
	"github.com/opensibyl/sibyl2/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/parser"
	"github.com/williamfzc/srctx/parser/lsif"
)

const (
	nodeLevelFunc = "func"
	nodeLevelFile = "file"
	nodeLevelDir  = "dir"
)

type Options struct {
	// required
	Src      string `json:"src,omitempty"`
	RepoRoot string `json:"repoRoot,omitempty"`
	Before   string `json:"before,omitempty"`
	After    string `json:"after,omitempty"`
	LsifZip  string `json:"lsifZip,omitempty"`
	ScipFile string `json:"scipFile,omitempty"`

	// output
	OutputJson string `json:"outputJson,omitempty"`
	OutputCsv  string `json:"outputCsv,omitempty"`
	OutputDot  string `json:"outputDot,omitempty"`
	OutputHtml string `json:"outputHtml,omitempty"`

	// options
	NodeLevel string `json:"nodeLevel,omitempty"`
	WithIndex bool   `json:"withIndex,omitempty"`
	CacheType string `json:"cacheType,omitempty"`
	Lang      string `json:"lang,omitempty"`
	NoDiff    bool   `json:"noDiff,omitempty"`
}

func NewOptionsFromCliFlags(c *cli.Context) Options {
	return Options{
		Src:        c.String("src"),
		RepoRoot:   c.String("repoRoot"),
		Before:     c.String("before"),
		After:      c.String("after"),
		LsifZip:    c.String("lsif"),
		ScipFile:   c.String("scip"),
		OutputJson: c.String("outputJson"),
		OutputCsv:  c.String("outputCsv"),
		OutputDot:  c.String("outputDot"),
		OutputHtml: c.String("outputHtml"),
		NodeLevel:  c.String("nodeLevel"),
		WithIndex:  c.Bool("withIndex"),
		CacheType:  c.String("cacheType"),
		Lang:       c.String("lang"),
		NoDiff:     c.Bool("noDiff"),
	}
}

func AddDiffCmd(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "src",
			Value: ".",
			Usage: "project path",
		},
		&cli.StringFlag{
			Name:  "repoRoot",
			Value: "",
			Usage: "root path of your repo",
		},
		&cli.StringFlag{
			Name:  "before",
			Value: "HEAD~1",
			Usage: "before rev",
		},
		&cli.StringFlag{
			Name:  "after",
			Value: "HEAD",
			Usage: "after rev",
		},
		&cli.StringFlag{
			Name:  "lsif",
			Value: "./dump.lsif",
			Usage: "lsif path, can be zip or origin file",
		},
		&cli.StringFlag{
			Name:  "scip",
			Value: "",
			Usage: "scip file",
		},
		&cli.StringFlag{
			Name:  "nodeLevel",
			Value: nodeLevelFunc,
			Usage: "graph level (file or func or dir)",
		},
		&cli.StringFlag{
			Name:  "outputJson",
			Value: "",
			Usage: "json output",
		},
		&cli.StringFlag{
			Name:  "outputCsv",
			Value: "",
			Usage: "csv output",
		},
		&cli.StringFlag{
			Name:  "outputDot",
			Value: "",
			Usage: "reference dot file output",
		},
		&cli.StringFlag{
			Name:  "outputHtml",
			Value: "",
			Usage: "render html report with g6",
		},
		&cli.BoolFlag{
			Name:  "withIndex",
			Value: false,
			Usage: "create indexes first if enabled, currently support golang only",
		},
		&cli.StringFlag{
			Name:  "cacheType",
			Value: lsif.CacheTypeFile,
			Usage: "mem or file",
		},
		&cli.StringFlag{
			Name:  "lang",
			Value: string(core.LangUnknown),
			Usage: "language of repo",
		},
		&cli.BoolFlag{
			Name:  "noDiff",
			Value: false,
			Usage: "will not calc git diff if enabled",
		},
	}

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
			log.Infof("start diffing: %v", src)

			if opts.CacheType != lsif.CacheTypeFile {
				parser.UseMemCache()
			}

			lang := core.LangType(opts.Lang)

			var lineMap diff.AffectedLineMap
			totalLineCountMap := make(map[string]int)
			if !opts.NoDiff {
				// prepare
				lineMap, err = diff.GitDiff(src, opts.Before, opts.After)
				if err != nil {
					return err
				}

				// Maybe the src does not start from the root of repo.
				// But the lineMap will start from the root of repo.
				if opts.RepoRoot != "" {
					repoRoot, err := filepath.Abs(opts.RepoRoot)
					if err != nil {
						return err
					}

					log.Infof("path sync from %s to %s", repoRoot, src)
					lineMap, err = diff.PathOffset(repoRoot, src, lineMap)
					if err != nil {
						return err
					}

					for eachPath := range lineMap {
						totalLineCountMap[eachPath], err = lineCounter(filepath.Join(src, eachPath))
						if err != nil {
							return err
						}
					}
				}
			} else {
				log.Infof("noDiff enabled")
				lineMap = make(diff.AffectedLineMap)
			}

			// metadata
			var funcGraph *function.FuncGraph

			if opts.ScipFile != "" {
				// using SCIP
				log.Infof("using SCIP as index")
				funcGraph, err = function.CreateFuncGraphFromDirWithSCIP(src, opts.ScipFile, lang)
			} else {
				// using LSIF
				log.Infof("using LSIF as index")
				if opts.WithIndex {
					switch lang {
					case core.LangGo:
						funcGraph, err = function.CreateFuncGraphFromGolangDir(src)
					default:
						return errors.New("did not specify `--lang`")
					}
				} else {
					funcGraph, err = function.CreateFuncGraphFromDirWithLSIF(src, opts.LsifZip, lang)
				}
			}

			if err != nil {
				return err
			}

			// look up start points
			startPoints := make([]*function.FuncVertex, 0)
			for path, lines := range lineMap {
				curPoints := funcGraph.GetFunctionsByFileLines(path, lines)
				if len(curPoints) == 0 {
					log.Infof("file %s line %v hit no func", path, lines)
				} else {
					startPoints = append(startPoints, curPoints...)
				}
			}

			// start scan
			stats := make([]*function.VertexStat, 0)
			for _, eachPtr := range startPoints {
				eachStat := funcGraph.Stat(eachPtr)
				stats = append(stats, eachStat)
				log.Infof("start point: %v, refed: %d, ref: %d", eachPtr.Id(), eachStat.Referenced, eachStat.Reference)
			}
			log.Infof("diff finished.")

			// tag (with priority)
			for _, eachStat := range stats {
				for _, eachVisited := range eachStat.VisitedIds() {
					err := funcGraph.FillWithYellow(eachVisited)
					if err != nil {
						return err
					}
				}
			}
			for _, eachStat := range stats {
				for _, eachId := range eachStat.ReferencedIds {
					err := funcGraph.FillWithOrange(eachId)
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

			// output
			if opts.OutputDot != "" {
				log.Infof("creating dot file: %v", opts.OutputDot)

				switch opts.NodeLevel {
				case nodeLevelFunc:
					err := funcGraph.DrawDot(opts.OutputDot)
					if err != nil {
						return err
					}
				case nodeLevelFile:
					fileGraph, err := funcGraph.ToFileGraph()
					if err != nil {
						return err
					}
					err = fileGraph.DrawDot(opts.OutputDot)
					if err != nil {
						return err
					}
				case nodeLevelDir:
					dirGraph, err := funcGraph.ToDirGraph()
					if err != nil {
						return err
					}
					err = dirGraph.DrawDot(opts.OutputDot)
					if err != nil {
						return err
					}
				}
			}

			if opts.OutputCsv != "" || opts.OutputJson != "" {
				fileMap := make(map[string]*FileVertex)
				for _, eachStat := range stats {
					path := eachStat.Root.FuncPos.Path

					if cur, ok := fileMap[path]; ok {
						cur.AffectedReferenceIds = append(cur.AffectedReferenceIds, eachStat.VisitedIds()...)
					} else {
						totalLine := totalLineCountMap[path]
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
					v.AffectedLinePercent = 0.0
					if v.TotalLines != 0 {
						v.AffectedLinePercent = float32(v.AffectedLines) / float32(v.TotalLines)
					}

					m := make(map[string]struct{})
					for _, each := range v.AffectedReferenceIds {
						m[each] = struct{}{}
					}
					v.AffectedReferences = len(m)
					v.AffectedReferencePercent = 0.0
					if v.TotalReferences != 0 {
						v.AffectedReferencePercent = float32(v.AffectedReferences) / float32(v.TotalReferences)
					}
					v.AffectedFunctionPercent = 0.0
					if v.TotalFunctions != 0 {
						v.AffectedFunctionPercent = float32(v.AffectedFunctions) / float32(v.TotalFunctions)
					}

					fileList = append(fileList, v)
				}

				if opts.OutputCsv != "" {
					log.Infof("creating output csv: %v", opts.OutputCsv)
					csvFile, err := os.OpenFile(opts.OutputCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
					if err != nil {
						return err
					}
					defer csvFile.Close()
					if err := gocsv.MarshalFile(&fileList, csvFile); err != nil {
						return err
					}
				}

				if opts.OutputJson != "" {
					log.Infof("creating output json: %s", opts.OutputJson)
					contentBytes, err := json.Marshal(&fileList)
					if err != nil {
						return err
					}
					err = os.WriteFile(opts.OutputJson, contentBytes, os.ModePerm)
					if err != nil {
						return err
					}
				}
			}

			if opts.OutputHtml != "" {
				log.Infof("creating output html: %s", opts.OutputHtml)

				var g6data *g6.Data
				if opts.NodeLevel != nodeLevelFunc {
					fileGraph, err := funcGraph.ToFileGraph()
					if err != nil {
						return err
					}
					g6data, err = fileGraph.ToG6Data()
					if err != nil {
						return err
					}

					for _, eachStat := range stats {
						for _, eachId := range eachStat.TransitiveReferencedIds {
							f, err := funcGraph.GetById(eachId)
							if err != nil {
								return err
							}
							g6data.FillWithYellow(f.Path)
						}
					}
					for _, eachStat := range stats {
						g6data.FillWithRed(eachStat.Root.Path)
					}
				} else {
					g6data, err = funcGraph.ToG6Data()
					if err != nil {
						return err
					}
				}

				err = g6data.RenderHtml(opts.OutputHtml)
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
