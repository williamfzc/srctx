package g6

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
)

const g6template = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <title>srctx report</title>
	<style>
        body {
            margin: 0;
            padding: 0;
            font-family: Arial, sans-serif;
        }
        #mountNode {
            width: 100%%;
            height: 100vh;
        }
        #toggleLayoutButton {
            position: absolute;
            top: 10px;
            right: 10px;
            padding: 10px;
            font-size: 16px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        #toggleLayoutButton:hover {
            background-color: #3e8e41;
        }
    </style>
</head>
<body>
<button id="toggleLayoutButton">Change Layout</button>
<div id="mountNode"></div>
<script src="https://gw.alipayobjects.com/os/lib/antv/g6/4.8.7/dist/g6.min.js"></script>

<script>
	const grid = new G6.Grid()
	const toolbar = new G6.ToolBar()
	const edgeBundling = new G6.Bundling({
      bundleThreshold: 0.6,
      K: 100,
    })

    const data = %s

    const graph = new G6.Graph({
        container: 'mountNode',
        width: window.innerWidth,
        height: window.innerHeight,
		layout: {
            type: 'comboCombined',
			preventOverlap: true,
            spacing: 5,
			linkDistance: 100,
            nodeSize: 100,
            outerLayout: new G6.Layout['forceAtlas2']({
                kr: 50,
                factor: 10,
            })
        },
		modes: {
            default: ['drag-canvas', 'zoom-canvas', 'drag-node', 'activate-relations', 'collapse-expand-combo'],
        },
        defaultNode: {
            size: 60,
            style: {
                lineWidth: 1,
            },
        },
        defaultEdge: {
            style: {
                opacity: 0.6,
                stroke: 'black',
                startArrow: true,
            },
            labelCfg: {
                autoRotate: true,
            },
        },
		defaultCombo: {
            type: 'rect',
            collapse: true,
            style: {
                stroke: 'black',
            }
        },
		groupByTypes: false,
		plugins: [grid, toolbar, edgeBundling]
    });
    graph.data(data);
    graph.render();

	const toggleLayoutButton = document.getElementById('toggleLayoutButton');
    const layoutTypes = ['comboCombined'];

    let currentLayoutIndex = 0;
    toggleLayoutButton.addEventListener('click', function() {
        currentLayoutIndex = (currentLayoutIndex + 1) %% layoutTypes.length;
        const layoutType = layoutTypes[currentLayoutIndex];
        graph.updateLayout({ type: layoutType });
    });
</script>
</body>
</html>
`

type G6Node struct {
	Id      string       `json:"id"`
	Label   string       `json:"label,omitempty"`
	ComboId string       `json:"comboId,omitempty"`
	Style   *G6NodeStyle `json:"style"`
}

// G6NodeStyle https://g6.antv.antgroup.com/api/shape-properties
type G6NodeStyle struct {
	Fill string `json:"fill,omitempty"`
}

type G6Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type G6Combo struct {
	Id        string `json:"id"`
	Label     string `json:"label"`
	Collapsed bool   `json:"collapsed,omitempty"`
}

// G6Data https://g6.antv.antgroup.com/api/graph-func/data
type G6Data struct {
	Nodes  []*G6Node  `json:"nodes"`
	Edges  []*G6Edge  `json:"edges"`
	Combos []*G6Combo `json:"combos"`
}

func EmptyG6Data() *G6Data {
	return &G6Data{
		Nodes:  make([]*G6Node, 0),
		Edges:  make([]*G6Edge, 0),
		Combos: make([]*G6Combo, 0),
	}
}

func (g *G6Data) RenderHtml(filename string) error {
	// render
	dataRaw, err := json.Marshal(g)
	if err != nil {
		return nil
	}
	htmlContent := fmt.Sprintf(g6template, dataRaw)
	err = os.WriteFile(filename, []byte(htmlContent), 0o666)
	if err != nil {
		return err
	}

	return nil
}

func (g *G6Data) FillWithYellow(label string) {
	for _, each := range g.Nodes {
		if each.Label == label {
			each.Style.Fill = "yellow"
			break
		}
	}
}

func (g *G6Data) FillWithRed(label string) {
	for _, each := range g.Nodes {
		if each.Label == label {
			each.Style.Fill = "red"
			break
		}
	}
}
