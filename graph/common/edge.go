package common

type EdgeStorage struct {
	RefLines map[int]struct{}
}

func NewEdgeStorage() *EdgeStorage {
	return &EdgeStorage{RefLines: make(map[int]struct{})}
}
