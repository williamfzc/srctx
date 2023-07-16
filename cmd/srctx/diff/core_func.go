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
	log.Infof("func level main entry")
	// metadata
	funcGraph, err := createFuncGraph(opts)
	if err != nil {
		return err
	}

	// look up start points
	startPoints := make([]*function.Vertex, 0)
	for path, lines := range lineMap {
		curPoints := funcGraph.GetFunctionsByFileLines(path, lines)
		if len(curPoints) == 0 {
			log.Infof("file %s line %v hit no func", path, lines)
		} else {
			startPoints = append(startPoints, curPoints...)
		}
	}

	// start scan
	stats := make([]*ImpactUnitWithFile, 0)
	for _, eachPtr := range startPoints {
		eachStat := funcGraph.Stat(eachPtr)
		wrappedStat := WrapImpactUnitWithFile(eachStat)

		totalLineCount := len(eachPtr.GetSpan().Lines())
		affectedLineCount := 0

		if lines, ok := lineMap[eachPtr.Path]; ok {
			for _, eachLine := range lines {
				if eachPtr.GetSpan().ContainLine(eachLine) {
					affectedLineCount++
				}
			}
		}
		wrappedStat.TotalLineCount = totalLineCount
		wrappedStat.AffectedLineCount = affectedLineCount

		stats = append(stats, wrappedStat)
		log.Infof("start point: %v, refed: %d, ref: %d", eachPtr.Id(), len(eachStat.ReferencedIds), len(eachStat.ReferenceIds))
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
		err := funcGraph.FillWithRed(eachStat.UnitName)
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
		if opts.OutputCsv != "" {
			log.Infof("creating output csv: %v", opts.OutputCsv)
			csvFile, err := os.OpenFile(opts.OutputCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer csvFile.Close()
			if err := gocsv.MarshalFile(&stats, csvFile); err != nil {
				return err
			}
		}

		if opts.OutputJson != "" {
			log.Infof("creating output json: %s", opts.OutputJson)
			contentBytes, err := json.Marshal(&stats)
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
		g6data, err = funcGraph.ToG6Data()
		if err != nil {
			return err
		}

		err = g6data.RenderHtml(opts.OutputHtml)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFuncGraph(opts *Options) (*function.Graph, error) {
	var funcGraph *function.Graph
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
