package lexer

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	Newline        = '\n'
	CarriageReturn = '\r'
	Tab            = '\t'
	WhiteSpace     = ' '
	Comment        = '#'
	DecPoint       = '.'
)

type LexType int

const (
	OParen LexType = iota
	CParen
	Add
	Sub
	Mul
	Div
	Ident
	Number
	Unknown
	EOF
)

type Token struct {
	Typeof LexType
	Value  any

	Line   int
	Column int
}

type scanner struct {
	source  string
	tokens  []*Token
	unknown bool

	offset int // the total offset from the begining of a string
	start  int // the start of a lexeme
	line   int
}

func (sc *scanner) isAtEnd() bool {
	return sc.offset >= len(sc.source)
}

func (sc *scanner) next() rune {
	r, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
	sc.offset = sc.start + w
	return r
}

func (sc *scanner) peek() rune {
	if sc.isAtEnd() {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(sc.source[sc.offset:])
	return r
}

func (sc *scanner) peekNext() rune {
	_, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
	offset := sc.offset + w
	if offset >= len(sc.source) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(sc.source[offset:])
	return r
}

func (sc *scanner) match(e rune) bool {
	if sc.isAtEnd() {
		return false
	}

	r, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
	if r != e {
		return false
	}

	sc.offset = sc.start + w
	return true
}

func (sc *scanner) addToken(t LexType, v any) {
	sc.tokens = append(sc.tokens, &Token{
		Typeof: t,
		Value:  v,
		Line:   sc.line,
		Column: sc.start,
	})
}

func (sc *scanner) scanToken() {
	r := sc.next()
	// numbers
	if unicode.IsDigit(r) {
		for unicode.IsDigit(sc.peek()) {
			sc.next()
		}
		if sc.peek() == DecPoint && unicode.IsDigit(sc.peekNext()) {
			sc.next()
			for unicode.IsDigit(sc.peek()) {
				sc.next()
			}
		}
		text := string(sc.source[sc.start:sc.offset])
		num, err := strconv.ParseFloat(text, 64)
		if err != nil {
			err = fmt.Errorf("parsing float %q: %w", text, err)
			sc.addToken(Unknown, err)
		} else {
			sc.addToken(Number, num)
		}
		return
	}
	// punctuators
	switch r {
	case WhiteSpace, CarriageReturn, Tab:
		return
	case Newline:
		sc.line += 1
		return
	case '(':
		sc.addToken(OParen, r)
		return
	case ')':
		sc.addToken(CParen, r)
		return
	case '-':
		sc.addToken(Sub, r)
		return
	case '+':
		sc.addToken(Add, r)
		return
	case '*':
		sc.addToken(Mul, r)
		return
	case '/':
		sc.addToken(Div, r)
		return
	case Comment:
		for sc.peek() != Newline && !sc.isAtEnd() {
			sc.next()
		}
		return
	default:
		sc.unknown = true
		sc.addToken(Unknown, fmt.Errorf("unknown: %v", r))
		return
	}
}

// The Lexer API: drives the scanner.
func Scan(t string) ([]*Token, bool) {
	sc := scanner{
		source:  t,
		tokens:  make([]*Token, 0),
		unknown: false,
		offset:  0,
		start:   0,
		line:    1,
	}
	for !sc.isAtEnd() {
		sc.start = sc.offset
		sc.scanToken()
	}
	sc.tokens = append(sc.tokens, &Token{
		Typeof: EOF,
		Value:  nil,
		Line:   sc.line,
		Column: sc.offset,
	})
	return sc.tokens, sc.unknown
}
