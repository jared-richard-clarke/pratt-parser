package parser

import "github/jared-richard-clarke/pratt/lexer"

type Node interface {
	ast()
}

type Literal struct {
	Typeof       lexer.LexType
	Value        string
	Line, Column int
}

type Unary struct {
	Op           string
	X            Node
	Line, Column int
}

type Binary struct {
	Op           string
	X, Y         Node
	Line, Column int
}

func (l *Literal) ast() {}
func (u *Unary) ast()   {}
func (b *Binary) ast()  {}
