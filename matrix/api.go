package matrix

import (
	"github.com/dominikbraun/graph"
	"github.com/williamfzc/srctx/graph/common"
	"gonum.org/v1/gonum/mat"
)

const (
	InvalidIndex = -1
	InvalidValue = -1.0
)

/*
Matrix

Standard relationship representation layer based on matrix
*/
type Matrix struct {
	IndexMap        map[string]int
	ReverseIndexMap map[int]string
	Data            *mat.Dense
}

func (m *Matrix) Size() int {
	return len(m.IndexMap)
}

func (m *Matrix) Id(s string) int {
	if item, ok := m.IndexMap[s]; ok {
		return item
	}
	return InvalidIndex
}

func (m *Matrix) ById(id int) string {
	if item, ok := m.ReverseIndexMap[id]; ok {
		return item
	}
	return ""
}

func (m *Matrix) ForEach(s string, f func(i int, v float64)) {
	index := m.Id(s)
	if index == InvalidIndex {
		return
	}

	for i := 0; i < m.Size(); i++ {
		f(i, m.Data.At(i, index))
	}
}

func CreateMatrixFromGraph[T string, U any](g graph.Graph[T, U]) (*Matrix, error) {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	nodeCount := len(adjacencyMap)
	data := mat.NewDense(nodeCount, nodeCount, nil)

	indexMap := make(map[string]int)
	i := 0
	for node := range adjacencyMap {
		indexMap[string(node)] = i
		i++
	}

	for source, edges := range adjacencyMap {
		sourceIndex := indexMap[string(source)]

		data.Set(sourceIndex, sourceIndex, float64(len(edges)))

		for target, edge := range edges {
			storage := edge.Properties.Data.(*common.EdgeStorage)
			targetIndex := indexMap[string(target)]

			currentValue := data.At(targetIndex, sourceIndex)
			data.Set(targetIndex, sourceIndex, currentValue+float64(len(storage.RefLines)))
		}
	}

	ret := &Matrix{
		IndexMap:        indexMap,
		ReverseIndexMap: reverseMap(indexMap),
		Data:            data,
	}
	return ret, nil
}

func reverseMap(originalMap map[string]int) map[int]string {
	reversedMap := make(map[int]string)

	for key, value := range originalMap {
		reversedMap[value] = key
	}

	return reversedMap
}
