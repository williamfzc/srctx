package graph

import (
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

type FgStorage struct {
	VertexIds map[int]string           `json:"vertexIds"`
	GEdges    map[int][]int            `json:"gEdges"`
	RGEdges   map[int][]int            `json:"rgEdges"`
	Cache     map[string][]*FuncVertex `json:"cache"`
}

func (fg *FuncGraph) DumpFile(fp string) error {
	storage, err := fg.Dump()
	if err != nil {
		return err
	}
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := msgpack.NewEncoder(file)
	encoder.SetCustomStructTag("json")
	if err := encoder.Encode(storage); err != nil {
		return err
	}
	return nil
}

func (fg *FuncGraph) Dump() (*FgStorage, error) {
	ret := &FgStorage{
		VertexIds: make(map[int]string),
		GEdges:    make(map[int][]int),
		RGEdges:   make(map[int][]int),
		Cache:     nil,
	}
	ret.Cache = fg.cache

	allVertices := make([]*FuncVertex, 0)
	for _, vertices := range ret.Cache {
		for _, each := range vertices {
			allVertices = append(allVertices, each)
		}
	}
	reverseMapping := make(map[string]int)
	for index, each := range allVertices {
		eachId := each.Id()
		ret.VertexIds[index] = eachId
		reverseMapping[eachId] = index
	}

	adjacencyMap, err := fg.g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	for src, v := range adjacencyMap {
		for tar := range v {
			srcIndex := reverseMapping[src]
			tarIndex := reverseMapping[tar]
			ret.GEdges[srcIndex] = append(ret.GEdges[srcIndex], tarIndex)
		}
	}

	// so does rg
	adjacencyMap, err = fg.rg.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	for src, v := range adjacencyMap {
		for tar := range v {
			srcIndex := reverseMapping[src]
			tarIndex := reverseMapping[tar]
			ret.RGEdges[srcIndex] = append(ret.RGEdges[srcIndex], tarIndex)
		}
	}

	return ret, err
}

func Load(fgs *FgStorage) (*FuncGraph, error) {
	ret := NewEmptyFuncGraph()

	// vertex building
	ret.cache = fgs.Cache
	for _, eachFile := range ret.cache {
		for _, eachFunc := range eachFile {
			_ = ret.g.AddVertex(eachFunc)
			_ = ret.rg.AddVertex(eachFunc)
		}
	}

	// edge building
	mapping := fgs.VertexIds
	for srcId, targets := range fgs.GEdges {
		for _, tarId := range targets {
			src := mapping[srcId]
			tar := mapping[tarId]
			_ = ret.g.AddEdge(src, tar)
		}
	}
	for srcId, targets := range fgs.RGEdges {
		for _, tarId := range targets {
			src := mapping[srcId]
			tar := mapping[tarId]
			_ = ret.rg.AddEdge(src, tar)
		}
	}
	return ret, nil
}

func LoadFile(fp string) (*FuncGraph, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := msgpack.NewDecoder(file)
	storage := &FgStorage{}
	decoder.SetCustomStructTag("json")
	if err := decoder.Decode(storage); err != nil {
		return nil, err
	}

	fg, err := Load(storage)
	if err != nil {
		return nil, err
	}

	return fg, nil
}
