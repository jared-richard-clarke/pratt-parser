package pratt

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	result := Scan("2k \r	+ y + \nz")
	for _, v := range result {
		fmt.Println(v)
	}
}

func TestScan(t *testing.T) {

	text := "1 + 2 * 3"

	expect := []token{
		{
			typeof: lexNumber,
			value:  1.0,
			line:   1,
			column: 1,
		},
		{
			typeof: lexAdd,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: lexNumber,
			value:  2.0,
			line:   1,
			column: 5,
		},
		{
			typeof: lexMul,
			value:  '*',
			line:   1,
			column: 7,
		},
		{
			typeof: lexNumber,
			value:  3.0,
			line:   1,
			column: 9,
		},
		{
			typeof: lexEOF,
			line:   1,
			column: 10,
		},
	}

	result := Scan(text)
	compare(expect, result, t, "TestScan")
}

func TestParens(t *testing.T) {

	text := "(1 + 2) * 3"

	expect := []token{
		{
			typeof: lexOpenParen,
			value:  '(',
			line:   1,
			column: 1,
		},
		{
			typeof: lexNumber,
			value:  1.0,
			line:   1,
			column: 2,
		},
		{
			typeof: lexAdd,
			value:  '+',
			line:   1,
			column: 4,
		},
		{
			typeof: lexNumber,
			value:  2.0,
			line:   1,
			column: 6,
		},
		{
			typeof: lexCloseParen,
			value:  ')',
			line:   1,
			column: 7,
		},
		{
			typeof: lexMul,
			value:  '*',
			line:   1,
			column: 9,
		},
		{
			typeof: lexNumber,
			value:  3.0,
			line:   1,
			column: 11,
		},
		{
			typeof: lexEOF,
			line:   1,
			column: 12,
		},
	}

	result := Scan(text)
	compare(expect, result, t, "TestParens")
}

func TestNewlines(t *testing.T) {
	text := "1 + 2\n * 3"

	expect := []token{
		{
			typeof: lexNumber,
			value:  1.0,
			line:   1,
			column: 1,
		},
		{
			typeof: lexAdd,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: lexNumber,
			value:  2.0,
			line:   1,
			column: 5,
		},
		{
			typeof: lexMul,
			value:  '*',
			line:   2,
			column: 2,
		},
		{
			typeof: lexNumber,
			value:  3.0,
			line:   2,
			column: 4,
		},
		{
			typeof: lexEOF,
			line:   2,
			column: 5,
		},
	}

	result := Scan(text)
	compare(expect, result, t, "TestNewlines (Part 1)")

	text = `1 +
	        2 *
	        3`

	expect = []token{
		{
			typeof: lexNumber,
			value:  1.0,
			line:   1,
			column: 1,
		},
		{
			typeof: lexAdd,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: lexNumber,
			value:  2.0,
			line:   2,
			column: 10,
		},
		{
			typeof: lexMul,
			value:  '*',
			line:   2,
			column: 12,
		},
		{
			typeof: lexNumber,
			value:  3.0,
			line:   3,
			column: 10,
		},
		{
			typeof: lexEOF,
			line:   3,
			column: 11,
		},
	}

	result = Scan(text)
	compare(expect, result, t, "TestNewlines (Part 2)")
}

func TestIdent(t *testing.T) {

	text := "x + wyvern * 3"

	expect := []token{
		{
			typeof: lexIdent,
			value:  "x",
			line:   1,
			column: 1,
		},
		{
			typeof: lexAdd,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: lexIdent,
			value:  "wyvern",
			line:   1,
			column: 5,
		},
		{
			typeof: lexMul,
			value:  '*',
			line:   1,
			column: 12,
		},
		{
			typeof: lexNumber,
			value:  3.0,
			line:   1,
			column: 14,
		},
		{
			typeof: lexEOF,
			line:   1,
			column: 15,
		},
	}

	result := Scan(text)
	compare(expect, result, t, "TestIdent (Part 1)")

	text = "x + wyvern/hamster"

	expect = []token{
		{
			typeof: lexIdent,
			value:  "x",
			line:   1,
			column: 1,
		},
		{
			typeof: lexAdd,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: lexIdent,
			value:  "wyvern",
			line:   1,
			column: 5,
		},
		{
			typeof: lexDiv,
			value:  '/',
			line:   1,
			column: 11,
		},
		{
			typeof: lexIdent,
			value:  "hamster",
			line:   1,
			column: 12,
		},
		{
			typeof: lexEOF,
			line:   1,
			column: 19,
		},
	}

	result = Scan(text)
	compare(expect, result, t, "TestIdent (Part 2)")
}

// utility functions

func compare(expect []token, result []*token, t *testing.T, name string) {
	for i, c := range result {
		e := expect[i]
		g := *c
		if e != g {
			t.Errorf("Test %s failed. Expected: %v, Got: %v", name, e, g)
		}
	}
}
