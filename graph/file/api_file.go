package file

import (
	"path/filepath"

	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

type FileGraph struct {
	// reference graph (called graph)
	G graph.Graph[string, *FileVertex]
	// reverse reference graph (call graph)
	Rg graph.Graph[string, *FileVertex]
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
		G:  graph.New((*FileVertex).Id, graph.Directed()),
		Rg: graph.New((*FileVertex).Id, graph.Directed()),
	}

	// building edges
	err := fileGraph2FileGraph(fg.G, fileGraph.G)
	if err != nil {
		return nil, err
	}
	err = fileGraph2FileGraph(fg.Rg, fileGraph.Rg)
	if err != nil {
		return nil, err
	}

	nodeCount, err := fileGraph.G.Order()
	if err != nil {
		return nil, err
	}
	edgeCount, err := fileGraph.G.Size()
	if err != nil {
		return nil, err
	}
	log.Infof("dir graph ready. nodes: %d, edges: %d", nodeCount, edgeCount)

	return fileGraph, nil
}

func (fg *FileGraph) GetById(id string) *FileVertex {
	v, err := fg.G.Vertex(id)
	if err != nil {

		log.Warnf("no vertex: %v", id)
		return nil
	}
	return v
}

func path2dir(fp string) string {
	return filepath.ToSlash(filepath.Dir(fp))
}

func Path2vertex(fp string) *FileVertex {
	return &FileVertex{Path: fp}
}

func fileGraph2FileGraph(f graph.Graph[string, *FileVertex], g graph.Graph[string, *FileVertex]) error {
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
