package main

type fileVertex struct {
	Name     string
	Refs     []string
	Directly bool
}

func (vertex *fileVertex) Id() string {
	return vertex.Name
}
