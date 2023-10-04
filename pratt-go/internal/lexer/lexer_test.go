package lexer

import (
	"fmt"
	"testing"
)

func (t Token) String() string {
	switch {
	case t.Typeof == ImpMul:
		return fmt.Sprintf("punct: \"imp-*\" :%d:%d", t.Line, t.Column)
	case t.Typeof < Number:
		return fmt.Sprintf("punct: %q :%d:%d", t.Value, t.Line, t.Column)
	case t.Typeof == Number:
		return fmt.Sprintf("number: %q :%d:%d", t.Value, t.Line, t.Column)
	case t.Typeof == Symbol:
		return fmt.Sprintf("symbol: %q :%d:%d", t.Value, t.Line, t.Column)
	default:
		return fmt.Sprintf("<eof> :%d:%d", t.Line, t.Column)
	}
}

func TestScan(t *testing.T) {
	text := "1 + 2 * 3"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "1",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 5,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   1,
			Column: 7,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   1,
			Column: 9,
		},
		mkEof(1, 10),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Scan")
}

func TestEmpty(t *testing.T) {
	text := " \n\t"
	expect := []Token{mkEof(2, 1)}
	result, _ := Scan(text)
	compare(expect, result, t, "Empty")
}

func TestParens(t *testing.T) {
	text := "(1 + 2) * 3"
	expect := []Token{
		{
			Typeof: OpenParen,
			Value:  "(",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Number,
			Value:  "1",
			Line:   1,
			Column: 2,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 4,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 6,
		},
		{
			Typeof: CloseParen,
			Value:  ")",
			Line:   1,
			Column: 7,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   1,
			Column: 9,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   1,
			Column: 11,
		},
		mkEof(1, 12),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Parens")
}

func TestComma(t *testing.T) {
	text := "op(2, 5)"
	expect := []Token{
		{
			Typeof: Symbol,
			Value:  "op",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: OpenParen,
			Value:  "(",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 4,
		},
		{
			Typeof: Comma,
			Value:  ",",
			Line:   1,
			Column: 5,
		},
		{
			Typeof: Number,
			Value:  "5",
			Line:   1,
			Column: 7,
		},
		{
			Typeof: CloseParen,
			Value:  ")",
			Line:   1,
			Column: 8,
		},
		mkEof(1, 9),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Comma")
}

func TestNewlines(t *testing.T) {
	text := "1 + 2\n * 3"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "1",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 5,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   2,
			Column: 2,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   2,
			Column: 4,
		},
		mkEof(2, 5),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Newlines (1)")

	text = `1 +
	        2 *
	        3`
	expect = []Token{
		{
			Typeof: Number,
			Value:  "1",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   2,
			Column: 9,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   2,
			Column: 11,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   3,
			Column: 9,
		},
		mkEof(3, 10),
	}
	result, _ = Scan(text)
	compare(expect, result, t, "Newlines (2)")
}

func TestSymbol(t *testing.T) {
	text := "x + wyvern * 3"
	expect := []Token{
		{
			Typeof: Symbol,
			Value:  "x",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Symbol,
			Value:  "wyvern",
			Line:   1,
			Column: 5,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   1,
			Column: 12,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   1,
			Column: 14,
		},
		mkEof(1, 15),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Symbol (1)")

	text = "x + wyvern/hamster"
	expect = []Token{
		{
			Typeof: Symbol,
			Value:  "x",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Symbol,
			Value:  "wyvern",
			Line:   1,
			Column: 5,
		},
		{
			Typeof: Div,
			Value:  "/",
			Line:   1,
			Column: 11,
		},
		{
			Typeof: Symbol,
			Value:  "hamster",
			Line:   1,
			Column: 12,
		},
		mkEof(1, 19),
	}
	result, _ = Scan(text)
	compare(expect, result, t, "Symbol (2)")
}

func TestNumbers(t *testing.T) {
	text := "7.5/2"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "7.5",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Div,
			Value:  "/",
			Line:   1,
			Column: 4,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 5,
		},
		mkEof(1, 6),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Numbers (1)")

	text = "1024"
	expect = []Token{
		{
			Typeof: Number,
			Value:  "1024",
			Line:   1,
			Column: 1,
		},
		mkEof(1, 5),
	}
	result, _ = Scan(text)
	compare(expect, result, t, "Numbers (2)")
}

func TestSub(t *testing.T) {
	text := "1 - -2"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "1",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Sub,
			Value:  "-",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Sub,
			Value:  "-",
			Line:   1,
			Column: 5,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 6,
		},
		mkEof(1, 7),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Sub")
}

// A Token of type ImpMul has no meaningful position.
// These Tokens carry Line and Column fields only because
// a unique Token type is not worth the added complexity.
func TestImpMul(t *testing.T) {
	text := "2x"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: ImpMul,
			Value:  "*",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Symbol,
			Value:  "x",
			Line:   1,
			Column: 2,
		},
		mkEof(1, 3),
	}
	result, _ := Scan(text)
	compare(result, expect, t, "ImpMul (1)")

	text = "(x)y"
	expect = []Token{
		{
			Typeof: OpenParen,
			Value:  "(",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Symbol,
			Value:  "x",
			Line:   1,
			Column: 2,
		},
		{
			Typeof: CloseParen,
			Value:  ")",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: ImpMul,
			Value:  "*",
			Line:   1,
			Column: 3,
		},
		{
			Typeof: Symbol,
			Value:  "y",
			Line:   1,
			Column: 4,
		},
		mkEof(1, 5),
	}
	result, _ = Scan(text)
	compare(result, expect, t, "ImpMul (2)")
}

func TestPow(t *testing.T) {
	text := "4^2"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "4",
			Line:   1,
			Column: 1,
		},
		{
			Typeof: Pow,
			Value:  "^",
			Line:   1,
			Column: 2,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 3,
		},
		mkEof(1, 4),
	}
	result, _ := Scan(text)
	compare(expect, result, t, "Pow")
}

// utility functions

func mkEof(l, c int) Token {
	return Token{
		Typeof: EOF,
		Line:   l,
		Column: c,
	}
}

func compare(e []Token, r []Token, t *testing.T, n string) {
	if len(e) != len(r) {
		t.Errorf("Test %s failed. Token slices of unequal length.", n)
	} else {
		for i := range e {
			exp := e[i]
			got := r[i]
			if exp != got {
				t.Errorf("Test %s failed. Expected: %v, Got: %v", n, exp, got)
			}
		}
	}
}
