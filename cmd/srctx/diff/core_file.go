package diff

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/williamfzc/srctx/graph/common"

	"github.com/gocarina/gocsv"
	"github.com/opensibyl/sibyl2/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/graph/file"
)

func fileLevelMain(opts *Options, lineMap diff.ImpactLineMap, totalLineCountMap map[string]int) error {
	log.Infof("file level main entry")
	fileGraph, err := createFileGraph(opts)
	if err != nil {
		return err
	}

	// look up start points
	startPoints := make([]*file.Vertex, 0)
	for path := range lineMap {
		// this vertex may not exist in graph! like: push.yml, README.md
		pv := fileGraph.GetById(path)
		if pv != nil {
			startPoints = append(startPoints, pv)
		}
	}
	// start scan
	stats := make([]*ImpactUnitWithFile, 0)
	globalStat := fileGraph.GlobalStat(startPoints)

	for _, eachStat := range globalStat.ImpactUnitsMap {
		wrappedStat := WrapImpactUnitWithFile(eachStat)

		// fill with file info
		if totalLineCount, ok := totalLineCountMap[eachStat.FileName]; ok {
			wrappedStat.TotalLineCount = totalLineCount
		}
		if impactLineCount, ok := lineMap[eachStat.FileName]; ok {
			wrappedStat.ImpactLineCount = len(impactLineCount)
		}
		stats = append(stats, wrappedStat)
		log.Infof("start point: %v, refed: %d, ref: %d", eachStat.UnitName, len(eachStat.ReferencedIds), len(eachStat.ReferenceIds))
	}
	log.Infof("diff finished.")

	// tag
	for _, eachStartPoint := range startPoints {
		err = fileGraph.FillWithRed(eachStartPoint.Id())
		if err != nil {
			return err
		}
	}

	if opts.OutputDot != "" {
		log.Infof("creating dot file: %v", opts.OutputDot)

		err := fileGraph.DrawDot(opts.OutputDot)
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

		g6data, err := fileGraph.ToG6Data()
		if err != nil {
			return err
		}
		err = g6data.RenderHtml(opts.OutputHtml)
		if err != nil {
			return err
		}
	}

	if opts.StatJson != "" {
		log.Infof("creating stat json: %s", opts.StatJson)
		contentBytes, err := json.Marshal(&globalStat)
		if err != nil {
			return err
		}
		err = os.WriteFile(opts.StatJson, contentBytes, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func createFileGraph(opts *Options) (*file.Graph, error) {
	var fileGraph *file.Graph
	var err error

	graphOptions := common.DefaultGraphOptions()
	graphOptions.Src = opts.Src
	graphOptions.NoEntries = opts.NoEntries

	if opts.ScipFile != "" {
		// using SCIP
		log.Infof("using SCIP as index")
		graphOptions.ScipFile = opts.ScipFile
		fileGraph, err = file.CreateFileGraphFromDirWithSCIP(graphOptions)
	} else {
		// using LSIF
		log.Infof("using LSIF as index")
		if opts.WithIndex {
			switch core.LangType(opts.Lang) {
			case core.LangGo:
				graphOptions.GenGolangIndex = true
				fileGraph, err = file.CreateFileGraphFromGolangDir(graphOptions)
			default:
				return nil, errors.New("did not specify `--lang`")
			}
		} else {
			graphOptions.LsifFile = opts.LsifZip
			fileGraph, err = file.CreateFileGraphFromDirWithLSIF(graphOptions)
		}
	}
	if err != nil {
		return nil, err
	}
	return fileGraph, nil
}
