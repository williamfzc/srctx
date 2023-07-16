package function

import (
	"fmt"

	"github.com/williamfzc/srctx/graph/common"
)

func (fg *Graph) GetFunctionsByFile(fileName string) []*Vertex {
	if item, ok := fg.Cache[fileName]; ok {
		return item
	}
	return make([]*Vertex, 0)
}

func (fg *Graph) GetFunctionsByFileLines(fileName string, lines []int) []*Vertex {
	ret := make([]*Vertex, 0)
	functions := fg.Cache[fileName]
	if len(functions) == 0 {
		return ret
	}

	for _, eachFunc := range functions {
		// append these def lines
		if eachFunc.GetSpan().ContainAnyLine(lines...) {
			ret = append(ret, eachFunc)
		}
	}
	return ret
}

func (fg *Graph) GetById(id string) (*Vertex, error) {
	if item, ok := fg.IdCache[id]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("id not found in graph: %s", id)
}

func (fg *Graph) FuncCount() int {
	return len(fg.IdCache)
}

func (fg *Graph) ListFunctions() []*Vertex {
	return fg.FilterFunctions(func(funcVertex *Vertex) bool {
		return true
	})
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
	return fg.FilterFunctions(func(funcVertex *Vertex) bool {
		return funcVertex.ContainTag(TagEntry)
	})
}

func (fg *Graph) RelationBetween(a string, b string) (*common.EdgeStorage, error) {
	edge, err := fg.g.Edge(a, b)
	if err != nil {
		return nil, err
	}
	if ret, ok := edge.Properties.Data.(*common.EdgeStorage); ok {
		return ret, nil
	}
	return nil, fmt.Errorf("failed to convert %v", edge.Properties.Data)
}
