package graph

import (
	"os"

	"github.com/dominikbraun/graph/draw"
	log "github.com/sirupsen/logrus"
)

func (fg *FuncGraph) GetFunctionsByFile(f string) []*FuncVertex {
	return fg.cache[f]
}

func (fg *FuncGraph) GetById(id string) (*FuncVertex, error) {
	return fg.g.Vertex(id)
}

func (fg *FuncGraph) ReferencedCount(f *FuncVertex) int {
	return len(fg.DirectlyReferenced(f))
}

func (fg *FuncGraph) DirectlyReferenced(f *FuncVertex) []string {
	adjacencyMap, err := fg.g.AdjacencyMap()
	if err != nil {
		log.Warnf("failed to get adjacency map: %v", f)
		return nil
	}
	m := adjacencyMap[f.Id()]
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func (fg *FuncGraph) DirectlyReference(f *FuncVertex) []string {
	predecessorMap, err := fg.g.PredecessorMap()
	if err != nil {
		log.Warnf("failed to get predecessor map: %v", f)
		return nil
	}
	m := predecessorMap[f.Id()]
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func (fg *FuncGraph) DrawDot(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = draw.DOT(fg.g, file)
	if err != nil {
		return err
	}
	return nil
}
