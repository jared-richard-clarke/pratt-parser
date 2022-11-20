package parser

import "github/jared-richard-clarke/pratt/lexer"

type (
	Nud     func(lexer.Token) Node       // Null denotation
	Led     func(Node, lexer.Token) Node // Left denotation
	NudMap  map[lexer.LexType]Nud        // lexeme -> Nud
	LedMap  map[lexer.LexType]Led        // lexeme -> Led
	NudBind map[lexer.LexType]int        // prefix binding powers
	LedBind map[lexer.LexType]int        // infix binding powers
)
