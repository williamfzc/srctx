package function

func (fg *FuncGraph) GetFunctionsByFile(fileName string) []*FuncVertex {
	if item, ok := fg.Cache[fileName]; ok {
		return item
	}
	return make([]*FuncVertex, 0)
}

func (fg *FuncGraph) GetFunctionsByFileLines(fileName string, lines []int) []*FuncVertex {
	ret := make([]*FuncVertex, 0)
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

func (fg *FuncGraph) GetById(id string) (*FuncVertex, error) {
	return fg.g.Vertex(id)
}
