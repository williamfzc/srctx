package file

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/graph/common"
)

func (fg *Graph) GetById(id string) *Vertex {
	v, err := fg.G.Vertex(id)
	if err != nil {

		logrus.Warnf("no vertex: %v", id)
		return nil
	}
	return v
}

func (fg *Graph) ListFiles() []*Vertex {
	ret := make([]*Vertex, 0, len(fg.IdCache))
	for _, each := range fg.IdCache {
		ret = append(ret, each)
	}
	return ret
}

func (fg *Graph) FilterFunctions(f func(*Vertex) bool) []*Vertex {
	ret := make([]*Vertex, 0)
	for _, each := range fg.IdCache {
		if f(each) {
			ret = append(ret, each)
		}
	}
	return ret
}

func (fg *Graph) ListEntries() []*Vertex {
	return fg.FilterFunctions(func(vertex *Vertex) bool {
		return vertex.ContainTag(TagEntry)
	})
}

func (fg *Graph) RelationBetween(a string, b string) (*common.EdgeStorage, error) {
	edge, err := fg.G.Edge(a, b)
	if err != nil {
		return nil, err
	}
	if ret, ok := edge.Properties.Data.(*common.EdgeStorage); ok {
		return ret, nil
	}
	return nil, fmt.Errorf("failed to convert %v", edge.Properties.Data)
}
