package function

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

// DirectReferencedCount
// This function returns the number of direct references to a given function vertex in the function graph.
// It does so by counting the length of the slice of IDs of the function vertices that directly reference the given function vertex.
func (fg *FuncGraph) DirectReferencedCount(f *FuncVertex) int {
	return len(fg.DirectReferencedIds(f))
}

func (fg *FuncGraph) DirectReferencedIds(f *FuncVertex) []string {
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

func (fg *FuncGraph) DirectReferenceIds(f *FuncVertex) []string {
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

// TransitiveReferencedIds
// This function takes a FuncGraph and a FuncVertex as input and returns a slice of strings containing all the transitive referenced ids.
// It uses a map to store the referenced ids and a BFS algorithm to traverse the graph and add the referenced ids to the map.
// Finally, it returns the keys of the map as a slice of strings.
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

func (fg *FuncGraph) EntryIds(f *FuncVertex) []string {
	ret := make([]string, 0)
	all := fg.TransitiveReferencedIds(f)
	for _, eachId := range all {
		item := fg.IdCache[eachId]
		if item.ContainTag(TagEntry) {
			ret = append(ret, eachId)
		}
	}
	return ret
}
