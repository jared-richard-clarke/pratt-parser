package pratt

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	newline        = '\n'
	formFeed       = '\f'
	carriageReturn = '\r'
	tab            = '\t'
	whiteSpace     = ' '
	decimalPoint   = '.'
	underscore     = '_'
)

type lexType int

const (
	lexOpenParen lexType = iota
	lexCloseParen
	lexAdd
	lexSub
	lexMul
	lexDiv
	lexExp
	lexNumber
	lexIdent
	lexError
	lexEOF
)

type token struct {
	typeof lexType
	value  any
	line   int
	column int
}

func (t *token) String() string {
	switch {
	case t.typeof < lexNumber:
		return fmt.Sprintf("Punctuator: %c", t.value)
	case t.typeof == lexNumber:
		return fmt.Sprintf("Float64: %f", t.value)
	case t.typeof == lexIdent:
		return fmt.Sprintf("Identifier: %s", t.value)
	case t.typeof == lexError:
		return fmt.Sprintf("Error: %s", t.value)
	case t.typeof == lexEOF:
		return "EOF"
	default:
		return fmt.Sprintf("%v", t.value)
	}
}

type scanner struct {
	source string
	tokens []token
	length int // The number of bytes in the source string. Compute only once.

	offset int // The total offset from the beginning of a string. Counts runes by byte size.
	start  int // The start of a lexeme within source string. Counts runes by byte size.
	line   int // Counts newlines ('\n').
	column int // Tracks the start of a lexeme within a newline. Counts runes by 1.
}

func (sc *scanner) isAtEnd() bool {
	return sc.offset >= sc.length
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == underscore
}

func (sc *scanner) next() rune {
	r, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
	sc.column += 1
	sc.offset += w
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
	runeCount := utf8.RuneCountInString(sc.source[sc.start:sc.offset])
	lexOffset := 0
	// Center column on character representation of rune.
	// [ w, o, r, d ] not [ w, o, r, d ]
	//   ^                 ^
	if runeCount > 1 {
		lexOffset = runeCount - 1
	}
	// Ensure column count begins at start of lexeme.
	column := sc.column - lexOffset

	sc.tokens = append(sc.tokens, token{
		typeof: t,
		value:  v,
		line:   sc.line,
		column: column,
	})
}

func (sc *scanner) scanToken() {
	r := sc.next()
	switch {
	// whitespace
	case r == whiteSpace, r == formFeed, r == carriageReturn, r == tab:
		return
	case r == newline:
		sc.column = 0
		sc.line += 1
		return
	// punctuators
	case r == '(':
		sc.addToken(lexOpenParen, r)
		return
	case r == ')':
		sc.addToken(lexCloseParen, r)
		return
	case r == '-':
		sc.addToken(lexSub, r)
		return
	case r == '+':
		sc.addToken(lexAdd, r)
		return
	case r == '*':
		sc.addToken(lexMul, r)
		return
	case r == '/':
		sc.addToken(lexDiv, r)
		return
	case r == '^':
		sc.addToken(lexExp, r)
		return
	// numbers
	case unicode.IsDigit(r):
		for unicode.IsDigit(sc.peek()) {
			sc.next()
		}
		if sc.peek() == decimalPoint && unicode.IsDigit(sc.peekNext()) {
			sc.next()
			for unicode.IsDigit(sc.peek()) {
				sc.next()
			}
		}
		text := sc.source[sc.start:sc.offset]
		num, err := strconv.ParseFloat(text, 64)
		if err != nil {
			sc.addToken(lexError, fmt.Sprintf("parsing float %q: %s", text, err.Error()))
		} else {
			sc.addToken(lexNumber, num)
		}
		return
	// identifiers
	case unicode.IsLetter(r):
		for isAlphaNumeric(sc.peek()) {
			sc.next()
		}
		text := sc.source[sc.start:sc.offset]
		sc.addToken(lexIdent, text)
		return
	// unknowns
	default:
		sc.addToken(lexError, fmt.Sprintf("%c", r))
		return
	}
}

// The Lexer API: drives the scanner.
func Scan(t string) []token {
	sc := scanner{
		source: t,
		tokens: make([]token, 0),
		length: len(t),
		offset: 0,
		start:  0,
		line:   1,
		column: 0,
	}
	for !sc.isAtEnd() {
		sc.start = sc.offset
		sc.scanToken()
	}
	sc.tokens = append(sc.tokens, token{
		typeof: lexEOF,
		line:   sc.line,
		column: sc.column + 1,
	})
	return sc.tokens
}
