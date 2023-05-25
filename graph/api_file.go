package graph

import (
	"path/filepath"

	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

type FileGraph struct {
	// reference graph (called graph)
	g graph.Graph[string, *FileVertex]
	// reverse reference graph (call graph)
	rg graph.Graph[string, *FileVertex]
}

type FileVertex struct {
	Path       string
	Referenced int
}

func (fv *FileVertex) Id() string {
	return fv.Path
}

func (fg *FileGraph) ToDirGraph() (*FileGraph, error) {
	// create graph
	fileGraph := &FileGraph{
		g:  graph.New((*FileVertex).Id, graph.Directed()),
		rg: graph.New((*FileVertex).Id, graph.Directed()),
	}

	// building edges
	err := fileGraph2FileGraph(fg.g, fileGraph.g)
	if err != nil {
		return nil, err
	}
	err = fileGraph2FileGraph(fg.rg, fileGraph.rg)
	if err != nil {
		return nil, err
	}

	nodeCount, err := fileGraph.g.Order()
	if err != nil {
		return nil, err
	}
	edgeCount, err := fileGraph.g.Size()
	if err != nil {
		return nil, err
	}
	log.Infof("dir graph ready. nodes: %d, edges: %d", nodeCount, edgeCount)

	return fileGraph, nil
}

func path2dir(fp string) string {
	return filepath.ToSlash(filepath.Dir(fp))
}

func path2vertex(fp string) *FileVertex {
	return &FileVertex{Path: fp}
}

func fileGraph2FileGraph(f graph.Graph[string, *FileVertex], g graph.Graph[string, *FileVertex]) error {
	m, err := f.PredecessorMap()
	if err != nil {
		return err
	}
	// add all the vertices
	for k := range m {
		_ = g.AddVertex(path2vertex(path2dir(k)))
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

func funcGraph2FileGraph(f graph.Graph[string, *FuncVertex], g graph.Graph[string, *FileVertex]) error {
	m, err := f.PredecessorMap()
	if err != nil {
		return err
	}
	// add all the vertices
	for k := range m {
		v, err := f.Vertex(k)
		if err != nil {
			return err
		}
		_ = g.AddVertex(path2vertex(v.Path))
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
		target, err := f.Vertex(eachEdge.Target)
		if err != nil {
			log.Warnf("vertex not found: %v", eachEdge.Target)
			continue
		}

		// ignore self ptr
		if source.Path == target.Path {
			continue
		}

		sourceFile, err := g.Vertex(source.Path)
		if err != nil {
			return err
		}
		targetFile, err := g.Vertex(target.Path)
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
