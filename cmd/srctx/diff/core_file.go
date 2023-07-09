package diff

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/graph/file"
)

func fileLevelMain(opts *Options, lineMap diff.AffectedLineMap) error {
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

	// tag
	for _, eachStartPoint := range startPoints {
		err = fileGraph.FillWithRed(eachStartPoint.Id())
		if err != nil {
			return err
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
