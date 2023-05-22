package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

func File2Tokens(fileName string, line int) ([]chroma.Token, error) {
	lexer := lexers.Match(fileName)
	if lexer == nil {
		return nil, errors.New("no lexer matches " + fileName)
	}
	contents, err := ReadLine(fileName, line)
	if err != nil {
		return nil, err
	}
	tokens, err := lexer.Tokenise(nil, contents)
	if err != nil {
		return nil, err
	}
	return tokens.Tokens(), nil
}

func ReadLine(fileName string, lineNum int) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for i := 1; scanner.Scan(); i++ {
		if i == lineNum {
			return scanner.Text(), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("line %d not found in file %s", lineNum, fileName)
}
