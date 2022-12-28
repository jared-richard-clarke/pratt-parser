package parser

import (
	"fmt"
	"github/jared-richard-clarke/pratt/internal/lexer"
	"strconv"
)

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
	index  int           // src[index]
	end    int           // src[len(src) - 1]
	*table               // parser and binding lookup
}

func (p *parser) next() lexer.Token {
	// From final index onwards, returns final token — usually EOF.
	if p.index >= p.end {
		return p.src[p.end]
	}
	t := p.src[p.index]
	p.index += 1
	return t
}

func (p *parser) peek() lexer.LexType {
	return p.src[p.index].Typeof
}

func (p *parser) match(expect lexer.LexType) bool {
	return p.peek() == expect
}

func (p *parser) expression(rbp int) (Node, error) {
	token := p.next()
	nud, ok := p.nuds[token.Typeof]
	if !ok {
		msg := "undefined prefix operation for %q line:%d column:%d"
		return nil, fmt.Errorf(msg, token.Value, token.Line, token.Column)
	}
	left, err := nud(token)
	if err != nil {
		return nil, err
	}
	for rbp < p.binds[p.peek()] {
		token := p.next()
		led, ok := p.leds[token.Typeof]
		if !ok {
			msg := "undefined infix operation for %q line:%d column:%d"
			return nil, fmt.Errorf(msg, token.Value, token.Line, token.Column)
		}
		left, err = led(left, token)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (p *parser) eof(token lexer.Token) (Node, error) {
	if len(p.src) == 1 {
		return &Empty{}, nil
	}
	msg := "incomplete expression, unexpected <EOF> line:%d column:%d"
	err := fmt.Errorf(msg, token.Line, token.Column)
	return nil, err
}

func (p *parser) number(token lexer.Token) (Node, error) {
	num, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		msg := "invalid number: %s line:%d column:%d"
		return nil, fmt.Errorf(msg, token.Value, token.Line, token.Column)
	}
	return &Number{
		Value:  num,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

// Always returns Node. Has error type to satisfy "nud".
func (p *parser) symbol(token lexer.Token) (Node, error) {
	return &Symbol{
		Value:  token.Value,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

func (p *parser) unary(token lexer.Token) (Node, error) {
	node, err := p.expression(p.prebinds[token.Typeof])
	if err != nil {
		return nil, err
	}
	return &Unary{
		Op:     token.Value,
		X:      node,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

func (p *parser) binary(left Node, token lexer.Token) (Node, error) {
	right, err := p.expression(p.binds[token.Typeof])
	if err != nil {
		return nil, err
	}
	if token.Typeof == lexer.ImpMul {
		return &ImpBinary{
			Op: token.Value,
			X:  left,
			Y:  right,
		}, nil
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

func (p *parser) paren(token lexer.Token) (Node, error) {
	position := fmt.Sprintf("line:%d column:%d", token.Line, token.Column)
	node, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	if !p.match(lexer.CloseParen) {
		msg := "for '(' at %s, missing matching ')'"
		return nil, fmt.Errorf(msg, position)
	}
	p.next()
	return node, nil
}

func (p *parser) call(left Node, token lexer.Token) (Node, error) {
	position := fmt.Sprintf("line:%d column:%d", token.Line, token.Column)
	if p.match(lexer.CloseParen) {
		p.next()
		return &Call{Fun: left}, nil
	}
	var args []Node
	for {
		node, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		args = append(args, node)
		if !p.match(lexer.Comma) {
			break
		}
		p.next()
	}
	if !p.match(lexer.CloseParen) {
		msg := "for function call at %s, missing matching ')'"
		return nil, fmt.Errorf(msg, position)
	}
	p.next()
	return &Call{
		Fun:    left,
		Args:   args,
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
	set := func(t lexer.LexType, n nud) {
		pratt.nuds[t] = n
		pratt.binds[t] = 0
	}
	prefix := func(bp int, n nud, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.nuds[t] = n
			pratt.prebinds[t] = bp
		}
	}
	infix := func(bp int, l led, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.leds[t] = l
			pratt.binds[t] = bp
		}
	}
	infixr := func(bp int, l led, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.leds[t] = l
			pratt.binds[t] = bp
		}
	}

	// Initialize lookup tables.
	set(lexer.EOF, pratt.eof)
	set(lexer.Number, pratt.number)
	set(lexer.Symbol, pratt.symbol)
	set(lexer.OpenParen, pratt.paren)
	infix(10, pratt.binary, lexer.Add, lexer.Sub)
	infix(20, pratt.binary, lexer.Mul, lexer.Div)
	infixr(30, pratt.binaryr, lexer.Pow)
	infix(40, pratt.binary, lexer.ImpMul)
	prefix(50, pratt.unary, lexer.Add, lexer.Sub)
	infix(60, pratt.call, lexer.OpenParen)
}

// Parser API: inputs string, outputs either AST or Error
func Parse(s string) (Node, error) {
	// Transform string into tokens
	ts, err := lexer.Scan(s)
	if err != nil {
		return nil, err
	}
	// Set parser state.
	pratt = parser{
		src:   ts,
		index: 0,
		end:   len(ts) - 1,
		table: pratt.table, // Reuse lookup table from package initialization.
	}
	// Weave tokens into abstract syntax tree.
	node, err := pratt.expression(0)
	if err != nil {
		return nil, err
	}
	return node, nil
}
