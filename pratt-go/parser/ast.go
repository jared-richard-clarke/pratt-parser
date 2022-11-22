package parser

import "github/jared-richard-clarke/pratt/lexer"

type Node interface {
	ast()
}

type Number struct {
	Value        float64
	Line, Column int
}

type Ident struct {
	Value        string
	Line, Column int
}

type Unary struct {
	Op           lexer.LexType
	X            Node
	Line, Column int
}

type Binary struct {
	Op           lexer.LexType
	X, Y         Node
	Line, Column int
}

func (n *Number) ast() {}
func (u *Unary) ast()  {}
func (b *Binary) ast() {}
func (i *Ident) ast()  {}
