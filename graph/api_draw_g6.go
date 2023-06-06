package graph

import (
	"fmt"
	"os"
	"strconv"

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
<script src="https://gw.alipayobjects.com/os/lib/antv/g6/4.3.11/dist/g6.min.js"></script>

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
			type: 'gForce',
            preventOverlap: true,
            linkDistance: 100,
            nodeSize: 100
		},
		modes: {
            default: ['drag-canvas', 'zoom-canvas', 'drag-node', 'activate-relations'],
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
		plugins: [grid, toolbar, edgeBundling]
    });
    graph.data(data);
    graph.render();

	const toggleLayoutButton = document.getElementById('toggleLayoutButton');
    const layoutTypes = ['gForce', 'circular', 'dagre', 'radial', 'random', 'concentric'];

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

// G6Node https://g6.antv.antgroup.com/api/shape-properties
type G6Node struct {
	Id    string `json:"id"`
	Label string `json:"label,omitempty"`
	Fill  string `json:"fill,omitempty"`
}

type G6Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// G6Data https://g6.antv.antgroup.com/api/graph-func/data
type G6Data struct {
	Nodes []*G6Node `json:"nodes"`
	Edges []*G6Edge `json:"edges"`
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

func (fg *FuncGraph) ToG6Data() (*G6Data, error) {
	storage, err := fg.Dump()
	if err != nil {
		return nil, err
	}

	data := &G6Data{
		Nodes: make([]*G6Node, 0, len(storage.VertexIds)),
		Edges: make([]*G6Edge, 0),
	}
	// Nodes
	for nodeId, funcId := range storage.VertexIds {
		curNode := &G6Node{
			Id:    strconv.Itoa(nodeId),
			Label: funcId,
		}
		data.Nodes = append(data.Nodes, curNode)
	}
	// Edges
	for src, targets := range storage.GEdges {
		for _, target := range targets {
			curEdge := &G6Edge{
				Source: strconv.Itoa(src),
				Target: strconv.Itoa(target),
			}
			data.Edges = append(data.Edges, curEdge)
		}
	}
	return data, nil
}

func (fg *FuncGraph) DrawG6Html(filename string) error {
	g6data, err := fg.ToG6Data()
	if err != nil {
		return err
	}

	err = g6data.RenderHtml(filename)
	if err != nil {
		return err
	}
	return nil
}
