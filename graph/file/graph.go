package file

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser"
)

type Graph struct {
	// reference graph (called graph)
	G graph.Graph[string, *Vertex]
	// reverse reference graph (call graph)
	Rg graph.Graph[string, *Vertex]
	// k: id, v: file
	IdCache map[string]*Vertex
}

type Vertex struct {
	Path       string
	Referenced int
}

func (fv *Vertex) Id() string {
	return fv.Path
}

func (fg *Graph) ToDirGraph() (*Graph, error) {
	// create graph
	fileGraph := &Graph{
		G:  graph.New((*Vertex).Id, graph.Directed()),
		Rg: graph.New((*Vertex).Id, graph.Directed()),
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

func NewEmptyFileGraph() *Graph {
	return &Graph{
		G:       graph.New((*Vertex).Id, graph.Directed()),
		Rg:      graph.New((*Vertex).Id, graph.Directed()),
		IdCache: make(map[string]*Vertex),
	}
}

func CreateFileGraphFromDirWithLSIF(src string, lsifFile string) (*Graph, error) {
	sourceContext, err := parser.FromLsifFile(lsifFile, src)
	if err != nil {
		return nil, err
	}
	log.Infof("fact ready. creating file graph ...")
	return CreateFileGraph(sourceContext)
}

func CreateFileGraph(relationship *object.SourceContext) (*Graph, error) {
	g := NewEmptyFileGraph()

	// nodes
	for each := range relationship.FileMapping {
		v := Path2vertex(each)
		err := g.G.AddVertex(v)
		if err != nil {
			return nil, err
		}

		err = g.Rg.AddVertex(v)
		if err != nil {
			return nil, err
		}

		g.IdCache[each] = v
	}

	for eachSrcFile := range relationship.FileMapping {
		refs, err := relationship.RefsByFileName(eachSrcFile)
		if err != nil {
			return nil, err
		}
		for _, eachRef := range refs {
			defs, err := relationship.RefsFromDefId(eachRef.Id())
			if err != nil {
				return nil, err
			}

			for _, eachDef := range defs {
				targetFile := relationship.FileName(eachDef.FileId)

				// avoid itself
				if eachSrcFile == targetFile {
					continue
				}
				_ = g.G.AddEdge(eachSrcFile, targetFile)
				_ = g.Rg.AddEdge(targetFile, eachSrcFile)
			}
		}
	}

	nodeCount, err := g.G.Order()
	if err != nil {
		return nil, err
	}
	edgeCount, err := g.G.Size()
	if err != nil {
		return nil, err
	}
	log.Infof("file graph ready. nodes: %d, edges: %d", nodeCount, edgeCount)

	return g, nil
}
