package diff

import (
	"encoding/json"
	"errors"
	"os"

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
		pv := file.Path2vertex(path)
		startPoints = append(startPoints, pv)
	}
	// start scan
	stats := make([]*ImpactUnitWithFile, 0)
	for _, eachPtr := range startPoints {
		eachStat := fileGraph.Stat(eachPtr)
		wrappedStat := WrapImpactUnitWithFile(eachStat)

		// fill with file info
		if totalLineCount, ok := totalLineCountMap[eachStat.FileName]; ok {
			wrappedStat.TotalLineCount = totalLineCount
		}
		if impactLineCount, ok := lineMap[eachStat.FileName]; ok {
			wrappedStat.ImpactLineCount = len(impactLineCount)
		}
		stats = append(stats, wrappedStat)
		log.Infof("start point: %v, refed: %d, ref: %d", eachPtr.Id(), len(eachStat.ReferencedIds), len(eachStat.ReferenceIds))
	}
	log.Infof("diff finished.")

	// tag
	for _, eachStartPoint := range startPoints {
		err = fileGraph.FillWithRed(eachStartPoint.Id())
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
	return nil
}

func createFileGraph(opts *Options) (*file.Graph, error) {
	var fileGraph *file.Graph
	var err error

	if opts.ScipFile != "" {
		// using SCIP
		log.Infof("using SCIP as index")
		fileGraph, err = file.CreateFileGraphFromDirWithSCIP(opts.Src, opts.ScipFile)
	} else {
		// using LSIF
		log.Infof("using LSIF as index")
		if opts.WithIndex {
			switch core.LangType(opts.Lang) {
			case core.LangGo:
				fileGraph, err = file.CreateFileGraphFromGolangDir(opts.Src)
			default:
				return nil, errors.New("did not specify `--lang`")
			}
		} else {
			fileGraph, err = file.CreateFileGraphFromDirWithLSIF(opts.Src, opts.LsifZip)
		}
	}
	if err != nil {
		return nil, err
	}
	return fileGraph, nil
}
