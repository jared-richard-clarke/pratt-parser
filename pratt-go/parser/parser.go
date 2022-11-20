package parser

import "github/jared-richard-clarke/pratt/lexer"

type nud func(lexer.Token) Node       // Null denotation
type led func(Node, lexer.Token) Node // Left denotation

type parseTable struct {
	nudParse map[lexer.LexType]nud // lexeme -> Nud
	ledParse map[lexer.LexType]led // lexeme -> Led

	nudBind map[lexer.LexType]int // lexeme -> prefix precedence
	ledBind map[lexer.LexType]int // lexeme -> infix precedence
}

type parser struct {
	src   []lexer.Token
	index int
}
