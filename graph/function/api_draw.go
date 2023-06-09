package function

import (
	"os"

	"github.com/dominikbraun/graph/draw"
)

type Drawable interface {
	DrawDot(fileName string) error
	FillWithYellow(vertexHash string) error
	FillWithRed(vertexHash string) error
}

func (fg *FuncGraph) DrawDot(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// draw the call graph
	err = draw.DOT(fg.g, file, draw.GraphAttribute("rankdir", "LR"))
	if err != nil {
		return err
	}
	return nil
}

func (fg *FuncGraph) FillWithYellow(vertexHash string) error {
	err := fg.setProperty(vertexHash, "style", "filled")
	if err != nil {
		return err
	}
	err = fg.setProperty(vertexHash, "color", "yellow")
	if err != nil {
		return err
	}
	return nil
}

func (fg *FuncGraph) FillWithRed(vertexHash string) error {
	err := fg.setProperty(vertexHash, "style", "filled")
	if err != nil {
		return err
	}
	err = fg.setProperty(vertexHash, "color", "red")
	if err != nil {
		return err
	}
	return nil
}

func (fg *FuncGraph) setProperty(vertexHash string, propertyK string, propertyV string) error {
	_, properties, err := fg.g.VertexWithProperties(vertexHash)
	if err != nil {
		return err
	}
	properties.Attributes[propertyK] = propertyV
	return nil
}
