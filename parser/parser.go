package parser

import (
	"fmt"
	"strconv"

	"github.com/jimmyjames85/gohandlebars/ast"
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

var ErrNotEnoughTokens = fmt.Errorf("not enough tokens")

func parseMainFunc(tokens []*lexer.Token) (int, *ast.MainFunc, error) {

	if len(tokens) < 5 {
		return 0, nil, ErrNotEnoughTokens
	}

	expected := []lexer.TokenType{
		lexer.Int,
		lexer.Identifier, //TODO verify identifier is main
		lexer.OpenParen,
		lexer.CloseParen,
		lexer.OpenBrace,
	}
	got := make([]*lexer.Token, 0)

	var i int
	for i = 0; i < len(expected) && i < len(tokens); i++ {
		t := tokens[i]

		if t.Typ == expected[i] {
			got = append(got, t)
		} else {
			return 0, nil, fmt.Errorf("expecting %s: got %s: %q", expected[i], t.Typ, t.Val)
		}
	}

	if len(got) != len(expected) {
		return 0, nil, fmt.Errorf("mismatch expected: TODO")
	}

	advance, ret, err := parseStatement(tokens[len(expected):])
	if err != nil {
		return 0, nil, err
	}

	i += advance
	if tokens[i].Typ != lexer.CloseBrace {
		return 0, nil, fmt.Errorf("expecting %s: got %s: %q", lexer.CloseBrace, tokens[i].Typ, tokens[i].Val)
	}
	i++

	return (i + advance), ast.NewMainFunc(ret), nil
}

func parseStatement(tokens []*lexer.Token) (int, *ast.Return, error) {
	//right now stmt is "return 3;"
	if len(tokens) < 3 {
		return 0, nil, ErrNotEnoughTokens
	}

	expected := []lexer.TokenType{lexer.Return, lexer.IntLiteral, lexer.Semicolon}

	var advance, count, i int

	var intLiteral *lexer.Token

	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		if count >= len(expected) {
			break
		}
		advance++

		if t.Typ == expected[count] {
			if t.Typ == lexer.IntLiteral {
				intLiteral = t
			}
			count++
		} else {
			return 0, nil, fmt.Errorf("ps: expecting %s: got %s: %q", expected[count], t.Typ, t.Val)
		}

	}

	v, err := strconv.Atoi(intLiteral.Val)
	if err != nil {
		return 0, nil, fmt.Errorf("unable to parse int: %s", intLiteral.Val)
	}

	return advance, ast.NewReturn(v), nil
}

func ParseReturn2(data []byte) error {
	tokens, err := lexer.Scan(data)
	if err != nil {
		return errors.Wrap(err, "lexing error")
	}

	_, mainFunc, err := parseMainFunc(tokens)
	if err != nil {
		return errors.Wrap(err, "parse")
	}
	fmt.Printf("%s\n", mainFunc.ToASM())

	// advance, ret, err := parseStatement(tokens)
	// if err != nil {
	// 	return errors.Wrap(err, "parse")
	// }
	// tokens = tokens[advance+1:]
	// fmt.Printf("%s\n", ret.ToASM())

	return nil
}
