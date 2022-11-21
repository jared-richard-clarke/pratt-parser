package parser

import (
	"fmt"
	"github/jared-richard-clarke/pratt/lexer"
	"strconv"
)

// === Under Heavy Construction ===

type nud func(lexer.Token) Node       // Null denotation
type led func(Node, lexer.Token) Node // Left denotation

type table struct {
	nuds map[lexer.LexType]nud // lexeme -> Nud
	leds map[lexer.LexType]led // lexeme -> Led

	rbp map[lexer.LexType]int // lexeme -> right (prefix) binding power
	lbp map[lexer.LexType]int // lexeme -> left (infix) binding power
}

type parser struct {
	src    []lexer.Token
	length int
	index  int
	table
}

func (p *parser) next() lexer.Token {
	t := p.src[p.index]
	p.index += 1
	return t
}

func (p *parser) expression(rbp int) Node {
	token := p.next()
	nud, ok := p.table.nuds[token.Typeof]
	if !ok {
		panic(fmt.Errorf("could not parse %s", token.Value))
	}
	left := nud(token)
	for rbp < p.table.lbp[token.Typeof] {
		token := p.next()
		led := p.table.leds[token.Typeof]
		left = led(left, token)
	}
	return left
}

func (p *parser) parseNumber(t lexer.Token) Node {
	num, _ := strconv.ParseFloat(t.Value, 64)
	return &Number{
		Value: num,
	}
}

func (p *parser) parseIdent(t lexer.Token) Node {
	return &Ident{
		Value: t.Value,
	}
}

func (p *parser) parseUnary(t lexer.Token) Node {
	x := p.expression(p.rbp[t.Typeof])
	return &Unary{
		Op: t.Typeof,
		X:  x,
	}
}

func (p *parser) parseBinary(left Node, t lexer.Token) Node {
	right := p.expression(p.table.lbp[t.Typeof])
	return &Binary{
		Op: t.Typeof,
		X:  left,
		Y:  right,
	}
}

func (p *parser) parseBinaryRight(left Node, t lexer.Token) Node {
	right := p.expression(p.table.lbp[t.Typeof] - 1)
	return &Binary{
		Op: t.Typeof,
		X:  left,
		Y:  right,
	}
}
