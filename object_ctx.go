package srctx

import "github.com/dominikbraun/graph"

type SourceContext struct {
	FactGraph graph.Graph[int, *FactVertex]
	RelGraph  graph.Graph[int, *RelVertex]
}
