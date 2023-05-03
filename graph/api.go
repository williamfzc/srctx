package graph

import (
	"github.com/dominikbraun/graph"
	log "github.com/sirupsen/logrus"
)

func (fg *FuncGraph) InfluenceCount(f *FuncVertex) int {
	ret := 0
	startPoint := f.Id()
	_ = graph.BFS(fg.g, startPoint, func(s string) bool {
		if startPoint == s {
			return false
		}
		if _, err := fg.g.Edge(startPoint, s); err != nil {
			return true
		}

		log.Infof("direct ref: %s", s)
		ret++
		return false
	})
	return ret
}
