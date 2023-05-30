package graph

import (
	"encoding/json"
	"os"
)

type FgStorage struct {
	GEdges  map[string]string        `json:"GEdges,omitempty"`
	RGEdges map[string]string        `json:"RGEdges,omitempty"`
	Cache   map[string][]*FuncVertex `json:"cache,omitempty"`
}

func (fg *FuncGraph) DumpJsonFile(fp string) error {
	storage, err := fg.Dump()
	if err != nil {
		return err
	}
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "")
	if err := encoder.Encode(storage); err != nil {
		return err
	}
	return nil
}

func (fg *FuncGraph) Dump() (*FgStorage, error) {
	ret := &FgStorage{
		GEdges:  make(map[string]string),
		RGEdges: make(map[string]string),
		Cache:   nil,
	}
	ret.Cache = fg.cache

	adjacencyMap, err := fg.g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	for src, v := range adjacencyMap {
		for tar := range v {
			ret.GEdges[src] = tar
		}
	}

	// so does rg
	adjacencyMap, err = fg.rg.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	for src, v := range adjacencyMap {
		for tar := range v {
			ret.RGEdges[src] = tar
		}
	}

	return ret, err
}
