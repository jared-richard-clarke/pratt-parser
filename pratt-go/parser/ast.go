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
	msg := "Number{ Value: %f, Line: %d, Column: %d }"
	return fmt.Sprintf(msg, n.Value, n.Line, n.Column)
}

type Symbol struct {
	Value        string
	Line, Column int
}

func (s *Symbol) String() string {
	msg := "Symbol{ Value: %s, Line: %d, Column: %d }"
	return fmt.Sprintf(msg, s.Value, s.Line, s.Column)
}

type Unary struct {
	Op           string
	X            Node
	Line, Column int
}

func (u *Unary) String() string {
	msg := "Unary{ Op: %s, X: %s, Line: %d, Column: %d }"
	return fmt.Sprintf(msg, u.Op, u.X, u.Line, u.Column)
}

type Binary struct {
	Op           string
	X, Y         Node
	Line, Column int
}

func (b *Binary) String() string {
	msg := "Binary{ Op: %s, X: %s, Y: %s, Line: %d, Column: %d }"
	return fmt.Sprintf(msg, b.Op, b.X, b.Y, b.Line, b.Column)
}

// Like Binary but without positional information.
type ImpliedBinary struct {
	Op   string
	X, Y Node
}

func (i *ImpliedBinary) String() string {
	msg := "ImpliedBinary{ Op: %s, X: %s, Y: %s }"
	return fmt.Sprintf(msg, i.Op, i.X, i.Y)
}

type Call struct {
	Callee       Node
	Args         []Node
	Line, Column int
}

func (c *Call) String() string {
	msg := "Call{ Callee: %v, Args: %v, Line: %d, Column: %d }"
	return fmt.Sprintf(msg, c.Callee, c.Args, c.Line, c.Column)
}

func (e *Empty) ast()         {}
func (n *Number) ast()        {}
func (s *Symbol) ast()        {}
func (u *Unary) ast()         {}
func (b *Binary) ast()        {}
func (i *ImpliedBinary) ast() {}
func (c *Call) ast()          {}
