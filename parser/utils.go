package parser

import (
	"github.com/williamfzc/srctx/parser/lsif"
)

func reverseMap(m map[lsif.Id]lsif.Id) map[lsif.Id][]lsif.Id {
	n := make(map[lsif.Id][]lsif.Id)
	for k, v := range m {
		if _, ok := n[v]; !ok {
			n[v] = make([]lsif.Id, 0)
		}
		n[v] = append(n[v], k)
	}
	return n
}
