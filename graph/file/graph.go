package file

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/graph/common"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser"
)

const TagEntry = "entry"

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

	// https://github.com/williamfzc/srctx/issues/41
	Tags map[string]struct{} `json:"tags,omitempty"`
}

func (fv *Vertex) Id() string {
	return fv.Path
}

func (fv *Vertex) ContainTag(tag string) bool {
	if _, ok := fv.Tags[tag]; ok {
		return true
	}
	return false
}

func (fv *Vertex) AddTag(tag string) {
	fv.Tags[tag] = struct{}{}
}

func (fv *Vertex) RemoveTag(tag string) {
	delete(fv.Tags, tag)
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

func CreateFileGraphFromDirWithLSIF(opts *common.GraphOptions) (*Graph, error) {
	sourceContext, err := parser.FromLsifFile(opts.LsifFile, opts.Src)
	if err != nil {
		return nil, err
	}
	log.Infof("fact ready. creating file graph ...")
	return CreateFileGraph(sourceContext, opts)
}

func CreateFileGraphFromGolangDir(opts *common.GraphOptions) (*Graph, error) {
	sourceContext, err := parser.FromGolangSrc(opts.Src)
	if err != nil {
		return nil, err
	}
	return CreateFileGraph(sourceContext, opts)
}

func CreateFileGraphFromDirWithSCIP(opts *common.GraphOptions) (*Graph, error) {
	sourceContext, err := parser.FromScipFile(opts.ScipFile, opts.Src)
	if err != nil {
		return nil, err
	}
	return CreateFileGraph(sourceContext, opts)
}

func CreateFileGraph(relationship *object.SourceContext, opts *common.GraphOptions) (*Graph, error) {
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

	current := 0
	total := len(relationship.FileMapping)
	for eachSrcFile := range relationship.FileMapping {
		log.Infof("processing: %d / %d", current, total)
		current++

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

				refLineNumber := eachRef.LineNumber()
				if edge, err := g.G.Edge(eachSrcFile, targetFile); err == nil {
					storage := edge.Properties.Data.(*common.EdgeStorage)
					storage.RefLines[refLineNumber] = struct{}{}
				} else {
					_ = g.G.AddEdge(eachSrcFile, targetFile, graph.EdgeData(common.NewEdgeStorage()))
				}

				if edge, err := g.Rg.Edge(targetFile, eachSrcFile); err == nil {
					storage := edge.Properties.Data.(*common.EdgeStorage)
					storage.RefLines[refLineNumber] = struct{}{}
				} else {
					_ = g.Rg.AddEdge(targetFile, eachSrcFile, graph.EdgeData(common.NewEdgeStorage()))
				}
			}
		}
	}

	// entries tag
	// calculation takes a few minutes if target project is large
	if !opts.NoEntries {
		entries := g.FilterFunctions(func(vertex *Vertex) bool {
			return len(g.DirectReferencedIds(vertex)) == 0
		})
		log.Infof("detect entries: %d", len(entries))
		for _, entry := range entries {
			entry.AddTag(TagEntry)
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
