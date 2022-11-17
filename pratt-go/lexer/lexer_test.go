package lexer

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	result := Scan("# 7..0/5wyvern\n 100")
	for _, v := range result {
		fmt.Println(v)
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
			Length: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 5,
			Length: 1,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   1,
			Column: 7,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   1,
			Column: 9,
			Length: 1,
		},
		mkEof(1, 10),
	}
	result := Scan(text)
	compare(expect, result, t, "TestScan")
}

func TestEmpty(t *testing.T) {
	text := " \n\t"
	expect := []Token{mkEof(2, 2)}
	result := Scan(text)
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
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "1",
			Line:   1,
			Column: 2,
			Length: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 4,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 6,
			Length: 1,
		},
		{
			Typeof: CloseParen,
			Value:  ")",
			Line:   1,
			Column: 7,
			Length: 1,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   1,
			Column: 9,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   1,
			Column: 11,
			Length: 1,
		},
		mkEof(1, 12),
	}
	result := Scan(text)
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
			Length: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 5,
			Length: 1,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   2,
			Column: 2,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   2,
			Column: 4,
			Length: 1,
		},
		mkEof(2, 5),
	}
	result := Scan(text)
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
			Length: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   2,
			Column: 10,
			Length: 1,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   2,
			Column: 12,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   3,
			Column: 10,
			Length: 1,
		},
		mkEof(3, 11),
	}
	result = Scan(text)
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
			Length: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
			Length: 1,
		},
		{
			Typeof: Ident,
			Value:  "wyvern",
			Line:   1,
			Column: 5,
			Length: 6,
		},
		{
			Typeof: Mul,
			Value:  "*",
			Line:   1,
			Column: 12,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "3",
			Line:   1,
			Column: 14,
			Length: 1,
		},
		mkEof(1, 15),
	}
	result := Scan(text)
	compare(expect, result, t, "TestIdent (Test 1)")

	text = "x + wyvern/hamster"
	expect = []Token{
		{
			Typeof: Ident,
			Value:  "x",
			Line:   1,
			Column: 1,
			Length: 1,
		},
		{
			Typeof: Add,
			Value:  "+",
			Line:   1,
			Column: 3,
			Length: 1,
		},
		{
			Typeof: Ident,
			Value:  "wyvern",
			Line:   1,
			Column: 5,
			Length: 6,
		},
		{
			Typeof: Div,
			Value:  "/",
			Line:   1,
			Column: 11,
			Length: 1,
		},
		{
			Typeof: Ident,
			Value:  "hamster",
			Line:   1,
			Column: 12,
			Length: 7,
		},
		mkEof(1, 19),
	}
	result = Scan(text)
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
			Length: 3,
		},
		{
			Typeof: Div,
			Value:  "/",
			Line:   1,
			Column: 4,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 5,
			Length: 1,
		},
		mkEof(1, 6),
	}
	result := Scan(text)
	compare(expect, result, t, "TestNumbers (Test 1)")

	text = "7.5.0"
	expect = []Token{
		{
			Typeof: Number,
			Value:  "7.5",
			Line:   1,
			Column: 1,
			Length: 3,
		},
		{
			Typeof: Error,
			Value:  ".",
			Line:   1,
			Column: 4,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "0",
			Line:   1,
			Column: 5,
			Length: 1,
		},
		mkEof(1, 6),
	}
	result = Scan(text)
	compare(expect, result, t, "TestNumbers (Test 2)")
}

func TestSub(t *testing.T) {
	text := "1 - -2"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "1",
			Line:   1,
			Column: 1,
			Length: 1,
		},
		{
			Typeof: Sub,
			Value:  "-",
			Line:   1,
			Column: 3,
			Length: 1,
		},
		{
			Typeof: Sub,
			Value:  "-",
			Line:   1,
			Column: 5,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 6,
			Length: 1,
		},
		mkEof(1, 7),
	}
	result := Scan(text)
	compare(expect, result, t, "TestSub")
}

func TestExp(t *testing.T) {
	text := "4^2"
	expect := []Token{
		{
			Typeof: Number,
			Value:  "4",
			Line:   1,
			Column: 1,
			Length: 1,
		},
		{
			Typeof: Exp,
			Value:  "^",
			Line:   1,
			Column: 2,
			Length: 1,
		},
		{
			Typeof: Number,
			Value:  "2",
			Line:   1,
			Column: 3,
			Length: 1,
		},
		mkEof(1, 4),
	}
	result := Scan(text)
	compare(expect, result, t, "TestExp")
}

// utility functions

func mkEof(l, c int) Token {
	return Token{
		Typeof: EOF,
		Line:   l,
		Column: c,
		Length: 1,
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
