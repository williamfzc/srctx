package file

import log "github.com/sirupsen/logrus"

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
