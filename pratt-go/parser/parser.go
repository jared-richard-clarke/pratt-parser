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
	table                // parser and binding lookup
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

func (p *parser) peek() lexer.LexType {
	return p.src[p.index].Typeof
}

func (p *parser) match(e lexer.LexType) bool {
	g := p.src[p.index].Typeof
	if g != e {
		return false
	}
	p.next()
	return true
}

func (p *parser) parseExpr(rbp int) (Node, error) {
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

func (p *parser) parseNumber(t lexer.Token) (Node, error) {
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

func (p *parser) parseIdent(t lexer.Token) (Node, error) {
	return &Ident{
		Value:  t.Value,
		Line:   t.Line,
		Column: t.Column,
	}, nil
}

func (p *parser) parseUnary(t lexer.Token) (Node, error) {
	x, err := p.parseExpr(p.rbp[t.Typeof])
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

func (p *parser) parseBinary(left Node, t lexer.Token) (Node, error) {
	right, err := p.parseExpr(p.lbp[t.Typeof])
	if err != nil {
		return nil, err
	}
	return &Binary{
		Op:     t.Value,
		X:      left,
		Y:      right,
		Line:   t.Line,
		Column: t.Column,
	}, nil
}

func (p *parser) parseBinaryR(left Node, t lexer.Token) (Node, error) {
	right, err := p.parseExpr(p.lbp[t.Typeof] - 1)
	if err != nil {
		return nil, err
	}
	return &Binary{
		Op:     t.Value,
		X:      left,
		Y:      right,
		Line:   t.Line,
		Column: t.Column,
	}, nil
}

func (p *parser) parseParenExpr(t lexer.Token) (Node, error) {
	x, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	if !p.match(lexer.CloseParen) {
		return nil, fmt.Errorf("expected ')', got '%s'", t.Value)
	}
	return x, nil
}
