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
