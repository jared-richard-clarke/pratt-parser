package parser

import "github/jared-richard-clarke/pratt/internal/lexer"

type Node interface {
	ast()
}

type Literal struct {
	Typeof       lexer.LexType
	Value        string
	Line, Column int
}

func (l *Literal) String() string {
	s := "Literal{ Value: %s }"
	return fmt.Sprintf(s, l.Value)
}

type Unary struct {
	Op           string
	X            Node
	Line, Column int
}

func (u *Unary) String() string {
	s := "Unary{ Op: %s, X: %s }"
	return fmt.Sprintf(s, u.Op, u.X)
}

type Binary struct {
	Op           string
	X, Y         Node
	Line, Column int
}

func (b *Binary) String() string {
	s := "Binary{ Op: %s, X: %s, Y: %s }"
	return fmt.Sprintf(s, b.Op, b.X, b.Y)
}

func (l *Literal) ast() {}
func (u *Unary) ast()   {}
func (b *Binary) ast()  {}
