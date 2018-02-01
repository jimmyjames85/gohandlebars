package lexer

import (
	"fmt"
	"regexp"
)

type TokenType string

func (t *TokenType) String() string { return string(*t) }

const (
	LineComment      = TokenType("LineComment")
	MultiLineComment = TokenType("MultiLineComment")

	OpenBrace  = TokenType("OpenBrace")
	CloseBrace = "CloseBrace"
	OpenParen  = TokenType("OpenParen")
	CloseParen = "CloseParen"

	// comparators
	EQ  = TokenType("EQ")
	LT  = TokenType("LT")
	LTE = TokenType("LTE")
	GT  = TokenType("GT")
	GTE = TokenType("GTE")

	// keywords
	If     = TokenType("if")
	Else   = TokenType("else")
	Int    = TokenType("int")
	Float  = TokenType("float")
	Return = TokenType("return")

	// literals todo is literal the right name?
	IntLiteral    = TokenType("IntLiteral")
	FloatLiteral  = TokenType("FloatLiteral")
	StringLiteral = TokenType("StringLiteral")

	// todo better naming for the following...and what classification?
	Assignment = TokenType("Assignment")
	Asterix    = TokenType("Asterix")
	Dot        = TokenType("Dot")
	Identifier = TokenType("Identifier")
	Semicolon  = TokenType("Semicolon")
	Whitespace = TokenType("Whitespace")
)

// tokens are terminal statements
var tokenMap = map[TokenType]*regexp.Regexp{

	// TODO optimization test: add a ^ in front each one of these
	// regexes for early exit. We enforce that we must match at
	// the first char of data, so do an optimization test to see
	// if regexes do early exits.

	LineComment:      regexp.MustCompile(`//.*`),
	MultiLineComment: regexp.MustCompile(`(?Us)/\*.*\*/`), // (?Us) Ungreedy multi-line mode See https://golang.org/pkg/regexp/syntax/

	OpenBrace:  regexp.MustCompile(`{`),
	CloseBrace: regexp.MustCompile(`}`),
	OpenParen:  regexp.MustCompile(`\(`),
	CloseParen: regexp.MustCompile(`\)`),

	// comparators
	EQ:  regexp.MustCompile(`==`),
	LT:  regexp.MustCompile(`<`),
	LTE: regexp.MustCompile(`<=`),
	GT:  regexp.MustCompile(`>`),
	GTE: regexp.MustCompile(`>=`),

	// keywords
	If:     regexp.MustCompile(`if`),
	Else:   regexp.MustCompile(`else`),
	Int:    regexp.MustCompile(`int`),
	Float:  regexp.MustCompile(`float`),
	Return: regexp.MustCompile(`return`),

	// literals todo is literal the right name?
	IntLiteral:    regexp.MustCompile(`[0-9]+`),
	FloatLiteral:  regexp.MustCompile(`[0-9]+\.[0-9]+`),
	StringLiteral: regexp.MustCompile(`"(?:[^"\\]|\\.)*"`), // this probably needs to be tested

	// todo better naming for the following...and what classification?
	Assignment: regexp.MustCompile(`=`),
	Asterix:    regexp.MustCompile(`\*`),
	Dot:        regexp.MustCompile(`\.`),
	Identifier: regexp.MustCompile(`[a-zA-Z]\w*`), // TODO need to make sure these don't clash with keywords
	Semicolon:  regexp.MustCompile(`;`),
	Whitespace: regexp.MustCompile(`\s+`),
}

type Token struct {
	Typ TokenType
	Val string
}

func Scan(data []byte) ([]*Token, error) {
	var ret []*Token

	buf := data
	advance, token, err := nextToken(buf)
	for advance > 0 {
		if err != nil {
			return nil, err
		}
		if token != nil {
			ret = append(ret, token)
		}
		buf = buf[advance:]
		advance, token, err = nextToken(buf)
	}

	if err != nil {
		return nil, err
	}
	return ret, nil
}

func nextToken(data []byte) (advance int, token *Token, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	// TODO make this run all regexes at once with goroutines
	var cur []int
	var typ TokenType
	for t, re := range tokenMap {

		i := re.FindIndex(data)
		if i == nil {
			// doesn't match
			continue
		}

		tokenLen := i[1] - i[0]
		if tokenLen == 0 {
			// matches but doesn't give a valid token (prob bad regex)
			return 0, nil, fmt.Errorf("token type %d: returned empty string as token", t)
		}

		var curLen int
		if cur != nil {
			curLen = cur[1] - cur[0]
		}

		if cur == nil || i[0] < cur[0] {
			// if cur hasn't been assigned or matches
			// with an index (i[0]) closer then cur[0]
			cur = i
			typ = t
		} else if cur != nil && (i[0] == cur[0] && tokenLen > curLen) {
			// if we matched at the same spot as
			// cur[0] but we can consume more then
			// accept this token ( e.g. accept <= instead of < )
			cur = i
			typ = t
		}
	}

	if cur == nil {
		// no matches
		return 0, nil, fmt.Errorf("unexpected token: %q", string(data))
	} else if cur[0] != 0 {
		// Prevent matches that don't consume data at the very
		// first char. We don't want to skip unrecognized
		// tokens. TODO see todo above and add ^ to regexes to
		// see if it improves prefomance
		return 0, nil, fmt.Errorf("unexpected token: %q", data[0:cur[0]+1])
	}

	// one of the tokens matched!!
	if typ == "Whitespace" {
		// whitspace found, but don't return it
		return cur[1], nil, nil
	}

	tkn := &Token{
		Typ: typ,
		Val: string(data[cur[0]:cur[1]]),
	}
	return cur[1], tkn, nil
}
