package graph

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

func (fg *FuncGraph) FuncCount() int {
	ret, _ := fg.g.Order()
	return ret
}

func (fg *FuncGraph) ReferencedCount(f *FuncVertex) int {
	return len(fg.ReferencedIds(f))
}

func (fg *FuncGraph) ReferencedIds(f *FuncVertex) []string {
	adjacencyMap, err := fg.rg.AdjacencyMap()
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

func (fg *FuncGraph) TransitiveReferencedIds(f *FuncVertex) []string {
	m := make(map[string]struct{}, 0)
	start := f.Id()
	graph.BFS(fg.g, start, func(cur string) bool {
		if cur == start {
			return false
		}
		m[cur] = struct{}{}
		return false
	})
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func (fg *FuncGraph) ReferenceIds(f *FuncVertex) []string {
	predecessorMap, err := fg.rg.PredecessorMap()
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

func (fg *FuncGraph) TransitiveReferenceIds(f *FuncVertex) []string {
	m := make(map[string]struct{}, 0)
	start := f.Id()
	graph.BFS(fg.rg, start, func(cur string) bool {
		if cur == start {
			return false
		}
		m[cur] = struct{}{}
		return false
	})
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}
