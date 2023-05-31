package graph

import (
	"path/filepath"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	log "github.com/sirupsen/logrus"
)

func CreateFact(root string, lang core.LangType) (*FactStorage, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	conf := sibyl2.DefaultConfig()
	conf.LangType = lang
	functionFiles, err := sibyl2.ExtractFunction(abs, conf)
	if err != nil {
		return nil, err
	}
	symbolFiles, err := sibyl2.ExtractSymbol(abs, conf)
	if err != nil {
		return nil, err
	}

	fact := &FactStorage{
		cache:       make(map[string]*extractor.FunctionFileResult, len(functionFiles)),
		symbolCache: make(map[string]*extractor.SymbolFileResult, len(symbolFiles)),
	}
	for _, eachFunc := range functionFiles {
		log.Debugf("create func file for: %v", eachFunc.Path)
		fact.cache[eachFunc.Path] = eachFunc
	}
	for _, eachSymbol := range symbolFiles {
		log.Debugf("create symbol file for: %v", eachSymbol.Path)
		fact.symbolCache[eachSymbol.Path] = eachSymbol
	}

	return fact, nil
}

// FactStorage
// fact is some extra metadata extracted from source code
// something like: function definitions with their annotations/params/receiver ...
// these data can be used for enhancing relationship
type FactStorage struct {
	cache       map[string]*extractor.FunctionFileResult
	symbolCache map[string]*extractor.SymbolFileResult
}

func (fs *FactStorage) GetFunctionsByFile(fileName string) *extractor.FunctionFileResult {
	return fs.cache[fileName]
}

func (fs *FactStorage) GetSymbolsByFileAndLine(fileName string, line int) []*extractor.Symbol {
	item, ok := fs.symbolCache[fileName]
	if !ok {
		log.Warnf("failed to get symbol: %v", fileName)
		return nil
	}
	ret := make([]*extractor.Symbol, 0)
	for _, eachUnit := range item.Units {
		if eachUnit.GetSpan().ContainLine(line - 1) {
			ret = append(ret, eachUnit)
		}
	}
	return ret
}
