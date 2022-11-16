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
	expect := []token{
		{
			typeof: lexNumber,
			value:  "1",
			line:   1,
			column: 1,
			length: 1,
		},
		{
			typeof: lexAdd,
			value:  "+",
			line:   1,
			column: 3,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "2",
			line:   1,
			column: 5,
			length: 1,
		},
		{
			typeof: lexMul,
			value:  "*",
			line:   1,
			column: 7,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "3",
			line:   1,
			column: 9,
			length: 1,
		},
		eof(1, 10),
	}
	result := Scan(text)
	compare(expect, result, t, "TestScan")
}

func TestEmpty(t *testing.T) {
	text := " \n\t"
	expect := []token{eof(2, 2)}
	result := Scan(text)
	compare(expect, result, t, "TestEmpty")
}

func TestParens(t *testing.T) {
	text := "(1 + 2) * 3"
	expect := []token{
		{
			typeof: lexOpenParen,
			value:  "(",
			line:   1,
			column: 1,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "1",
			line:   1,
			column: 2,
			length: 1,
		},
		{
			typeof: lexAdd,
			value:  "+",
			line:   1,
			column: 4,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "2",
			line:   1,
			column: 6,
			length: 1,
		},
		{
			typeof: lexCloseParen,
			value:  ")",
			line:   1,
			column: 7,
			length: 1,
		},
		{
			typeof: lexMul,
			value:  "*",
			line:   1,
			column: 9,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "3",
			line:   1,
			column: 11,
			length: 1,
		},
		eof(1, 12),
	}
	result := Scan(text)
	compare(expect, result, t, "TestParens")
}

func TestNewlines(t *testing.T) {
	text := "1 + 2\n * 3"
	expect := []token{
		{
			typeof: lexNumber,
			value:  "1",
			line:   1,
			column: 1,
			length: 1,
		},
		{
			typeof: lexAdd,
			value:  "+",
			line:   1,
			column: 3,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "2",
			line:   1,
			column: 5,
			length: 1,
		},
		{
			typeof: lexMul,
			value:  "*",
			line:   2,
			column: 2,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "3",
			line:   2,
			column: 4,
			length: 1,
		},
		eof(2, 5),
	}
	result := Scan(text)
	compare(expect, result, t, "TestNewlines (Test 1)")

	text = `1 +
	        2 *
	        3`
	expect = []token{
		{
			typeof: lexNumber,
			value:  "1",
			line:   1,
			column: 1,
			length: 1,
		},
		{
			typeof: lexAdd,
			value:  "+",
			line:   1,
			column: 3,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "2",
			line:   2,
			column: 10,
			length: 1,
		},
		{
			typeof: lexMul,
			value:  "*",
			line:   2,
			column: 12,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "3",
			line:   3,
			column: 10,
			length: 1,
		},
		eof(3, 11),
	}
	result = Scan(text)
	compare(expect, result, t, "TestNewlines (Test 2)")
}

func TestIdent(t *testing.T) {
	text := "x + wyvern * 3"
	expect := []token{
		{
			typeof: lexIdent,
			value:  "x",
			line:   1,
			column: 1,
			length: 1,
		},
		{
			typeof: lexAdd,
			value:  "+",
			line:   1,
			column: 3,
			length: 1,
		},
		{
			typeof: lexIdent,
			value:  "wyvern",
			line:   1,
			column: 5,
			length: 6,
		},
		{
			typeof: lexMul,
			value:  "*",
			line:   1,
			column: 12,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "3",
			line:   1,
			column: 14,
			length: 1,
		},
		eof(1, 15),
	}
	result := Scan(text)
	compare(expect, result, t, "TestIdent (Test 1)")

	text = "x + wyvern/hamster"
	expect = []token{
		{
			typeof: lexIdent,
			value:  "x",
			line:   1,
			column: 1,
			length: 1,
		},
		{
			typeof: lexAdd,
			value:  "+",
			line:   1,
			column: 3,
			length: 1,
		},
		{
			typeof: lexIdent,
			value:  "wyvern",
			line:   1,
			column: 5,
			length: 6,
		},
		{
			typeof: lexDiv,
			value:  "/",
			line:   1,
			column: 11,
			length: 1,
		},
		{
			typeof: lexIdent,
			value:  "hamster",
			line:   1,
			column: 12,
			length: 7,
		},
		eof(1, 19),
	}
	result = Scan(text)
	compare(expect, result, t, "TestIdent (Test 2)")
}

func TestNumbers(t *testing.T) {
	text := "7.5/2"
	expect := []token{
		{
			typeof: lexNumber,
			value:  "7.5",
			line:   1,
			column: 1,
			length: 3,
		},
		{
			typeof: lexDiv,
			value:  "/",
			line:   1,
			column: 4,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "2",
			line:   1,
			column: 5,
			length: 1,
		},
		eof(1, 6),
	}
	result := Scan(text)
	compare(expect, result, t, "TestNumbers (Test 1)")

	text = "7.5.0"
	expect = []token{
		{
			typeof: lexNumber,
			value:  "7.5",
			line:   1,
			column: 1,
			length: 3,
		},
		{
			typeof: lexError,
			value:  ".",
			line:   1,
			column: 4,
			length: 1,
		},
		{
			typeof: lexNumber,
			value:  "0",
			line:   1,
			column: 5,
			length: 1,
		},
		eof(1, 6),
	}
	result = Scan(text)
	compare(expect, result, t, "TestNumbers (Test 2)")
}

// utility functions

func eof(l, c int) token {
	return token{
		typeof: lexEOF,
		line:   l,
		column: c,
		length: 1,
	}
}

func compare(e []token, r []token, t *testing.T, n string) {
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
