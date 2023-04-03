package lexer

import (
	"errors"
	"os"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// not good but temp
var tokenCache = make(map[string][][]chroma.Token)

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
	tokenCache[fileName] = ret
	return ret, nil
}
