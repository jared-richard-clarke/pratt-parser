package parser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	text := "1 + 2 * 3"
	expect := &Binary{
		Op: "+",
		X: &Number{
			Value:  1.0,
			Line:   1,
			Column: 1,
		},
		Y: &Binary{
			Op: "*",
			X: &Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Y: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestBasic failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestEmpty(t *testing.T) {
	text := ""
	expect := &Empty{}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestEmpty (1) failed. Expected: %s, Got: %s", expect, err)
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestEmpty (1) failed. Expected: %s, Got: %s", expect, result)
	}
	text = "\r\n   "
	if err != nil {
		t.Errorf("TestEmpty (2) failed. Expected: %s, Got: %s", expect, err)
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestEmpty (2) failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestSymbol(t *testing.T) {
	text := "wyvern ^ 11"
	expect := &Binary{
		Op: "^",
		X: &Symbol{
			Value:  "wyvern",
			Line:   1,
			Column: 1,
		},
		Y: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestSymbol failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestLeftAssociative(t *testing.T) {
	text := "1 + 2 + 3"
	expect := &Binary{
		Op: "+",
		X: &Binary{
			Op: "+",
			X: &Number{
				Value:  1.0,
				Line:   1,
				Column: 1,
			},
			Y: &Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Line:   1,
			Column: 3,
		},
		Y: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestLeftAssociative failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestParens(t *testing.T) {
	text := "((1 + (2)))"
	expect := &Binary{
		Op: "+",
		X: &Number{
			Value:  1.0,
			Line:   1,
			Column: 3,
		},
		Y: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestParens failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestUnary(t *testing.T) {
	text := "--7"
	expect := &Unary{
		Op: "-",
		X: &Unary{
			Op: "-",
			X: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestUnary failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestMinus(t *testing.T) {
	text := "7--7"
	expect := &Binary{
		Op: "-",
		X: &Number{
			Value:  7.0,
			Line:   1,
			Column: 1,
		},
		Y: &Unary{
			Op: "-",
			X: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestMinus failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestExponent(t *testing.T) {
	text := "1 ^ 2 ^ 3"
	expect := &Binary{
		Op: "^",
		X: &Number{
			Value:  1.0,
			Line:   1,
			Column: 1,
		},
		Y: &Binary{
			Op: "^",
			X: &Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Y: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestExponent failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestCall(t *testing.T) {
	var text string
	var expect Node

	text = "square(5) + 2"
	expect = &Binary{
		Op: "+",
		X: &Call{
			Callee: &Symbol{
				Value:  "square",
				Line:   1,
				Column: 1,
			},
			Args: []Node{
				&Number{
					Value:  5.0,
					Line:   1,
					Column: 8,
				},
			},
			Line:   1,
			Column: 7,
		},
		Y: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestCall (1) failed. Expected: %s, Got: %s", expect, result)
	}

	text = "random()"
	expect = &Call{
		Callee: &Symbol{
			Value:  "random",
			Line:   1,
			Column: 1,
		},
		Args:   nil,
		Line:   1,
		Column: 7,
	}
	result, err = Parse(text)
	if err != nil {
		t.Errorf("TestCall (2) failed. Expected: %s, Got: %s", expect, err)
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestCall (2) failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestImpliedBinary(t *testing.T) {
	text := "7x"
	expect := &ImpliedBinary{
		Op: "*",
		X: &Number{
			Value:  7.0,
			Line:   1,
			Column: 1,
		},
		Y: &Symbol{
			Value:  "x",
			Line:   1,
			Column: 2,
		},
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestImpliedBinary failed. Expected: %s, Got: %s", expect, err)
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestImpliedBinary failed. Expected: %s, Got: %s", expect, result)
	}
}

func TestAltOperators(t *testing.T) {
	text := "1 × 2 ÷ 3"
	expect := &Binary{
		Op: "/",
		X: &Binary{
			Op: "*",
			X: &Number{
				Value:  1.0,
				Line:   1,
				Column: 1,
			},
			Y: &Number{
				Value:  2.0,
				Line:   1,
				Column: 5,
			},
			Line:   1,
			Column: 3,
		},
		Y: &Number{
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
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestAltOperators failed. Expected: %s, Got: %s", expect, result)
	}
}
