package parser

import (
	"fmt"
	"github/jared-richard-clarke/pratt/lexer"
)

// === Under Heavy Construction ===

type nud func(lexer.Token) (Node, error)       // Null denotation
type led func(Node, lexer.Token) (Node, error) // Left denotation

type table struct {
	nuds     map[lexer.LexType]nud // lexeme -> Nud
	leds     map[lexer.LexType]led // lexeme -> Led
	prebinds map[lexer.LexType]int // lexeme -> prefix binding power
	binds    map[lexer.LexType]int // lexeme -> binding power
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
	for rbp < p.binds[token.Typeof] {
		token := p.next()
		led := p.leds[token.Typeof]
		left, err = led(left, token)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

// Always returns error. Has Node type to satisfy "nud".
func (p *parser) error(t lexer.Token) (Node, error) {
	err := fmt.Errorf("unexpected lexeme %s :%d:%d", t.Value, t.Line, t.Column)
	return nil, err
}

// Always returns Node. Has error type to satisfy "nud".
func (p *parser) literal(t lexer.Token) (Node, error) {
	return &Literal{
		Typeof: t.Typeof,
		Value:  t.Value,
		Line:   t.Line,
		Column: t.Column,
	}, nil
}

func (p *parser) unary(t lexer.Token) (Node, error) {
	x, err := p.expression(p.prebinds[t.Typeof])
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
	right, err := p.expression(p.binds[token.Typeof])
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
	right, err := p.expression(p.binds[token.Typeof] - 1)
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
			nuds:     make(map[lexer.LexType]nud),
			leds:     make(map[lexer.LexType]led),
			prebinds: make(map[lexer.LexType]int),
			binds:    make(map[lexer.LexType]int),
		},
	}
	// Helper functions build lookup tables.
	register := func(t lexer.LexType, n nud) {
		pratt.nuds[t] = n
	}
	prefix := func(bp int, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.nuds[t] = pratt.unary
			pratt.prebinds[t] = bp
		}
	}
	infix := func(bp int, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.leds[t] = pratt.binary
			pratt.binds[t] = bp
		}
	}
	infixr := func(bp int, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.leds[t] = pratt.binaryr
			pratt.binds[t] = bp
		}
	}
	// Initialize lookup tables.
	register(lexer.Error, pratt.error)
	register(lexer.Number, pratt.literal)
	register(lexer.Ident, pratt.literal)
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
		table:  pratt.table, // Reuse lookup table from package initialization.
	}
	node, err := pratt.expression(0)
	if err != nil {
		return nil, err
	}
	return node, nil
}
