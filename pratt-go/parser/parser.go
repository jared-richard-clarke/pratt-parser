package parser

import "github/jared-richard-clarke/pratt/lexer"

type Expr interface {
	Nud() (Expr, error)     // Null Denotation
	Led(Expr) (Expr, error) // Left Denotation
}

type Number struct {
	Value   string
	Literal float64
}

type Unary struct {
	Op lexer.LexType
	X  Expr
}

type Binary struct {
	Op   lexer.LexType
	X, Y Expr
}

type parser struct {
	src   []lexer.Token // Token source.
	token lexer.Token   // Current token.
	index int           // Current index in src.
	ast   Expr          // Abstract Syntax Tree.
}
