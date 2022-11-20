package parser

import "github/jared-richard-clarke/pratt/lexer"

type (
	Nud  func(lexer.Token) Node              // Null denotation
	Led  func(left Node, t lexer.Token) Node // Left denotation
	Nuds map[lexer.LexType]Nud               // lexeme -> Nud
	Leds map[lexer.LexType]Led               // lexeme -> Led
)
