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
			Value:  "1",
			Line:   1,
			Column: 1,
		},
		Y: &Binary{
			Op: "*",
			X: &Number{
				Value:  "2",
				Line:   1,
				Column: 5,
			},
			Y: &Number{
				Value:  "3",
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

func TestParens(t *testing.T) {
	text := "((1 + (2)))"
	expect := &Binary{
		Op: "+",
		X: &Number{
			Value:  "1",
			Line:   1,
			Column: 3,
		},
		Y: &Number{
			Value:  "2",
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
				Value:  "7",
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
			Value:  "7",
			Line:   1,
			Column: 1,
		},
		Y: &Unary{
			Op: "-",
			X: &Number{
				Value:  "7",
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
			Value:  "1",
			Line:   1,
			Column: 1,
		},
		Y: &Binary{
			Op: "^",
			X: &Number{
				Value:  "2",
				Line:   1,
				Column: 5,
			},
			Y: &Number{
				Value:  "3",
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
