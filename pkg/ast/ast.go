package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jimmyjames85/gohandlebars/pkg/lexer"
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

////////////////////////////// Expresion //////////////////////////////

// <program> ::= <function>
// <function> ::= "int" <id> "(" ")" "{" <statement> "}"
// <statement> ::= "return" <exp> ";"

// <exp> ::= <term> { ("+" | "-") <term> }
// <term> ::= <factor> { ("*" | "/") <factor> }
// <factor> ::= "(" <exp> ")" | <unary_op> <factor> | <int>

// <unary_op> ::= "!" | "~" | "-"
// <binary_op> ::= "+" | "-" | "*" | "/"

type Exp interface {
	ExpString() string
}

type Constant struct {
	c int
}

func NewConstant(c int) *Constant {
	return &Constant{c: c}
}

func (c *Constant) ExpString() string {
	return fmt.Sprintf("%d", c.c)
}

type UnOp struct {
	// A unary operator should only be applied to a whole expression if:
	//
	//  - the expression is a single integer (e.g. ~4)
	//  - the expression is wrapped in parentheses (e.g. ~(1+1)), or
	//  - the expression is itself a unary operation (e.g. ~!8, -~(2+2)).

	op  lexer.TokenType
	exp Exp
}

func NewUnOp(op lexer.TokenType, exp Exp) *UnOp {
	// TODO verify op is valid e.g. ! - or ~
	return &UnOp{op: op, exp: exp}
}

func (u *UnOp) ExpString() string {
	return fmt.Sprintf("(%s %s)", u.op, u.exp.ExpString())
}

type BinOp struct {
	op   lexer.TokenType
	lExp Exp
	rExp Exp
}

func NewBinOp(op lexer.TokenType, lExp, rExp Exp) *BinOp {
	// TODO verify op is valid e.g. + or -
	return &BinOp{op: op, lExp: lExp, rExp: rExp}
}

func (b *BinOp) ExpString() string {
	return fmt.Sprintf("(%s %s %s)", b.lExp.ExpString(), b.op, b.rExp.ExpString())
}

///////////////////////////// Statement /////////////////////////////
type Statement struct {
	exp Exp
}

func NewStatement(exp Exp) *Statement {
	return &Statement{exp: exp}
}

func (s *Statement) StmtString() string {
	return fmt.Sprintf("return %s", s.exp.ExpString())
}

///////////////////////////// Statement /////////////////////////////
// <function> ::= "int" <id> "(" ")" "{" <statement> "}"

type Function struct {
	name string
	stmt *Statement
}

func NewFunction(name string, stmt *Statement) *Function {
	return &Function{name: name, stmt: stmt}
}

func (f *Function) FuncString() string {
	return fmt.Sprintf("int %s () { %s }", f.name, f.stmt.StmtString())
}

//////////////////////////////////////////////////////////////////////

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

func (m *MainFunc) ToASM() []byte {

	returnASM := strings.Split(m.ret.ToASM(), "\n")
	for i := 0; i < len(returnASM); i++ {
		returnASM[i] = fmt.Sprintf("\t%s", returnASM[i])
	}

	ret := strings.Join(returnASM, "\n")
	return []byte(fmt.Sprintf(".globl _%s\n_%s:\n%s\n", "main", "main", ret))
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
