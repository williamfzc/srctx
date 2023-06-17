package function

func (fg *FuncGraph) ListFiles() []string {
	ret := make([]string, 0, len(fg.Cache))
	for k := range fg.Cache {
		ret = append(ret, k)
	}
	return ret
}

func (fg *FuncGraph) FileCount() int {
	return len(fg.ListFiles())
}
