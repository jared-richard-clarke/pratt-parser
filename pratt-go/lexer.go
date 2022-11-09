package pratt

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	newline        = '\n'
	carriageReturn = '\r'
	tab            = '\t'
	whiteSpace     = ' '
	decPoint       = '.'
)

type lexType int

const (
	oParen lexType = iota
	cParen
	add
	sub
	mul
	div
	number
	unknown
	eof
)

type token struct {
	typeof lexType
	value  any
	line   int
	column int
}

func (t *token) String() string {
	switch {
	case t.typeof == eof:
		return "EOF"
	case t.typeof == unknown:
		return fmt.Sprintf("Unknown: %s", t.value)
	case t.typeof < number:
		return fmt.Sprintf("%c", t.value)
	default:
		return fmt.Sprintf("%d", t.value)
	}
}

type scanner struct {
	source  string
	tokens  []*token
	unknown bool

	offset int // The total offset from the beginning of a string. Counts runes by byte size.
	start  int // The start of a lexeme within source string. Counts runes by byte size.
	line   int // Counts newlines ('\n').
	column int // The start of a lexeme within a newline. Counts runes by 1.
}

func (sc *scanner) isAtEnd() bool {
	return sc.offset >= len(sc.source)
}

func (sc *scanner) next() rune {
	r, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
	sc.column += 1
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

	sc.column += 1
	sc.offset = sc.start + w
	return true
}

func (sc *scanner) addToken(t lexType, v any) {
	sc.tokens = append(sc.tokens, &token{
		typeof: t,
		value:  v,
		line:   sc.line,
		column: sc.column,
	})
}

func (sc *scanner) scanToken() {
	r := sc.next()
	// numbers
	if unicode.IsDigit(r) {
		for unicode.IsDigit(sc.peek()) {
			sc.next()
		}
		if sc.peek() == decPoint && unicode.IsDigit(sc.peekNext()) {
			sc.next()
			for unicode.IsDigit(sc.peek()) {
				sc.next()
			}
		}
		text := string(sc.source[sc.start:sc.offset])
		num, err := strconv.ParseFloat(text, 64)
		if err != nil {
			err = fmt.Errorf("parsing float %q: %w", text, err)
			sc.addToken(unknown, err)
		} else {
			sc.addToken(number, num)
		}
		return
	}
	// punctuators
	switch r {
	case whiteSpace, carriageReturn, tab:
		return
	case newline:
		sc.column = 0
		sc.line += 1
		return
	case '(':
		sc.addToken(oParen, r)
		return
	case ')':
		sc.addToken(cParen, r)
		return
	case '-':
		sc.addToken(sub, r)
		return
	case '+':
		sc.addToken(add, r)
		return
	case '*':
		sc.addToken(mul, r)
		return
	case '/':
		sc.addToken(div, r)
		return
	default:
		sc.unknown = true
		sc.addToken(unknown, fmt.Sprintf("unknown: %v", r))
		return
	}
}

// The Lexer API: drives the scanner.
func Scan(t string) ([]*token, bool) {
	sc := scanner{
		source:  t,
		tokens:  make([]*token, 0),
		unknown: false,
		offset:  0,
		start:   0,
		line:    1,
		column:  0,
	}
	for !sc.isAtEnd() {
		sc.start = sc.offset
		sc.scanToken()
	}
	sc.tokens = append(sc.tokens, &token{
		typeof: eof,
		value:  nil,
		line:   sc.line,
		column: sc.column + 1,
	})
	return sc.tokens, sc.unknown
}
