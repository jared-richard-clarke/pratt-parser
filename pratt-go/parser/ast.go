package parser

import (
	"fmt"
)

// The interface that all AST components must satisfy.
type Node interface {
	ast()
}

// An empty string creates an empty Node.
type Empty struct{}

func (e *Empty) String() string {
	return "Empty{}"
}

// Numbers parsed as 64-bit floating point.
type Number struct {
	Value        float64
	Line, Column int
}

func (n *Number) String() string {
	msg := "Number{ Value: %g }"
	return fmt.Sprintf(msg, n.Value)
}

// Any symbolic stand in for a value, function, or expression.
// Also known as an identifier.
type Symbol struct {
	Value        string
	Line, Column int
}

func (s *Symbol) String() string {
	msg := "Symbol{ Value: %q }"
	return fmt.Sprintf(msg, s.Value)
}

// Operations with one operand.
type Unary struct {
	Op           string
	X            Node
	Line, Column int
}

func (u *Unary) String() string {
	msg := "Unary{ Op: %q, X: %s }"
	return fmt.Sprintf(msg, u.Op, u.X)
}

// Operations with two operands.
type Binary struct {
	Op           string
	X, Y         Node
	Line, Column int
}

func (b *Binary) String() string {
	msg := "Binary{ Op: %q, X: %s, Y: %s }"
	return fmt.Sprintf(msg, b.Op, b.X, b.Y)
}

// Like Binary but without positional information.
type ImpliedBinary struct {
	Op   string
	X, Y Node
}

func (i *ImpliedBinary) String() string {
	msg := "ImpliedBinary{ Op: %q, X: %s, Y: %s }"
	return fmt.Sprintf(msg, i.Op, i.X, i.Y)
}

// A function call.
type Call struct {
	Callee       Node
	Args         []Node
	Line, Column int
}

func (c *Call) String() string {
	msg := "Call{ Callee: %v, Args: %v }"
	return fmt.Sprintf(msg, c.Callee, c.Args)
}

func (e *Empty) ast()         {}
func (n *Number) ast()        {}
func (s *Symbol) ast()        {}
func (u *Unary) ast()         {}
func (b *Binary) ast()        {}
func (i *ImpliedBinary) ast() {}
func (c *Call) ast()          {}
