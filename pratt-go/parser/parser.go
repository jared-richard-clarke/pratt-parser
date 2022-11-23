package parser

import (
	"github/jared-richard-clarke/pratt/lexer"
	"strconv"
)

// === Under Heavy Construction ===

type nud func(lexer.Token) (Node, error)       // Null denotation
type led func(Node, lexer.Token) (Node, error) // Left denotation

type table struct {
	nuds map[lexer.LexType]nud // lexeme -> Nud
	leds map[lexer.LexType]led // lexeme -> Led
	rbp  map[lexer.LexType]int // lexeme -> right binding power
	lbp  map[lexer.LexType]int // lexeme -> left binding power
}

type parser struct {
	src    []lexer.Token // token source
	length int           // len(src)
	index  int           // src[index]
	end    int           // src[len(src) - 1]
	*table               // parser and binding lookup
}

func (p *parser) next() lexer.Token {
	// From final index onwards, returns final token — usually EOF.
	if p.index >= p.length {
		return p.src[p.end]
	}
	t := p.src[p.index]
	p.index += 1
	return t
}

func (p *parser) expression(rbp int) (Node, error) {
	token := p.next()
	nud := p.nuds[token.Typeof]
	left, err := nud(token)
	if err != nil {
		return nil, err
	}
	for rbp < p.lbp[token.Typeof] {
		token := p.next()
		led := p.leds[token.Typeof]
		left, err = led(left, token)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (p *parser) number(t lexer.Token) (Node, error) {
	num, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return nil, err
	}
	return &Number{
		Value:  num,
		Line:   t.Line,
		Column: t.Column,
	}, nil
}

func (p *parser) ident(t lexer.Token) (Node, error) {
	return &Ident{
		Value:  t.Value,
		Line:   t.Line,
		Column: t.Column,
	}, nil
}

func (p *parser) unary(t lexer.Token) (Node, error) {
	x, err := p.expression(p.rbp[t.Typeof])
	if err != nil {
		return nil, err
	}
	return &Unary{
		Op:     t.Value,
		X:      x,
		Line:   t.Line,
		Column: t.Column,
	}, nil
}

func (p *parser) binary(left Node, token lexer.Token) (Node, error) {
	right, err := p.expression(p.lbp[token.Typeof])
	if err != nil {
		return nil, err
	}
	return &Binary{
		Op:     token.Value,
		X:      left,
		Y:      right,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

func (p *parser) binaryr(left Node, token lexer.Token) (Node, error) {
	right, err := p.expression(p.lbp[token.Typeof] - 1)
	if err != nil {
		return nil, err
	}
	return &Binary{
		Op:     token.Value,
		X:      left,
		Y:      right,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

// Carries parser's internal state. Should persist throughout package lifetime.
var pratt parser

func init() {
	// Build lookup table at package initialization.
	pratt = parser{
		table: &table{
			nuds: make(map[lexer.LexType]nud),
			leds: make(map[lexer.LexType]led),
			rbp:  make(map[lexer.LexType]int),
			lbp:  make(map[lexer.LexType]int),
		},
	}
	// Helper functions build lookup tables.
	addNud := func(bp int, t lexer.LexType, n nud) {
		pratt.nuds[t] = n
		pratt.rbp[t] = bp
	}
	addLed := func(t lexer.LexType, bp int, l led) {
		pratt.leds[t] = l
		pratt.lbp[t] = bp
	}
	symbol := func(t lexer.LexType, n nud) {
		addNud(0, t, n)
	}
	prefix := func(bp int, ts ...lexer.LexType) {
		for _, t := range ts {
			addNud(bp, t, pratt.unary)
		}
	}
	infix := func(bp int, ts ...lexer.LexType) {
		for _, t := range ts {
			addLed(t, bp, pratt.binary)
		}
	}
	infixr := func(bp int, ts ...lexer.LexType) {
		for _, t := range ts {
			addLed(t, bp, pratt.binaryr)
		}
	}
	// Initialize lookup tables.
	symbol(lexer.Number, pratt.number)
	symbol(lexer.Ident, pratt.ident)
	infix(50, lexer.Add, lexer.Sub)
	infix(60, lexer.Mul, lexer.Div)
	infixr(70, lexer.Pow)
	prefix(80, lexer.Add, lexer.Sub)
}

// Parser API: inputs tokens, outputs either AST or Error
func Parse(ts []lexer.Token) (Node, error) {
	// Set parser state.
	pratt = parser{
		src:    ts,
		length: len(ts),
		index:  0,
		end:    len(ts) - 1,
		table:  pratt.table, // Persist lookup table from package initialization.
	}
	return nil, nil
}
