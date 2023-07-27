package file

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

func (fg *Graph) DirectReferencedCount(f *Vertex) int {
	return len(fg.DirectReferencedIds(f))
}

func (fg *Graph) DirectReferencedIds(f *Vertex) []string {
	adjacencyMap, err := fg.G.AdjacencyMap()
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

func (fg *Graph) DirectReferenceIds(f *Vertex) []string {
	adjacencyMap, err := fg.Rg.AdjacencyMap()
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

func (fg *Graph) TransitiveReferencedIds(f *Vertex) []string {
	m := make(map[string]struct{}, 0)
	start := f.Id()
	graph.BFS(fg.G, start, func(cur string) bool {
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

func (fg *Graph) TransitiveReferenceIds(f *Vertex) []string {
	m := make(map[string]struct{}, 0)
	start := f.Id()
	graph.BFS(fg.Rg, start, func(cur string) bool {
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

func (fg *Graph) EntryIds(f *Vertex) []string {
	ret := make([]string, 0)
	// and also itself
	all := append(fg.TransitiveReferencedIds(f), f.Id())
	for _, eachId := range all {
		item := fg.IdCache[eachId]
		if item.ContainTag(TagEntry) {
			ret = append(ret, eachId)
		}
	}
	return ret
}
