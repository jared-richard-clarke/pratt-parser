package parser

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
	Op           string
	X            Node
	Line, Column int
}

type Binary struct {
	Op           string
	X, Y         Node
	Line, Column int
}

func (n *Number) ast() {}
func (i *Ident) ast()  {}
func (u *Unary) ast()  {}
func (b *Binary) ast() {}
