package diff

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/opensibyl/sibyl2/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/graph/function"
	"github.com/williamfzc/srctx/graph/visual/g6"
)

func funcLevelMain(opts *Options, lineMap diff.AffectedLineMap, totalLineCountMap map[string]int) error {
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

		err := funcGraph.DrawDot(opts.OutputDot)
		if err != nil {
			return err
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

		err = g6data.RenderHtml(opts.OutputHtml)
		if err != nil {
			return err
		}
	}
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
