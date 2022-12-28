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

type Symbol struct {
	Value        string
	Line, Column int
}

func (i *Symbol) String() string {
	msg := "Symbol{ Value: %s }"
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

// Like Binary but without positional information.
type ImpBinary struct {
	Op   string
	X, Y Node
}

func (i *ImpBinary) String() string {
	msg := "ImpBinary{ Op: %s, X: %s, Y: %s }"
	return fmt.Sprintf(msg, i.Op, i.X, i.Y)
}

type Call struct {
	Fun          Node
	Args         []Node
	Line, Column int
}

func (c *Call) String() string {
	msg := "Call{ Fun: %v, Args: %v }"
	return fmt.Sprintf(msg, c.Fun, c.Args)
}

func (e *Empty) ast()     {}
func (n *Number) ast()    {}
func (s *Symbol) ast()     {}
func (u *Unary) ast()     {}
func (b *Binary) ast()    {}
func (i *ImpBinary) ast() {}
func (f *Call) ast()      {}
