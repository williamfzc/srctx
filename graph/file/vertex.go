package file

type Vertex struct {
	Path       string
	Referenced int

	// https://github.com/williamfzc/srctx/issues/41
	Tags map[string]struct{} `json:"tags,omitempty"`
}

func (fv *Vertex) Id() string {
	return fv.Path
}

func (fv *Vertex) ContainTag(tag string) bool {
	if _, ok := fv.Tags[tag]; ok {
		return true
	}
	return false
}

func (fv *Vertex) AddTag(tag string) {
	fv.Tags[tag] = struct{}{}
}

func (fv *Vertex) RemoveTag(tag string) {
	delete(fv.Tags, tag)
}
