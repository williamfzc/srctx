package file

import (
	"path/filepath"

	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

func path2dir(fp string) string {
	return filepath.ToSlash(filepath.Dir(fp))
}

func Path2vertex(fp string) *Vertex {
	return &Vertex{Path: fp}
}

func fileGraph2FileGraph(f graph.Graph[string, *Vertex], g graph.Graph[string, *Vertex]) error {
	m, err := f.AdjacencyMap()
	if err != nil {
		return err
	}
	// add all the vertices
	for k := range m {
		_ = g.AddVertex(Path2vertex(path2dir(k)))
	}

	edges, err := f.Edges()
	if err != nil {
		return err
	}
	for _, eachEdge := range edges {
		source, err := f.Vertex(eachEdge.Source)
		if err != nil {
			log.Warnf("vertex not found: %v", eachEdge.Source)
			continue
		}
		dirSourcePath := path2dir(source.Path)

		target, err := f.Vertex(eachEdge.Target)
		if err != nil {
			log.Warnf("vertex not found: %v", eachEdge.Target)
			continue
		}
		dirTargetPath := path2dir(target.Path)

		// ignore self ptr
		if dirSourcePath == dirTargetPath {
			continue
		}

		sourceFile, err := g.Vertex(dirSourcePath)
		if err != nil {
			return err
		}
		targetFile, err := g.Vertex(dirTargetPath)
		if err != nil {
			return err
		}
		if sv, err := g.Vertex(sourceFile.Id()); err == nil {
			sv.Referenced++
		} else {
			_ = g.AddVertex(sourceFile)
		}
		if tv, err := g.Vertex(targetFile.Id()); err == nil {
			tv.Referenced++
		} else {
			_ = g.AddVertex(targetFile)
		}

		_ = g.AddEdge(sourceFile.Id(), targetFile.Id())
	}

	return nil
}
