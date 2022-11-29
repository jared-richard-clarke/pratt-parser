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
	result, _ := Parse(text)
	if !reflect.DeepEqual(expect, result) {
		t.Errorf("TestBase failed. Expected: %s, Got: %s", expect, result)
	}
}

// Recursively compares the concrete values of two structs that fulfill the Node interface.
// Currently using `reflect.DeepEqual(x, y any) bool` as I iron out the bugs in this particularly
// hairy piece of code.
//
// func equal(n1, n2 Node) bool {
// 	switch v1 := n1.(type) {
// 	case *Literal:
// 		v2, ok := n2.(*Literal)
// 		if !ok {
// 			return false
// 		}
// 		if v1.Typeof != v2.Typeof || v1.Value != v2.Value || v1.Line != v2.Line || v1.Column != v2.Column {
// 			return false
// 		}
// 		return true
// 	case *Unary:
// 		v2, ok := n2.(*Unary)
// 		if !ok {
// 			return false
// 		}
// 		if v1.Op != v2.Op || !equal(v1.X, v2.X) || v1.Line != v2.Line || v1.Column != v2.Column {
// 			return false
// 		}
// 		return true
// 	case *Binary:
// 		v2, ok := n2.(*Binary)
// 		if !ok {
// 			return false
// 		}
// 		if v1.Op != v2.Op || !equal(v1.X, v2.X) || !equal(v1.Y, v2.Y) || v1.Line != v2.Line || v1.Column != v2.Column {
// 			return false
// 		}
// 		return true
// 	default:
// 		return false
// 	}
// }
