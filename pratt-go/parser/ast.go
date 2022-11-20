package parser

import (
	"fmt"
	"github/jared-richard-clarke/pratt/lexer"
)

type position struct {
	line   int
	column int
	length int
}

type Node interface {
	Pos() string
}

type Number struct {
	Value float64
	position
}

func (n *Number) Pos() string {
	return fmt.Sprintf("lin %d, col %d, len %d", n.line, n.column, n.length)
}

type Unary struct {
	Op lexer.LexType
	X  Node
	position
}

func (u *Unary) Pos() string {
	return fmt.Sprintf("lin %d, col %d, len %d", u.line, u.column, u.length)
}

type Binary struct {
	Op lexer.LexType
	X  Node
	Y  Node
	position
}

func (b *Binary) Pos() string {
	return fmt.Sprintf("lin %d, col %d, len %d", b.line, b.column, b.length)
}

type Ident struct {
	Id string
	position
}

func (i *Ident) Pos() string {
	return fmt.Sprintf("lin %d, col %d, len %d", i.line, i.column, i.length)
}
