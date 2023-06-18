package function

import (
	"fmt"
)

const (
	TagYellow = "yellow"
	TagRed    = "red"
	TagOrange = "orange"
)

func (fg *FuncGraph) FillWithYellow(vertexHash string) error {
	item, ok := fg.IdCache[vertexHash]
	if !ok {
		return fmt.Errorf("no such vertex: %v", vertexHash)
	}
	item.AddTag(TagYellow)
	return nil
}

func (fg *FuncGraph) FillWithOrange(vertexHash string) error {
	item, ok := fg.IdCache[vertexHash]
	if !ok {
		return fmt.Errorf("no such vertex: %v", vertexHash)
	}
	item.AddTag(TagOrange)
	return nil
}

func (fg *FuncGraph) FillWithRed(vertexHash string) error {
	item, ok := fg.IdCache[vertexHash]
	if !ok {
		return fmt.Errorf("no such vertex: %v", vertexHash)
	}
	item.AddTag(TagRed)
	return nil
}
