package main

import (
	"encoding/json"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

func exportJson(outputFile string, lineStats []*LineStat) {
	data, err := json.Marshal(lineStats)
	panicIfErr(err)
	err = os.WriteFile(outputFile, data, 0644)
	panicIfErr(err)
	log.Infof("dump json to %s", outputFile)
}

func exportCsv(outputFile string, lineStats []*LineStat) {
	csvFile, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	panicIfErr(err)
	defer csvFile.Close()

	unsafeLines := make([]*LineStat, 0)
	for _, each := range lineStats {
		if !each.IsSafe() {
			unsafeLines = append(unsafeLines, each)
		}
	}

	if err := gocsv.MarshalFile(&unsafeLines, csvFile); err != nil { // Load clients from file
		panic(err)
	}
	log.Infof("dump csv to %s", outputFile)
}

func exportDot(outputFile string, fileRefMap map[string]*fileVertex) {
	// only create a file level graph
	fileGraph := graph.New((*fileVertex).Id, graph.Directed())
	for _, vertex := range fileRefMap {
		if vertex.Directly {
			_ = fileGraph.AddVertex(vertex, func(vertexProperties *graph.VertexProperties) {
				vertexProperties.Attributes["style"] = "filled"
				vertexProperties.Attributes["fillcolor"] = "yellow"
			})
		} else {
			_ = fileGraph.AddVertex(vertex)
		}
	}
	for _, vertex := range fileRefMap {
		for _, eachRef := range vertex.Refs {
			// ignore self ref
			if eachRef != vertex.Id() {
				_ = fileGraph.AddEdge(eachRef, vertex.Id())
			}
		}
	}
	f, _ := os.Create(outputFile)
	_ = draw.DOT(fileGraph, f)
	log.Infof("dump dot to %s", outputFile)
}
