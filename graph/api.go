package graph

func (fg *FuncGraph) GetFunctionsByFile(f string) []*FuncVertex {
	return fg.cache[f]
}

func (fg *FuncGraph) GetById(id string) (*FuncVertex, error) {
	return fg.g.Vertex(id)
}
