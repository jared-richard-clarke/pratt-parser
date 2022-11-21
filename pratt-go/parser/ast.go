package parser

import "github/jared-richard-clarke/pratt/lexer"

type Node interface {
	ast()
}

type Number struct {
	Value float64
}

type Unary struct {
	Op lexer.LexType
	X  Node
}

type Binary struct {
	Op lexer.LexType
	X  Node
	Y  Node
}

type Ident struct {
	Value string
}

func (n *Number) ast() {}
func (u *Unary) ast()  {}
func (b *Binary) ast() {}
func (i *Ident) ast()  {}
