package graph

import (
	"github.com/williamfzc/srctx/graph/utils"
)

func (fg *FuncGraph) RemoveNodeById(funcId string) error {
	err := utils.RemoveFromGraph(fg.g, funcId)
	if err != nil {
		return err
	}
	err = utils.RemoveFromGraph(fg.rg, funcId)
	if err != nil {
		return err
	}
	return nil
}
