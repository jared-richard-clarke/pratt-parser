package parser

import (
	"fmt"
)

type Node interface {
	ast()
}

type Number struct {
	Value        float64
	Line, Column int
}

func (n *Number) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type Ident struct {
	Value        string
	Line, Column int
}

func (i *Ident) String() string {
	return i.Value
}

type Unary struct {
	Op           string
	X            Node
	Line, Column int
}

func (u *Unary) String() string {
	s := "( Op: %s, X: %s )"
	return fmt.Sprintf(s, u.Op, u.X)
}

type Binary struct {
	Op           string
	X, Y         Node
	Line, Column int
}

func (b *Binary) String() string {
	s := "( Op: %s, X: %s, Y: %s )"
	return fmt.Sprintf(s, b.Op, b.X, b.Y)
}

type Function struct {
	Name         Ident  // Function names limited to identifiers.
	Args         []Node // All Nodes are expressions.
	Line, Column int
}

func (f *Function) String() string {
	s := "( Name: %v, Args: %v )"
	return fmt.Sprintf(s, f.Name, f.Args)
}

func (l *Number) ast()   {}
func (i *Ident) ast()    {}
func (u *Unary) ast()    {}
func (b *Binary) ast()   {}
func (f *Function) ast() {}
