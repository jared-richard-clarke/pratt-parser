package pratt

import (
	"testing"
)

func TestScan(t *testing.T) {

	text := "1 + 2 * 3"

	expect := []token{
		{
			typeof: number,
			value:  1.0,
			line:   1,
			column: 1,
		},
		{
			typeof: add,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: number,
			value:  2.0,
			line:   1,
			column: 5,
		},
		{
			typeof: mul,
			value:  '*',
			line:   1,
			column: 7,
		},
		{
			typeof: number,
			value:  3.0,
			line:   1,
			column: 9,
		},
		{
			typeof: eof,
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
			typeof: number,
			value:  1.0,
			line:   1,
			column: 1,
		},
		{
			typeof: add,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: number,
			value:  2.0,
			line:   1,
			column: 5,
		},
		{
			typeof: mul,
			value:  '*',
			line:   2,
			column: 2,
		},
		{
			typeof: number,
			value:  3.0,
			line:   2,
			column: 4,
		},
		{
			typeof: eof,
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
			typeof: number,
			value:  1.0,
			line:   1,
			column: 1,
		},
		{
			typeof: add,
			value:  '+',
			line:   1,
			column: 3,
		},
		{
			typeof: number,
			value:  2.0,
			line:   2,
			column: 10,
		},
		{
			typeof: mul,
			value:  '*',
			line:   2,
			column: 12,
		},
		{
			typeof: number,
			value:  3.0,
			line:   3,
			column: 10,
		},
		{
			typeof: eof,
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
			typeof: oParen,
			value:  '(',
			line:   1,
			column: 1,
		},
		{
			typeof: number,
			value:  1.0,
			line:   1,
			column: 2,
		},
		{
			typeof: add,
			value:  '+',
			line:   1,
			column: 4,
		},
		{
			typeof: number,
			value:  2.0,
			line:   1,
			column: 6,
		},
		{
			typeof: cParen,
			value:  ')',
			line:   1,
			column: 7,
		},
		{
			typeof: mul,
			value:  '*',
			line:   1,
			column: 9,
		},
		{
			typeof: number,
			value:  3.0,
			line:   1,
			column: 11,
		},
		{
			typeof: eof,
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
