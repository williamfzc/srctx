package function

import (
	"os"

	"github.com/dominikbraun/graph/draw"
)

func (fg *Graph) setProperty(vertexHash string, propertyK string, propertyV string) error {
	_, properties, err := fg.g.VertexWithProperties(vertexHash)
	if err != nil {
		return err
	}
	properties.Attributes[propertyK] = propertyV
	return nil
}

func (fg *Graph) DrawDot(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	for _, each := range fg.ListFunctions() {
		if each.ContainTag(TagYellow) {
			err := fg.setProperty(each.Id(), "style", "filled")
			if err != nil {
				return err
			}
			err = fg.setProperty(each.Id(), "color", "yellow")
			if err != nil {
				return err
			}
		}

		if each.ContainTag(TagRed) {
			err := fg.setProperty(each.Id(), "style", "filled")
			if err != nil {
				return err
			}
			err = fg.setProperty(each.Id(), "color", "red")
			if err != nil {
				return err
			}
		}
	}

	// draw the call graph
	err = draw.DOT(fg.g, file, draw.GraphAttribute("rankdir", "LR"))
	if err != nil {
		return err
	}
	return nil
}
