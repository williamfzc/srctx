package graph

func (fg *FileGraph) RemoveNodeById(path string) error {
	err := removeFromGraph(fg.g, path)
	if err != nil {
		return err
	}
	err = removeFromGraph(fg.rg, path)
	if err != nil {
		return err
	}
	return nil
}
