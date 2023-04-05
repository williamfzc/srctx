package collector

import (
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type FunctionMetaMap = map[string]*extractor.FunctionFileResult

func GetFunctionsMetadata(rootDir string) (FunctionMetaMap, error) {
	functionFileResults, err := sibyl2.ExtractFunction(rootDir, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}
	fmm := make(FunctionMetaMap, len(functionFileResults))
	for _, each := range functionFileResults {
		fmm[each.Path] = each
	}
	return fmm, nil
}
