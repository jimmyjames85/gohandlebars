package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jimmyjames85/gohandlebars/lexer"
	"github.com/pkg/errors"
)

// TODO What is recursive descent parsing?

// func NewStatement(tokens []*lexer.Token) (*Return, int) {
// 	if len(tokens) < 3 {
// 		return nil, -1
// 	}
//
// 	expected := []lexer.TokenType{lexer.Return, lexer.IntLiteral, lexer.Semicolon}
// 	for i, t := range expected {
// 		if tokens[i].Typ != t {
// 			return nil, -2
// 		}
// 	}
//
// 	v, err := strconv.Atoi(tokens[1].Val)
// 	if err != nil {
// 		return nil, -3
// 	}
//
// 	return &Return{returnCode: v}, 3
// }

type MainFunc struct {
	ret *Return
}

func NewMainFunc(ret *Return) *MainFunc {
	return &MainFunc{ret: ret}
}

type Return struct {
	returnCode int
}

func NewReturn(returnCode int) *Return {
	return &Return{returnCode: returnCode}
}

func (r *Return) ToASM() string {
	return fmt.Sprintf("movl $%d, %%eax\nret\n", r.returnCode)
}

func (m *MainFunc) ToASM() string {

	returnASM := strings.Split(m.ret.ToASM(), "\n")
	for i := 0; i < len(returnASM); i++ {
		returnASM[i] = fmt.Sprintf("\t%s", returnASM[i])
	}

	ret := strings.Join(returnASM, "\n")
	return fmt.Sprintf(".globl _%s\n_%s:\n%s\n", "main", "main", ret)
}

type node struct {
	eval string
	root *lexer.Token
	// children []*Node
}

// should include the return statement
func NewStatementAST(val *lexer.Token) (*node, error) {

	if val.Typ != lexer.IntLiteral {
		return nil, fmt.Errorf("expeting int literal")
	}

	v, err := strconv.Atoi(val.Val)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse integer: %q", val.Val)
	}

	return &node{
		root: val,
		eval: fmt.Sprintf("%d", v),
	}, nil
}

func (n *node) ToASM() {

}

func (n *node) Eval() string {
	return n.eval
}
