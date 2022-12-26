package lexer

import (
	"fmt"
	"testing"
)

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
	compare(expect, result, t, "TestScan")
}

func TestEmpty(t *testing.T) {
	text := " \n\t"
	expect := []Token{mkEof(2, 2)}
	result, _ := Scan(text)
	compare(expect, result, t, "TestEmpty")
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
	compare(expect, result, t, "TestParens")
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
	compare(expect, result, t, "TestNewlines (Test 1)")

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
			Column: 10,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   2,
			Column: 12,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   3,
			Column: 10,
		},
		mkEof(3, 11),
	}
	result, _ = Scan(text)
	compare(expect, result, t, "TestNewlines (Test 2)")
}

func TestIdent(t *testing.T) {
	text := "x + wyvern * 3"
	expect := []Token{
		{
			Typeof: Ident,
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
			Typeof: Ident,
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
	compare(expect, result, t, "TestIdent (Test 1)")

	text = "x + wyvern/hamster"
	expect = []Token{
		{
			Typeof: Ident,
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
			Typeof: Ident,
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
			Typeof: Ident,
			Value:  "hamster",
			Line:   1,
			Column: 12,
		},
		mkEof(1, 19),
	}
	result, _ = Scan(text)
	compare(expect, result, t, "TestIdent (Test 2)")
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
	compare(expect, result, t, "TestNumbers (Test 1)")
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
	compare(expect, result, t, "TestSub")
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
	compare(expect, result, t, "TestPow")
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
