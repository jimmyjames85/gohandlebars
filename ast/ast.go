package ast

// TODO What is recursive descent parsing?

type Expression int

type Statement struct{}

type FunctionDeclaration struct{}

func NewFunctionDeclaration(name string, stmt Statement) (*FunctionDeclaration, error) {
	return nil, nil
}

type Program struct{}

func NewProgram(fd FunctionDeclaration) (*Program, error) {
	return nil, nil
}
