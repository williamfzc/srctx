package lexer

import (
	"errors"
	"os"
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/williamfzc/srctx/object"
)

// not good but temp
var tokenCache = make(map[string][][]chroma.Token)
var l sync.Mutex

func File2Tokens(fileName string) ([][]chroma.Token, error) {
	if r, ok := tokenCache[fileName]; ok {
		return r, nil
	}
	lexer := lexers.Match(fileName)
	if lexer == nil {
		return nil, errors.New("no lexer matches " + fileName)
	}
	contents, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	tokens, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		return nil, err
	}
	ret := chroma.SplitTokensIntoLines(tokens.Tokens())

	l.Lock()
	defer l.Unlock()
	tokenCache[fileName] = ret
	return ret, nil
}

func TypeFromTokens(tokens []chroma.Token) object.DefType {
	for _, token := range tokens {
		switch token.Type {
		case chroma.NameFunction:
			return object.DefFunction
		case chroma.NameClass:
			return object.DefClass
		case chroma.NameNamespace:
			return object.DefNamespace

		default:
			continue
		}
	}
	return object.DefUnknown
}
