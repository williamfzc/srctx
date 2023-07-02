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
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/graph/function"
	"github.com/williamfzc/srctx/graph/visual/g6"
	"github.com/williamfzc/srctx/parser"
	"github.com/williamfzc/srctx/parser/lsif"
)

// MainDiff allow accessing as a lib
func MainDiff(opts *Options) error {
	log.Infof("start diffing: %v", opts.Src)

	if opts.CacheType != lsif.CacheTypeFile {
		parser.UseMemCache()
	}

	// collect diff info
	lineMap, err := collectLineMap(opts)
	if err != nil {
		return err
	}

	// collect info from file (line number/size ...)
	totalLineCountMap, err := collectTotalLineCountMap(opts, opts.Src, lineMap)
	if err != nil {
		return err
	}

	// metadata
	funcGraph, err := createFuncGraph(opts)
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
}

func createFuncGraph(opts *Options) (*function.FuncGraph, error) {
	var funcGraph *function.FuncGraph
	var err error

	if opts.ScipFile != "" {
		// using SCIP
		log.Infof("using SCIP as index")
		funcGraph, err = function.CreateFuncGraphFromDirWithSCIP(opts.Src, opts.ScipFile, core.LangType(opts.Lang))
	} else {
		// using LSIF
		log.Infof("using LSIF as index")
		if opts.WithIndex {
			switch core.LangType(opts.Lang) {
			case core.LangGo:
				funcGraph, err = function.CreateFuncGraphFromGolangDir(opts.Src)
			default:
				return nil, errors.New("did not specify `--lang`")
			}
		} else {
			funcGraph, err = function.CreateFuncGraphFromDirWithLSIF(opts.Src, opts.LsifZip, core.LangType(opts.Lang))
		}
	}
	if err != nil {
		return nil, err
	}
	return funcGraph, nil
}

func collectLineMap(opts *Options) (diff.AffectedLineMap, error) {
	if !opts.NoDiff {
		lineMap, err := diff.GitDiff(opts.Src, opts.Before, opts.After)
		if err != nil {
			return nil, err
		}
		return lineMap, nil
	}
	log.Infof("noDiff enabled")
	return make(diff.AffectedLineMap), nil
}

func collectTotalLineCountMap(opts *Options, src string, lineMap diff.AffectedLineMap) (map[string]int, error) {
	totalLineCountMap := make(map[string]int)

	if opts.RepoRoot != "" {
		repoRoot, err := filepath.Abs(opts.RepoRoot)
		if err != nil {
			return nil, err
		}

		log.Infof("path sync from %s to %s", repoRoot, src)
		lineMap, err = diff.PathOffset(repoRoot, src, lineMap)
		if err != nil {
			return nil, err
		}

		for eachPath := range lineMap {
			totalLineCountMap[eachPath], err = lineCounter(filepath.Join(src, eachPath))
			if err != nil {
				return nil, err
			}
		}
	}

	return totalLineCountMap, nil
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
