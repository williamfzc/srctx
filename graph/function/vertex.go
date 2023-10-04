package function

import (
	"fmt"

	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type Vertex struct {
	*object.Function
	*FuncPos

	// https://github.com/williamfzc/srctx/issues/41
	Tags map[string]struct{} `json:"tags,omitempty"`
}

func (fv *Vertex) Id() string {
	return fmt.Sprintf("%v:#%d-#%d:%s", fv.Path, fv.Start, fv.End, fv.GetSignature())
}

func (fv *Vertex) PosKey() string {
	return fmt.Sprintf("%s#%d", fv.Path, fv.Start)
}

func (fv *Vertex) ListTags() []string {
	ret := make([]string, 0, len(fv.Tags))
	for each := range fv.Tags {
		ret = append(ret, each)
	}
	return ret
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

func CreateFuncVertex(f *object.Function, fr *extractor.FunctionFileResult) *Vertex {
	cur := &Vertex{
		Function: f,
		FuncPos: &FuncPos{
			Path: fr.Path,
			Lang: string(fr.Language),
			// sync with real lines
			Start: int(f.GetSpan().Start.Row + 1),
			End:   int(f.GetSpan().End.Row + 1),
		},
		Tags: make(map[string]struct{}),
	}
	return cur
}
