package graph

import "github.com/dominikbraun/graph"

func removeFromGraph[T string, U any](g graph.Graph[T, *U], hash T) error {
	// to this hash
	pm, err := g.PredecessorMap()
	if err != nil {
		return err
	}
	// from this hash
	am, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	v, ok := pm[hash]
	if ok {
		for s := range v {
			err = g.RemoveEdge(s, hash)
			if err != nil {
				return err
			}
		}
	}

	v, ok = am[hash]
	if ok {
		for s := range v {
			err = g.RemoveEdge(hash, s)
			if err != nil {
				return err
			}
		}
	}

	err = g.RemoveVertex(hash)
	if err != nil {
		return err
	}

	return nil
}

func (fg *FuncGraph) RemoveNodeById(funcId string) error {
	err := removeFromGraph(fg.g, funcId)
	if err != nil {
		return err
	}
	err = removeFromGraph(fg.rg, funcId)
	if err != nil {
		return err
	}
	return nil
}
