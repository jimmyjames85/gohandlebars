package parser

import (
	"fmt"

	"github.com/jimmyjames85/gohandlebars/lexer"
	"github.com/pkg/errors"
)

type NonTerminal interface{}

func Parse(data []byte) error {

	tokens, err := lexer.Scan(data)
	if err != nil {
		return errors.Wrap(err, "lexing error")
	}

	for _, tkn := range tokens {
		fmt.Printf("%s: %s\n\n", tkn.Typ, tkn.Val)
	}

	return nil
}

func parseStatement(tokens []*lexer.Token) error {
	//right now stmt is "return 3;"
	if len(tokens) < 3 {
		return fmt.Errorf("not enough arguments for statement")
	}

	expected := []lexer.TokenType{lexer.Return, lexer.IntLiteral, lexer.Semicolon}

	for i, e := range expected {
		if tokens[i].Typ != e {
			return fmt.Errorf("expecting %s: got %q", e, tokens[i].Typ)
		}

		fmt.Printf("%s ", tokens[i].Val)
	}
	fmt.Printf("\n")
	return nil
}

func ParseReturn2(data []byte) error {
	tokens, err := lexer.Scan(data)
	if err != nil {
		return errors.Wrap(err, "lexing error")
	}

	err = parseStatement(tokens)
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	return nil

}
