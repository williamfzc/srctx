package function

import (
	"path/filepath"
	"strconv"

	"github.com/williamfzc/srctx/graph/visual/g6"
)

func (fg *FuncGraph) ToG6Data() (*g6.Data, error) {
	storage, err := fg.Dump()
	if err != nil {
		return nil, err
	}

	data := g6.EmptyG6Data()
	// cache
	cache := make(map[string]*FuncVertex)
	// dir combos (#35)
	dirCombos := make(map[string]struct{})
	for eachFile, fs := range fg.Cache {
		for _, eachF := range fs {
			cache[eachF.Id()] = eachF
		}

		eachDir := filepath.Dir(eachFile)
		dirCombos[eachDir] = struct{}{}
		data.Combos = append(data.Combos, &g6.Combo{
			Id:        eachFile,
			Label:     eachFile,
			Collapsed: false,
			ParentId:  eachDir,
		})
	}

	for eachDir := range dirCombos {
		data.Combos = append(data.Combos, &g6.Combo{
			Id:        eachDir,
			Label:     eachDir,
			Collapsed: false,
		})
	}

	// Nodes
	for nodeId, funcId := range storage.VertexIds {
		curNode := &g6.Node{
			Id:      strconv.Itoa(nodeId),
			Label:   funcId,
			Style:   &g6.NodeStyle{},
			ComboId: cache[funcId].Path,
		}
		data.Nodes = append(data.Nodes, curNode)
	}
	// Edges
	for src, targets := range storage.GEdges {
		for _, target := range targets {
			curEdge := &g6.Edge{
				Source: strconv.Itoa(src),
				Target: strconv.Itoa(target),
			}
			data.Edges = append(data.Edges, curEdge)
		}
	}
	return data, nil
}

func (fg *FuncGraph) DrawG6Html(filename string) error {
	g6data, err := fg.ToG6Data()
	if err != nil {
		return err
	}

	err = g6data.RenderHtml(filename)
	if err != nil {
		return err
	}
	return nil
}
