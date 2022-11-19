package parser

import (
	"fmt"
	"pratt/lexer"
	"strconv"
)

type Expr interface {
	Nud() (Expr, error)     // Null Denotation
	Led(Expr) (Expr, error) // Left Denotation
}

type NumExpr struct {
	Lexeme string
	Value  float64
}

func (n *NumExpr) Nud() (Expr, error) {
	num, err := strconv.ParseFloat(n.Lexeme, 64)
	if err != nil {
		return n, err
	}
	n.Value = num
	return n, err
}
func (n *NumExpr) Led(e Expr) (Expr, error) {
	err := fmt.Errorf("expected operator, got number %s", n.Lexeme)
	return n, err
}

type IdentExpr struct {
	Lexeme string
	Value  Expr
}

type UnaryExpr struct {
	Op lexer.LexType
	X  Expr
}

func AddPrefix(x Expr) UnaryExpr {
	return UnaryExpr{
		Op: lexer.Add,
		X:  x,
	}
}

func SubPrefix(x Expr) UnaryExpr {
	return UnaryExpr{
		Op: lexer.Sub,
		X:  x,
	}
}

type BinaryExpr struct {
	Op   lexer.LexType
	X, Y Expr
	Lbp  int
}

func AddInfix(x, y Expr) BinaryExpr {
	return BinaryExpr{
		Op:  lexer.Add,
		X:   x,
		Y:   y,
		Lbp: 50,
	}
}

func SubInfix(x, y Expr) BinaryExpr {
	return BinaryExpr{
		Op:  lexer.Sub,
		X:   x,
		Y:   y,
		Lbp: 50,
	}
}

func MulInfix(x, y Expr) BinaryExpr {
	return BinaryExpr{
		Op:  lexer.Mul,
		X:   x,
		Y:   y,
		Lbp: 60,
	}
}

func DivInfix(x, y Expr) BinaryExpr {
	return BinaryExpr{
		Op:  lexer.Div,
		X:   x,
		Y:   y,
		Lbp: 60,
	}
}

func PowInfix(x, y Expr) BinaryExpr {
	return BinaryExpr{
		Op:  lexer.Pow,
		X:   x,
		Y:   y,
		Lbp: 70,
	}
}