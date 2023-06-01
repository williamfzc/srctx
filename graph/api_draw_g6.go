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
</head>
<body>
<div id="mountNode"></div>
<script src="https://gw.alipayobjects.com/os/lib/antv/g6/4.3.11/dist/g6.min.js"></script>

<script>
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
    });
    graph.data(data);
    graph.render();
</script>
</body>
</html>
`

type g6node struct {
	Id    string `json:"id"`
	Label string `json:"label,omitempty"`
}

type g6edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// https://g6.antv.antgroup.com/api/graph-func/data
type g6data struct {
	Nodes []*g6node `json:"nodes"`
	Edges []*g6edge `json:"edges"`
}

func (fg *FuncGraph) DrawG6Html(filename string) error {
	storage, err := fg.Dump()
	if err != nil {
		return err
	}

	data := &g6data{
		Nodes: make([]*g6node, 0, len(storage.VertexIds)),
		Edges: make([]*g6edge, 0),
	}
	// Nodes
	for nodeId, funcId := range storage.VertexIds {
		curNode := &g6node{
			Id:    strconv.Itoa(nodeId),
			Label: funcId,
		}
		data.Nodes = append(data.Nodes, curNode)
	}
	// Edges
	for src, targets := range storage.GEdges {
		for _, target := range targets {
			curEdge := &g6edge{
				Source: strconv.Itoa(src),
				Target: strconv.Itoa(target),
			}
			data.Edges = append(data.Edges, curEdge)
		}
	}
	// render
	dataRaw, err := json.Marshal(data)
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
