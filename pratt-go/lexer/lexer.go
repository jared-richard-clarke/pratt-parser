package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const (
	newline        = '\n'
	carriageReturn = '\r'
	tab            = '\t'
	whiteSpace     = ' '
	decimalPoint   = '.'
	underscore     = '_'
)

type lexType int

const (
	OpenParen lexType = iota
	CloseParen
	Add
	Sub
	Mul
	Div
	Exp
	Number
	Ident
	Error
	EOF
)

type Token struct {
	Typeof lexType
	Value  string
	Line   int
	Column int
	Length int
}

func (t Token) String() string {
	switch {
	case t.Typeof < Number:
		return fmt.Sprintf("Punct: %q :%d:%d:%d", t.Value, t.Line, t.Column, t.Length)
	case t.Typeof == Number:
		return fmt.Sprintf("Float: %q :%d:%d:%d", t.Value, t.Line, t.Column, t.Length)
	case t.Typeof == Ident:
		return fmt.Sprintf("Ident: %q :%d:%d:%d", t.Value, t.Line, t.Column, t.Length)
	case t.Typeof == Error:
		return fmt.Sprintf("Error: %q :%d:%d:%d", t.Value, t.Line, t.Column, t.Length)
	case t.Typeof == EOF:
		return fmt.Sprintf("EOF :%d:%d:%d", t.Line, t.Column, t.Length)
	default:
		return fmt.Sprintf("Undef: %q :%d:%d:%d", t.Value, t.Line, t.Column, t.Length)
	}
}

type scanner struct {
	source string
	tokens []Token
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

func (sc *scanner) addToken(t lexType, v string) {
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

	sc.tokens = append(sc.tokens, Token{
		Typeof: t,
		Value:  v,
		Line:   sc.line,
		Column: column,
		Length: runeCount,
	})
}

func (sc *scanner) scanToken() {
	r := sc.next()
	switch {
	// whitespace
	case r == whiteSpace, r == carriageReturn, r == tab:
		return
	case r == newline:
		sc.column = 0
		sc.line += 1
		return
	// punctuators
	case r == '(':
		sc.addToken(OpenParen, "(")
		return
	case r == ')':
		sc.addToken(CloseParen, ")")
		return
	case r == '-':
		sc.addToken(Sub, "-")
		return
	case r == '+':
		sc.addToken(Add, "+")
		return
	case r == '*':
		sc.addToken(Mul, "*")
		return
	case r == '/':
		sc.addToken(Div, "/")
		return
	case r == '^':
		sc.addToken(Exp, "^")
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
		sc.addToken(Number, text)
		return
	// identifiers
	case unicode.IsLetter(r):
		for isAlphaNumeric(sc.peek()) {
			sc.next()
		}
		text := sc.source[sc.start:sc.offset]
		sc.addToken(Ident, text)
		return
	// undefined
	default:
		sc.addToken(Error, string(r))
		return
	}
}

// The Lexer API: drives the scanner.
func Scan(t string) []Token {
	sc := scanner{
		source: t,
		tokens: make([]Token, 0),
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
	sc.tokens = append(sc.tokens, Token{
		Typeof: EOF,
		Line:   sc.line,
		Column: sc.column + 1,
		Length: 1,
	})
	return sc.tokens
}
