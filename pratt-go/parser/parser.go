package parser

import (
	"fmt"
	"github/jared-richard-clarke/pratt/internal/lexer"
	"strconv"
)

// Top down operator precedence parsing, as imagined by Vaughan Pratt,
// combines lexical semantics with functions. Each lexeme is assigned a
// function — its semantic code. To parse a string of lexemes is to execute
// the semantic code of each lexeme in turn from left to right.
//
// Sidenote: semantic code, in this case, is very much like a method
// in object-oriented programming.
//
// There are two types of semantic code:
// 1. null denotation ( nud ): a lexeme without a left expression.
// 2. left denotation ( led ): a lexeme with a left expression.

// Parses numbers, symbols, and unary operators
type nud func(lexer.Token) (Node, error)

// Parses binary operators
type led func(Node, lexer.Token) (Node, error)

type table struct {
	nud  map[lexer.LexType]nud // lexeme -> nud
	led  map[lexer.LexType]led // lexeme -> led
	bind map[lexer.LexType]int // lexeme -> binding power
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

// The engine of Pratt's technique, "parseExpression" drives the parser,
// calling the semantic code of each lexeme in turn from left to right.
// For every level of precedence — dictated by position and binding power —
// there is a call to "parseExpression" either through the "nud" or "led"
// of the associated lexeme. The resolution of "parseExpression" is to
// return either the branch of an abstract syntax tree or an error.
func (p *parser) parseExpression(rbp int) (Node, error) {
	token := p.next()
	nud, ok := p.nud[token.Typeof]
	if !ok {
		msg := "undefined prefix operation %q line:%d column:%d"
		return nil, fmt.Errorf(msg, token.Value, token.Line, token.Column)
	}
	left, err := nud(token)
	if err != nil {
		return nil, err
	}
	for rbp < p.bind[p.peek()] {
		token := p.next()
		led, ok := p.led[token.Typeof]
		if !ok {
			msg := "undefined infix operation %q line:%d column:%d"
			return nil, fmt.Errorf(msg, token.Value, token.Line, token.Column)
		}
		left, err = led(left, token)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

// Parses either empty or incomplete expressions.
func (p *parser) parseEOF(token lexer.Token) (Node, error) {
	if len(p.src) == 1 {
		return Empty{}, nil
	}
	msg := "incomplete parseExpression, unexpected <EOF> line:%d column:%d"
	err := fmt.Errorf(msg, token.Line, token.Column)
	return nil, err
}

// Parses numbers as 64-bit floating point.
func (p *parser) parseNumber(token lexer.Token) (Node, error) {
	num, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		msg := "invalid number: %s line:%d column:%d"
		return nil, fmt.Errorf(msg, token.Value, token.Line, token.Column)
	}
	return Number{
		Value:  num,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

// Parses symbols — otherwise known as identifiers.
// Always returns Node. Has error type to satisfy "nud".
func (p *parser) parseSymbol(token lexer.Token) (Node, error) {
	return Symbol{
		Value:  token.Value,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

// Parses unary expressions.
func (p *parser) parseUnary(token lexer.Token) (Node, error) {
	node, err := p.parseExpression(p.bind[token.Typeof])
	if err != nil {
		return nil, err
	}
	return Unary{
		Op:     token.Value,
		X:      node,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

// Parses binary expressions that associate left.
func (p *parser) parseBinaryLeft(left Node, token lexer.Token) (Node, error) {
	right, err := p.parseExpression(p.bind[token.Typeof])
	if err != nil {
		return nil, err
	}
	if token.Typeof == lexer.ImpMul {
		return ImpliedBinary{
			Op: token.Value,
			X:  left,
			Y:  right,
		}, nil
	}
	return Binary{
		Op:     token.Value,
		X:      left,
		Y:      right,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

// Parses binary expressions that associate right.
func (p *parser) parseBinaryRight(left Node, token lexer.Token) (Node, error) {
	right, err := p.parseExpression(p.bind[token.Typeof] - 1)
	if err != nil {
		return nil, err
	}
	return Binary{
		Op:     token.Value,
		X:      left,
		Y:      right,
		Line:   token.Line,
		Column: token.Column,
	}, nil
}

// Parses parenthetical expressions.
func (p *parser) parseGrouping(token lexer.Token) (Node, error) {
	position := fmt.Sprintf("line:%d column:%d", token.Line, token.Column)
	node, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	if !p.match(lexer.CloseParen) {
		msg := "for '(' %s, missing matching ')'"
		return nil, fmt.Errorf(msg, position)
	}
	p.next()
	return node, nil
}

// Parses function calls.
func (p *parser) parseCall(left Node, token lexer.Token) (Node, error) {
	// For now, the only valid function callees are symbols.
	s, ok := left.(Symbol)
	if !ok {
		msg := "%s is not a callable function"
		return nil, fmt.Errorf(msg, left)
	}
	if p.match(lexer.CloseParen) {
		p.next()
		return Call{
			Callee: left,
			Args:   make([]Node, 0), // Make an empty slice, not a nil slice. Makes comparisons simpler.
			Line:   token.Line,
			Column: token.Column,
		}, nil
	}
	var args []Node
	for {
		node, err := p.parseExpression(0)
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
		msg := "for function call %q, missing closing ')'"
		return nil, fmt.Errorf(msg, s.Value)
	}
	p.next()
	return Call{
		Callee: left,
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
			nud:  make(map[lexer.LexType]nud),
			led:  make(map[lexer.LexType]led),
			bind: make(map[lexer.LexType]int),
		},
	}
	// Helper functions build lookup tables.
	set := func(t lexer.LexType, n nud) {
		pratt.nud[t] = n
		pratt.bind[t] = 0
	}
	prefix := func(n nud, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.nud[t] = n
		}
	}
	affix := func(bp int, l led, ts ...lexer.LexType) {
		for _, t := range ts {
			pratt.led[t] = l
			pratt.bind[t] = bp
		}
	}

	// Initialize lookup table.
	set(lexer.EOF, pratt.parseEOF)
	set(lexer.Number, pratt.parseNumber)
	set(lexer.Symbol, pratt.parseSymbol)
	set(lexer.OpenParen, pratt.parseGrouping)
	prefix(pratt.parseUnary, lexer.Add, lexer.Sub)
	affix(10, pratt.parseBinaryLeft, lexer.Equal, lexer.NotEqual)
	affix(20, pratt.parseBinaryLeft, lexer.Add, lexer.Sub)
	affix(30, pratt.parseBinaryLeft, lexer.Mul, lexer.Div)
	affix(40, pratt.parseBinaryLeft, lexer.ImpMul)
	affix(50, pratt.parseBinaryRight, lexer.Pow)
	affix(60, pratt.parseCall, lexer.OpenParen)
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
	node, err := pratt.parseExpression(0)
	if err != nil {
		return nil, err
	}
	// If unused tokens following expression, return error.
	if pratt.index < pratt.end {
		token := pratt.src[pratt.index]
		msg := "starting line:%d, column:%d, unused tokens following expression"
		return nil, fmt.Errorf(msg, token.Line, token.Column)
	}
	return node, nil
}
