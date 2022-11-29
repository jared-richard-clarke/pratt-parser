package parser

import (
	"github/jared-richard-clarke/pratt/internal/lexer"
	"reflect"
	"testing"
)

func TestBase(t *testing.T) {
	text := "1 + 2"
	expect := &Binary{
		Op: "+",
		X: &Literal{
			Typeof: lexer.Number,
			Value:  "1",
			Line:   1,
			Column: 1,
		},
		Y: &Literal{
			Typeof: lexer.Number,
			Value:  "2",
			Line:   1,
			Column: 5,
		},
		Line:   1,
		Column: 3,
	}
	result, err := Parse(text)
	if err != nil {
		t.Errorf("TestBase failed. Expected: %s, Got: %s", expect, err)
	} else if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestBase failed. Expected: %s, Got: %s", expect, result)
	}
}
