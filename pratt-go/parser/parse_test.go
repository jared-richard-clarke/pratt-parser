package parser

import (
	"fmt"
	"testing"
)

// Recursively tests for structural equality between nodes.
func equal(n, m Node) bool {
	switch n := n.(type) {
	case Empty:
		m, ok := m.(Empty)
		if !ok {
			return false
		}
		return n == m
	case Number:
		m, ok := m.(Number)
		if !ok {
			return false
		}
		return n == m
	case Symbol:
		m, ok := m.(Symbol)
		if !ok {
			return false
		}
		return n == m
	case Unary:
		m, ok := m.(Unary)
		if !ok {
			return false
		}
		return n.Op == m.Op && equal(n.X, m.X) && n.Line == m.Line && n.Column == m.Column
	case Binary:
		m, ok := m.(Binary)
		if !ok {
			return false
		}
		return n.Op == m.Op && equal(n.X, m.X) && equal(n.Y, m.Y) && n.Line == m.Line && n.Column == m.Column
	case ImpliedBinary:
		m, ok := m.(ImpliedBinary)
		if !ok {
			return false
		}
		return n.Op == m.Op && equal(n.X, m.X) && equal(n.Y, m.Y)
	case Call:
		m, ok := m.(Call)
		if !ok {
			return false
		}
		if !equal(n.Callee, m.Callee) {
			return false
		}
		// Since the default value for an empty "Args" is an empty slice,
		// checking for length of each slice will not cause a panic.
		if len(n.Args) != len(m.Args) {
			return false
		}
		for i := range n.Args {
			if !equal(n.Args[i], m.Args[i]) {
				return false
			}
		}
		return n.Line == m.Line && n.Column == m.Column
	default:
		return false
	}
}

