package g6

import (
	"bytes"
	_ "embed"
	"os"
	"text/template"

	"github.com/goccy/go-json"
)

//go:embed template/report.html
var g6ReportTemplate string

type Node struct {
	Id      string     `json:"id"`
	Label   string     `json:"label,omitempty"`
	ComboId string     `json:"comboId,omitempty"`
	Style   *NodeStyle `json:"style"`
}

// NodeStyle https://g6.antv.antgroup.com/api/shape-properties
type NodeStyle struct {
	Fill string `json:"fill,omitempty"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Combo struct {
	Id        string `json:"id"`
	Label     string `json:"label"`
	Collapsed bool   `json:"collapsed,omitempty"`
	ParentId  string `json:"parentId,omitempty"`
}

// Data https://g6.antv.antgroup.com/api/graph-func/data
type Data struct {
	Nodes  []*Node  `json:"nodes"`
	Edges  []*Edge  `json:"edges"`
	Combos []*Combo `json:"combos"`
}

func EmptyG6Data() *Data {
	return &Data{
		Nodes:  make([]*Node, 0),
		Edges:  make([]*Edge, 0),
		Combos: make([]*Combo, 0),
	}
}

func (g *Data) RenderHtml(filename string) error {
	// render
	dataRaw, err := json.Marshal(g)
	if err != nil {
		return nil
	}

	parsed, err := template.New("").Parse(g6ReportTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = parsed.Execute(&buf, string(dataRaw))
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, buf.Bytes(), 0o666)
	if err != nil {
		return err
	}

	return nil
}

func (g *Data) FillWithYellow(label string) {
	for _, each := range g.Nodes {
		if each.Label == label {
			each.Style.Fill = "yellow"
			break
		}
	}
}

func (g *Data) FillWithRed(label string) {
	for _, each := range g.Nodes {
		if each.Label == label {
			each.Style.Fill = "red"
			break
		}
	}
}
