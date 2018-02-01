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

	for i, tkn := range tokens {
		fmt.Printf("===============%02d:[%02d]===============\n%s\n", i, tkn.Typ, tkn.Val)
	}

	return nil

}
