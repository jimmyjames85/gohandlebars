package parser

import (
	"fmt"
	"strconv"

	"github.com/jimmyjames85/gohandlebars/pkg/ast"
	"github.com/jimmyjames85/gohandlebars/pkg/lexer"
	"github.com/pkg/errors"
)

type NonTerminal interface{}

func Parse(data []byte) error {

	tokens, err := lexer.Scan(data)
	if err != nil {
		return errors.Wrap(err, "lexing error")
	}

	var tkn *lexer.Token
	for tkn = tokens.Pop(); tkn != nil; tkn = tokens.Pop() {
		fmt.Printf("%s: %s\n\n", tkn.Typ, tkn.Val)
	}

	return nil
}

var ErrNotEnoughTokens = fmt.Errorf("not enough tokens")

func ParseFactor(l lexer.TokenList) (ast.Exp, error) {
	t := l.Pop()
	if t == nil {
		return nil, ErrNotEnoughTokens
	}

	switch t.Typ {
	case lexer.OpenParen:
		ret, err := ParseExp(l)
		if err != nil {
			return nil, err
		}
		nextParen := l.Pop()
		if nextParen == nil {
			return nil, ErrNotEnoughTokens
		}
		if nextParen.Typ != lexer.CloseParen {
			return nil, fmt.Errorf("mismatch parens: expecting ')': got %s", nextParen.Typ)
		}
		return ret, nil
	case lexer.Negation, lexer.LogicalNegation, lexer.BitwiseCompliment:
		nextExp, err := ParseFactor(l)
		if err != nil {
			return nil, err
		}
		ret := ast.NewUnOp(t.Typ, nextExp)
		return ret, nil
	case lexer.IntLiteral:
		n, err := strconv.Atoi(t.Val)
		if err != nil {
			return nil, err
		}
		ret := ast.NewConstant(n)
		return ret, nil
	default:
		return nil, fmt.Errorf("parseFactor: unexpected type: %s", t.Typ)
	}

	return nil, fmt.Errorf("unexpected error")
}

func ParseTerm(l lexer.TokenList) (ast.Exp, error) {
	ret, err := ParseFactor(l)
	if err != nil {
		return nil, err
	}

	t := l.Peek()
	for t != nil && (t.Typ == lexer.Multiplication || t.Typ == lexer.Division) {
		l.Pop()
		nextFactor, err := ParseFactor(l)
		if err != nil {
			return nil, err
		}
		ret = ast.NewBinOp(t.Typ, ret, nextFactor)
		t = l.Peek()
	}
	return ret, nil
}

func ParseExp(l lexer.TokenList) (ast.Exp, error) {

	ret, err := ParseTerm(l)
	if err != nil {
		return nil, err
	}

	t := l.Peek()
	for t != nil && (t.Typ == lexer.Addition || t.Typ == lexer.Negation) {
		l.Pop()
		nextTerm, err := ParseTerm(l)
		if err != nil {
			return nil, err
		}
		ret = ast.NewBinOp(t.Typ, ret, nextTerm)
		t = l.Peek()
	}
	return ret, nil
}

func ParseStatement(l lexer.TokenList) (*ast.Statement, error) {
	t := l.Pop()
	if t == nil {
		return nil, ErrNotEnoughTokens
	}
	if t.Typ != lexer.Return {
		return nil, fmt.Errorf("expecting return: got %s: %q", t.Typ, t.Val)
	}

	exp, err := ParseExp(l)
	if err != nil {
		return nil, err
	}

	t = l.Pop()
	if t == nil {
		return nil, ErrNotEnoughTokens
	}
	if t.Typ != lexer.Semicolon {
		return nil, fmt.Errorf("missing semicolon: got %s: %q", t.Typ, t.Val)
	}
	return ast.NewStatement(exp), nil
}

func ParseFunction(l lexer.TokenList) (*ast.Function, error) {

	expected := []lexer.TokenType{
		lexer.Int,
		lexer.Identifier,
		lexer.OpenParen,
		lexer.CloseParen,
		lexer.OpenBrace,
	}

	var identifier *lexer.Token

	for _, e := range expected {
		t := l.Pop()
		if t == nil {
			return nil, ErrNotEnoughTokens
		}
		if t.Typ != e {
			return nil, fmt.Errorf("expecting %s: got %s: %q", e, t.Typ, t.Val)
		}
		if t.Typ == lexer.Identifier {
			identifier = t
		}
	}

	stmt, err := ParseStatement(l)
	if err != nil {
		return nil, err
	}

	t := l.Pop()
	if t == nil {
		return nil, ErrNotEnoughTokens
	} else if t.Typ != lexer.CloseBrace {
		return nil, fmt.Errorf("unbalanced curly braces")
	}

	return ast.NewFunction(identifier.Val, stmt), nil
}

func parseMainFunc(l lexer.TokenList) (int, *ast.MainFunc, error) {

	// // keepgoing := true
	// for tkn, err := l.Peek(); err == nil; tkn, err = l.Peek() {
	// 	tkn, err = l.Pop()
	// 	if err != nil {
	// 		return 0, nil, err
	// 	}
	// 	fmt.Printf("%s", tkn)
	// }
	// // tokn, err := l.Peek(); true {
	// switch {
	// }

	//////////////////////////////////////////////////////////////////////
	if l.Size() < 5 {
		return 0, nil, ErrNotEnoughTokens
	}

	expected := []lexer.TokenType{
		lexer.Int,
		lexer.Identifier, //TODO verify identifier is main
		lexer.OpenParen,
		lexer.CloseParen,
		lexer.OpenBrace,
	}

	tokens := l.TODO_ToSlice()
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

func ParseReturn2(data []byte) ([]byte, error) {
	tokens, err := lexer.Scan(data)
	if err != nil {
		return nil, errors.Wrap(err, "lexing error")
	}

	_, mainFunc, err := parseMainFunc(tokens)
	if err != nil {
		return nil, errors.Wrap(err, "parse")
	}
	return mainFunc.ToASM(), nil
}
