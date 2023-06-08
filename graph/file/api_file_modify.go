package file

import (
	"github.com/williamfzc/srctx/graph/utils"
)

func (fg *FileGraph) RemoveNodeById(path string) error {
	err := utils.RemoveFromGraph(fg.G, path)
	if err != nil {
		return err
	}
	err = utils.RemoveFromGraph(fg.Rg, path)
	if err != nil {
		return err
	}
	return nil
}
