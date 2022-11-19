package parser

import "github/jared-richard-clarke/pratt/lexer"

type parser struct {
	src   []lexer.Token // Token source.
	token lexer.Token   // Current token.
	index int           // Current index in src.
	ast   Expr          // Abstract Syntax Tree.
}
