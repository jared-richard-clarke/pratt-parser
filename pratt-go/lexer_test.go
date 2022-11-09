package pratt

import (
	"testing"
)

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
			value:  nil,
			line:   1,
			column: 10,
		},
	}

	result, _ := Scan(text)
	for i, c := range result {
		e := expect[i]
		g := *c
		if e != g {
			t.Errorf("Test Scan (Part 1) failed. Expected: %v, Got: %v", e, g)
		}
	}

	text = "1 + 2\n * 3"

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
			value:  nil,
			line:   2,
			column: 5,
		},
	}

	result, _ = Scan(text)
	for i, c := range result {
		e := expect[i]
		g := *c
		if e != g {
			t.Errorf("Test Scan (Part 2) failed. Expected: %v, Got: %v", e, g)
		}
	}

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
			value:  nil,
			line:   3,
			column: 11,
		},
	}

	result, _ = Scan(text)
	for i, c := range result {
		e := expect[i]
		g := *c
		if e != g {
			t.Errorf("Test Scan (Part 3) failed. Expected: %v, Got: %v", e, g)
		}
	}

	text = "(1 + 2) * 3"

	expect = []token{
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
			value:  nil,
			line:   1,
			column: 12,
		},
	}

	result, _ = Scan(text)
	for i, c := range result {
		e := expect[i]
		g := *c
		if e != g {
			t.Errorf("Test Scan (Part 4) failed. Expected: %v, Got: %v", e, g)
		}
	}
}
