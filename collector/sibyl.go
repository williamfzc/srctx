package collector

import (
	"fmt"
	"strings"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type FunctionMetaMap = map[string]*extractor.FunctionFileResult

func IsSupported(path string) bool {
	if strings.HasSuffix(path, ".java") {
		return true
	}
	if strings.HasSuffix(path, ".kt") {
		return true
	}
	if strings.HasSuffix(path, ".go") {
		return true
	}
	if strings.HasSuffix(path, ".py") {
		return true
	}
	if strings.HasSuffix(path, ".js") {
		return true
	}
	return false
}

type NotSupportLangError struct {
	msg string
}

func (e *NotSupportLangError) Error() string {
	return fmt.Sprintf("not supported lang: %s", e.msg)
}

func GetFunctionMetadataFromFile(targetFile string) (*extractor.FunctionFileResult, error) {
	if !IsSupported(targetFile) {
		return nil, &NotSupportLangError{targetFile}
	}
	res, err := sibyl2.ExtractFunction(targetFile, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return res[0], nil
}
