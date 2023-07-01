package file

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser"
)

func NewEmptyFileGraph() *FileGraph {
	return &FileGraph{
		G:  graph.New((*FileVertex).Id, graph.Directed()),
		Rg: graph.New((*FileVertex).Id, graph.Directed()),
	}
}

func CreateFileGraphFromDirWithLSIF(src string, lsifFile string) (*FileGraph, error) {
	sourceContext, err := parser.FromLsifFile(lsifFile, src)
	if err != nil {
		return nil, err
	}
	log.Infof("fact ready. creating file graph ...")
	return CreateFileGraph(sourceContext)
}

func CreateFileGraph(relationship *object.SourceContext) (*FileGraph, error) {
	g := NewEmptyFileGraph()

	// nodes
	for each := range relationship.FileMapping {
		err := g.G.AddVertex(&FileVertex{
			Path:       each,
			Referenced: 0,
		})
		if err != nil {
			return nil, err
		}

		err = g.Rg.AddVertex(&FileVertex{
			Path:       each,
			Referenced: 0,
		})
		if err != nil {
			return nil, err
		}
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
