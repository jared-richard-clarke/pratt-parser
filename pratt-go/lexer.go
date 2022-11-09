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
		return fmt.Sprintf("%s", t.value)
	case t.typeof < number:
		return fmt.Sprintf("%c", t.value)
	default:
		return fmt.Sprintf("%d", t.value)
	}
}

type scanner struct {
	source  string
	tokens  []*token
	length  int // The number of bytes in source string.
	unknown bool

	offset int // The total offset from the beginning of a string. Counts runes by byte size.
	start  int // The start of a lexeme within source string. Counts runes by byte size.
	line   int // Counts newlines ('\n').
	column int // The start of a lexeme within a newline. Counts runes by 1.
}

func (sc *scanner) isAtEnd() bool {
	return sc.offset >= sc.length
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
	if offset >= sc.length {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(sc.source[offset:])
	return r
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
	switch {
	case r == whiteSpace, r == carriageReturn, r == tab:
		return
	case r == newline:
		sc.column = 0
		sc.line += 1
		return
	case unicode.IsDigit(r):
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
	case r == '(':
		sc.addToken(oParen, r)
		return
	case r == ')':
		sc.addToken(cParen, r)
		return
	case r == '-':
		sc.addToken(sub, r)
		return
	case r == '+':
		sc.addToken(add, r)
		return
	case r == '*':
		sc.addToken(mul, r)
		return
	case r == '/':
		sc.addToken(div, r)
		return
	default:
		sc.unknown = true
		sc.addToken(unknown, fmt.Errorf("unknown: %c", r))
		return
	}
}

// The Lexer API: drives the scanner.
func Scan(t string) ([]*token, bool) {
	sc := scanner{
		source:  t,
		tokens:  make([]*token, 0),
		length:  len(t),
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
