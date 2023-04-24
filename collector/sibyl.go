package collector

import (
	"path/filepath"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	log "github.com/sirupsen/logrus"
)

func CreateFact(root string) (*FactStorage, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	functionFiles, err := sibyl2.ExtractFunction(abs, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}

	fact := &FactStorage{
		cache: make(map[string]*extractor.FunctionFileResult),
	}
	for _, eachFunc := range functionFiles {
		log.Infof("create func file for: %v", eachFunc.Path)
		fact.cache[eachFunc.Path] = eachFunc
	}
	return fact, nil
}

type FactStorage struct {
	cache map[string]*extractor.FunctionFileResult
}

func (fs *FactStorage) GetByFile(fileName string) *extractor.FunctionFileResult {
	return fs.cache[fileName]
}
