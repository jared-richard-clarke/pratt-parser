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

type LexType int

const (
	OpenParen LexType = iota
	CloseParen
	Comma
	Add
	Sub
	Mul
	ImpMul // implicit multiplier
	Div
	Pow
	Number
	Symbol
	EOF
)

type Token struct {
	Typeof LexType // Lexeme type, denoted by "LexType".
	Value  string  // Lexeme string value.
	Line   int     // Lexeme line number. Counts newlines ('\n').
	Column int     // Lexeme starting column within newline. Counts runes.
}

// Helper functions and constants

// Whereas lexer uses "EOF" to mark the end of an array of tokens,
// lexer uses "eof" internally to signal the end of a string or file.
// "eof" is the untyped int -1, which has no rune alias.
const eof = -1

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == underscore
}

// Internal scanner and methods

type scanner struct {
	source string  // Scanner input. Currently a string.
	tokens []Token // Array slice of accumulating tokens.
	length int     // Number of bytes in the source string.

	offset int // Total string offset. Counts bytes.
	start  int // Start of a lexeme within source string. Counts bytes.
	line   int // Counts newlines ('\n').
	column int // Tracks the start of a lexeme within a newline. Counts runes.
}

func (sc *scanner) end() bool {
	return sc.offset >= sc.length
}

// Skips whitespace: '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
func (sc *scanner) skip() {
	for {
		r, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
		if !unicode.IsSpace(r) {
			break
		}
		sc.column += 1
		sc.offset += w
		if r == newline {
			sc.line += 1
			sc.column = 0
		}
	}
}

func (sc *scanner) next() rune {
	r, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
	sc.column += 1
	sc.offset += w
	return r
}

func (sc *scanner) peek() rune {
	if sc.end() {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(sc.source[sc.offset:])
	return r
}

func (sc *scanner) peekNext() rune {
	_, w := utf8.DecodeRuneInString(sc.source[sc.offset:])
	offset := sc.offset + w
	if offset >= sc.length {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(sc.source[offset:])
	return r
}

func (sc *scanner) addToken(t LexType, v string) {
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
	})
}

func (sc *scanner) scanToken() error {
	r := sc.next()
	switch {
	// whitespace
	case r == whiteSpace, r == carriageReturn, r == tab:
		return nil
	case r == newline:
		sc.column = 0
		sc.line += 1
		return nil
	// punctuators
	case r == '(':
		sc.addToken(OpenParen, "(")
		return nil
	case r == ')':
		sc.addToken(CloseParen, ")")
		// Check for implied multiplication: (7+11)x, (7+11)(11+7), or (7+11)7
		sc.skip()
		c := sc.peek()
		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '(' {
			sc.addToken(ImpMul, "*")
		}
		return nil
	case r == ',':
		sc.addToken(Comma, ",")
		return nil
	case r == '-':
		sc.addToken(Sub, "-")
		return nil
	case r == '+':
		sc.addToken(Add, "+")
		return nil
	case r == '*' || r == '×':
		sc.addToken(Mul, "*")
		return nil
	case r == '/' || r == '÷':
		sc.addToken(Div, "/")
		return nil
	case r == '^':
		sc.addToken(Pow, "^")
		return nil
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
		// Check for implied multiplication: 7x or 7(7+11)
		sc.skip()
		c := sc.peek()
		if unicode.IsLetter(c) || c == '(' {
			sc.addToken(ImpMul, "*")
		}
		return nil
	// symbols
	case unicode.IsLetter(r):
		for isAlphaNumeric(sc.peek()) {
			sc.next()
		}
		text := sc.source[sc.start:sc.offset]
		sc.addToken(Symbol, text)
		return nil
	// undefined
	default:
		msg := "unexpected character: %q line:%d, column:%d"
		return fmt.Errorf(msg, r, sc.line, sc.column)
	}
}

// The Lexer API: drives the scanner.
func Scan(t string) ([]Token, error) {
	sc := scanner{
		source: t,
		tokens: make([]Token, 0),
		length: len(t),
		offset: 0,
		start:  0,
		line:   1,
		column: 0,
	}
	for !sc.end() {
		sc.start = sc.offset
		if err := sc.scanToken(); err != nil {
			return nil, err
		}
	}
	sc.tokens = append(sc.tokens, Token{
		Typeof: EOF,
		Line:   sc.line,
		Column: sc.column + 1,
	})
	return sc.tokens, nil
}
