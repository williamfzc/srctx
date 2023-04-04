package parser

import (
	"github.com/williamfzc/srctx/parser/lsif"
)

func reverseMap(m map[lsif.Id]lsif.Id) map[lsif.Id]lsif.Id {
	n := make(map[lsif.Id]lsif.Id, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
