package lexer

import (
	"bufio"
	"os"
)

func Tokenize(filePath string) ([]Token, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lexer := NewLexer(reader)

	for {
		tok := lexer.NextToken()
		lexer.Tokens = append(lexer.Tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	return lexer.Tokens, nil
}
