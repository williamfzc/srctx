package srctx

import "github.com/williamfzc/srctx/parser"

func reverseMap(m map[parser.Id]parser.Id) map[parser.Id]parser.Id {
	n := make(map[parser.Id]parser.Id, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
