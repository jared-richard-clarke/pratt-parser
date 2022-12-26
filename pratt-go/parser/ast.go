package parser

import (
	"fmt"
)

type Node interface {
	ast()
}

type Empty struct{}

func (e *Empty) String() string {
	return "Empty{}"
}

type Number struct {
	Value        float64
	Line, Column int
}

func (n *Number) String() string {
	msg := "Number{ Value: %f }"
	return fmt.Sprintf(msg, n.Value)
}

type Ident struct {
	Value        string
	Line, Column int
}

func (i *Ident) String() string {
	msg := "Ident{ Value: %s }"
	return fmt.Sprintf(msg, i.Value)
}

type Unary struct {
	Op           string
	X            Node
	Line, Column int
}

func (u *Unary) String() string {
	msg := "Unary{ Op: %s, X: %s }"
	return fmt.Sprintf(msg, u.Op, u.X)
}

type Binary struct {
	Op           string
	X, Y         Node
	Line, Column int
}

func (b *Binary) String() string {
	msg := "Binary{ Op: %s, X: %s, Y: %s }"
	return fmt.Sprintf(msg, b.Op, b.X, b.Y)
}

type Function struct {
	Name         Ident  // Function names limited to identifiers.
	Args         []Node // All Nodes are expressions.
	Line, Column int
}

func (f *Function) String() string {
	msg := "Function{ Name: %v, Args: %v }"
	return fmt.Sprintf(msg, f.Name, f.Args)
}

func (e *Empty) ast()    {}
func (l *Number) ast()   {}
func (i *Ident) ast()    {}
func (u *Unary) ast()    {}
func (b *Binary) ast()   {}
func (f *Function) ast() {}
