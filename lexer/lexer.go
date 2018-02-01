package lexer

import (
	"fmt"
	"regexp"
)

type TokenType int

// tokens are terminal statements
const (
	// Tokens
	tokenLineComment TokenType = iota
	tokenMultiLineComment
	tokenOpenBrace
	tokenCloseBrace
	tokenOpenParen
	tokenCloseParen
	tokenEquals
	tokenIsEquals
	tokenKWIf
	tokenKWElse
	tokenKWInt
	tokenKWReturn
	tokenIdentifier
	tokenIntLiteral
	tokenSemicolon
	tokenAsterix
	tokenDot
	tokenWhitespace
	tokenStringLiteral
	tokenLT
	tokenLTE
	tokenGT
	tokenGTE
)

var tokenMap = map[TokenType]*regexp.Regexp{

	// TODO optimization test: add a $ in front each one of these
	// regexes for early exit. We enforce that we must match at
	// the first char of data, so do an optimization test to see
	// if regexes do early exits.

	// See https://golang.org/pkg/regexp/syntax/
	tokenLineComment:      regexp.MustCompile(`//.*`),
	tokenMultiLineComment: regexp.MustCompile(`(?Us)/\*.*\*/`), // (?Us) Ungreedy multi-line mode
	tokenOpenBrace:        regexp.MustCompile(`{`),
	tokenCloseBrace:       regexp.MustCompile(`}`),
	tokenOpenParen:        regexp.MustCompile(`\(`),
	tokenCloseParen:       regexp.MustCompile(`\)`),
	tokenEquals:           regexp.MustCompile(`=`),
	tokenIsEquals:         regexp.MustCompile(`==`),
	tokenKWIf:             regexp.MustCompile(`if`),
	tokenKWElse:           regexp.MustCompile(`else`),
	tokenKWInt:            regexp.MustCompile(`int`),
	tokenKWReturn:         regexp.MustCompile(`return`),
	tokenIdentifier:       regexp.MustCompile(`[a-zA-Z]\w*`),
	tokenIntLiteral:       regexp.MustCompile(`[0-9]+`),
	tokenSemicolon:        regexp.MustCompile(`;`),
	tokenAsterix:          regexp.MustCompile(`\*`),
	tokenDot:              regexp.MustCompile(`\.`),
	tokenWhitespace:       regexp.MustCompile(`\s+`),
	tokenStringLiteral:    regexp.MustCompile(`"(?:[^"\\]|\\.)*"`), // this probably needs to be tested
	tokenLT:               regexp.MustCompile(`<`),
	tokenLTE:              regexp.MustCompile(`<=`),
	tokenGT:               regexp.MustCompile(`>`),
	tokenGTE:              regexp.MustCompile(`>=`),
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
	if typ == tokenWhitespace {
		// whitspace found, but don't return it
		return cur[1], nil, nil
	}

	tkn := &Token{
		Typ: typ,
		Val: string(data[cur[0]:cur[1]]),
	}
	return cur[1], tkn, nil
}
