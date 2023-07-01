package file

import (
	"os"

	"github.com/dominikbraun/graph/draw"
)

func (fg *Graph) DrawDot(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// draw the call graph
	err = draw.DOT(fg.G, file, draw.GraphAttribute("rankdir", "LR"))
	if err != nil {
		return err
	}
	return nil
}

func (fg *Graph) FillWithYellow(vertexHash string) error {
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

func (fg *Graph) FillWithRed(vertexHash string) error {
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

func (fg *Graph) setProperty(vertexHash string, propertyK string, propertyV string) error {
	_, properties, err := fg.G.VertexWithProperties(vertexHash)
	if err != nil {
		return err
	}
	properties.Attributes[propertyK] = propertyV
	return nil
}
