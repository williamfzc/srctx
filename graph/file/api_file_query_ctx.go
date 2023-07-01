package file

import log "github.com/sirupsen/logrus"

func (fg *FileGraph) DirectReferencedCount(f *FileVertex) int {
	return len(fg.DirectReferencedIds(f))
}

func (fg *FileGraph) DirectReferencedIds(f *FileVertex) []string {
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

func (fg *FileGraph) DirectReferenceIds(f *FileVertex) []string {
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
