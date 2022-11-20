package parser

import "github/jared-richard-clarke/pratt/lexer"

type Nud func(lexer.Token) Node

type Led func(left Node, t lexer.Token) Node
