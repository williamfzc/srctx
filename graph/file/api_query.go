package file

import "github.com/sirupsen/logrus"

func (fg *Graph) GetById(id string) *Vertex {
	v, err := fg.G.Vertex(id)
	if err != nil {

		logrus.Warnf("no vertex: %v", id)
		return nil
	}
	return v
}

func (fg *Graph) ListFiles() []*Vertex {
	ret := make([]*Vertex, 0, len(fg.IdCache))
	for _, each := range fg.IdCache {
		ret = append(ret, each)
	}
	return ret
}