func TestBasic(t *testing.T) {
	text := "1 + 2 * 3"
	expect := Binary{
		Op: "+",
		X: Number{
			Value:  1.0,
			Line:   1,
			Column: 1,
		},
		Y: Binary{
			Op: "*",
			X: Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Y: Number{
				Value:  3.0,
				Line:   1,
				Column: 9,
			},
			Line:   1,
			Column: 7,
		},
		Line:   1,
		Column: 3,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestBasic failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestBasic failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestUndefinedPrefixOp(t *testing.T) {
	text := "1 + * 7"
	result, err := Parse(text)
	if err == nil {
		msg := "TestUndefinedPrefixOp failed. Expected: error, Got: %s"
		t.Errorf(msg, result)
	}
}

func TestUnexpectedToken(t *testing.T) {
	text := "1.0.7 + 2"
	result, err := Parse(text)
	if err == nil {
		msg := "TestUnexpectedToken failed. Expected: error, Got: %s"
		t.Errorf(msg, result)
	}
}

func TestIncompleteExpression(t *testing.T) {
	text := "1 +"
	result, err := Parse(text)
	if err == nil {
		msg := "TestIncompleteExpression failed. Expected: error, Got: %s"
		t.Errorf(msg, result)
	}
}

func TestEmpty(t *testing.T) {
	text := ""
	expect := Empty{}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestEmpty (1) failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestEmpty (1) failed. Expected: %s, Got: %s", expect, result)
	}
	text = "\r\n   "
	if err != nil {
		t.Errorf("TestEmpty (2) failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestEmpty (2) failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestSymbol(t *testing.T) {
	text := "wyvern ^ 11"
	expect := Binary{
		Op: "^",
		X: Symbol{
			Value:  "wyvern",
			Line:   1,
			Column: 1,
		},
		Y: Number{
			Value:  11.0,
			Line:   1,
			Column: 10,
		},
		Line:   1,
		Column: 8,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestSymbol failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestSymbol failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestEqual(t *testing.T) {
	text := "7 + 4 = 11"
	expect := Binary{
		Op: "=",
		X: Binary{
			Op: "+",
			X: Number{
				Value:  7.0,
				Line:   1,
				Column: 1,
			},
			Y: Number{
				Value:  4.0,
				Line:   1,
				Column: 5,
			},
			Line:   1,
			Column: 3,
		},
		Y: Number{
			Value:  11.0,
			Line:   1,
			Column: 9,
		},
		Line:   1,
		Column: 7,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestEqual failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestEqual failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestLeftAssociative(t *testing.T) {
	text := "1 + 2 + 3"
	expect := Binary{
		Op: "+",
		X: Binary{
			Op: "+",
			X: Number{
				Value:  1.0,
				Line:   1,
				Column: 1,
			},
			Y: Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Line:   1,
			Column: 3,
		},
		Y: Number{
			Value:  3.0,
			Line:   1,
			Column: 9,
		},
		Line:   1,
		Column: 7,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestLeftAssociative failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestLeftAssociative failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestParens(t *testing.T) {
	text := "((1 + (2)))"
	expect := Binary{
		Op: "+",
		X: Number{
			Value:  1.0,
			Line:   1,
			Column: 3,
		},
		Y: Number{
			Value:  2.0,
			Line:   1,
			Column: 8,
		},
		Line:   1,
		Column: 5,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestParens failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestParens failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestMissingParen(t *testing.T) {
	text := "(1 + 2"
	result, err := Parse(text)
	if err == nil {
		msg := "TestMissingParen failed. Expected: error, Got: %s"
		t.Errorf(msg, result)
	}
}

func TestUnary(t *testing.T) {
	text := "--7"
	expect := Unary{
		Op: "-",
		X: Unary{
			Op: "-",
			X: Number{
				Value:  7.0,
				Line:   1,
				Column: 3,
			},
			Line:   1,
			Column: 2,
		},
		Line:   1,
		Column: 1,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestUnary failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestUnary failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestMinus(t *testing.T) {
	text := "7--7"
	expect := Binary{
		Op: "-",
		X: Number{
			Value:  7.0,
			Line:   1,
			Column: 1,
		},
		Y: Unary{
			Op: "-",
			X: Number{
				Value:  7.0,
				Line:   1,
				Column: 4,
			},
			Line:   1,
			Column: 3,
		},
		Line:   1,
		Column: 2,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestMinus failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestMinus failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestExponent(t *testing.T) {
	text := "1 ^ 2 ^ 3"
	expect := Binary{
		Op: "^",
		X: Number{
			Value:  1.0,
			Line:   1,
			Column: 1,
		},
		Y: Binary{
			Op: "^",
			X: Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Y: Number{
				Value:  3.0,
				Line:   1,
				Column: 9,
			},
			Line:   1,
			Column: 7,
		},
		Line:   1,
		Column: 3,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestExponent failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestExponent failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestCall(t *testing.T) {
	var text string
	var expect Node

	text = "square(5) + 2"
	expect = Binary{
		Op: "+",
		X: Call{
			Callee: Symbol{
				Value:  "square",
				Line:   1,
				Column: 1,
			},
			Args: []Node{
				Number{
					Value:  5.0,
					Line:   1,
					Column: 8,
				},
			},
			Line:   1,
			Column: 7,
		},
		Y: Number{
			Value:  2.0,
			Line:   1,
			Column: 13,
		},
		Line:   1,
		Column: 11,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestCall (1) failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestCall (1) failed. Expected: %s, Got: %s", expect, result)
	}

	text = "random()"
	expect = Call{
		Callee: Symbol{
			Value:  "random",
			Line:   1,
			Column: 1,
		},
		Args:   make([]Node, 0),
		Line:   1,
		Column: 7,
	}
	result, err = Parse(text)
	if err != nil {
		t.Errorf("TestCall (2) failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestCall (2) failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestUnclosedCall(t *testing.T) {
	text := "sin(7"
	result, err := Parse(text)
	if err == nil {
		msg := "TestUnclosedCall failed. Expected: error, Got: %s"
		t.Errorf(msg, result)
	}
}

func TestImpliedBinary(t *testing.T) {
	text := "7x"
	expect := ImpliedBinary{
		Op: "*",
		X: Number{
			Value:  7.0,
			Line:   1,
			Column: 1,
		},
		Y: Symbol{
			Value:  "x",
			Line:   1,
			Column: 2,
		},
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestImpliedBinary failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestImpliedBinary failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestAltOperators(t *testing.T) {
	text := "1 × 2 ÷ 3"
	expect := Binary{
		Op: "/",
		X: Binary{
			Op: "*",
			X: Number{
				Value:  1.0,
				Line:   1,
				Column: 1,
			},
			Y: Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Line:   1,
			Column: 3,
		},
		Y: Number{
			Value:  3.0,
			Line:   1,
			Column: 9,
		},
		Line:   1,
		Column: 7,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestAltOperators failed. Expected: %s, Got: %s", expect, err)
	}
	if !equal(expect, result) {
		t.Errorf("TestAltOperators failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestUnusedTokens(t *testing.T) {
	text := "1 + 2 3 +"
	result, err := Parse(text)
	if err == nil {
		msg := "TestUnusedTokens failed. Expected: error, Got: %s"
		t.Errorf(msg, result)
	}
}
