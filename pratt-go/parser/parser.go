package parser

import "github/jared-richard-clarke/pratt/lexer"

type (
	Nud     func(lexer.Token) Node              // Null denotation
	Led     func(left Node, t lexer.Token) Node // Left denotation
	NudMap  map[lexer.LexType]Nud               // token -> Nud
	LedMap  map[lexer.LexType]Led               // token -> Led
	NudBind map[lexer.LexType]int               // prefix binding powers
	LedBind map[lexer.LexType]int               // infix binding powers
)
