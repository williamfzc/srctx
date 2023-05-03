package graph

import (
	"os"

	"github.com/dominikbraun/graph/draw"
)

func (fg *FuncGraph) DrawDot(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = draw.DOT(fg.g, file)
	if err != nil {
		return err
	}
	return nil
}

func (fg *FuncGraph) SetProperty(vertexHash string, propertyK string, propertyV string) error {
	_, properties, err := fg.g.VertexWithProperties(vertexHash)
	if err != nil {
		return err
	}
	properties.Attributes[propertyK] = propertyV
	return nil
}

func (fg *FuncGraph) Highlight(vertexHash string) error {
	return fg.FillWithYellow(vertexHash)
}

func (fg *FuncGraph) FillWithYellow(vertexHash string) error {
	err := fg.SetProperty(vertexHash, "style", "filled")
	if err != nil {
		return err
	}
	err = fg.SetProperty(vertexHash, "color", "yellow")
	if err != nil {
		return err
	}
	return nil
}

func (fg *FuncGraph) FillWithRed(vertexHash string) error {
	err := fg.SetProperty(vertexHash, "style", "filled")
	if err != nil {
		return err
	}
	err = fg.SetProperty(vertexHash, "color", "red")
	if err != nil {
		return err
	}
	return nil
}
