package function

import "fmt"

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
	if item, ok := fg.IdCache[id]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("id not found in graph: %s", id)
}

func (fg *FuncGraph) FuncCount() int {
	return len(fg.IdCache)
}

func (fg *FuncGraph) ListFunctions() []*FuncVertex {
	return fg.FilterFunctions(func(funcVertex *FuncVertex) bool {
		return true
	})
}

func (fg *FuncGraph) FilterFunctions(f func(*FuncVertex) bool) []*FuncVertex {
	ret := make([]*FuncVertex, 0, len(fg.IdCache))
	for _, each := range ret {
		if f(each) {
			ret = append(ret, each)
		}
	}
	return ret
}
