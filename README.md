# Pratt Parser

I'm building a Pratt Parser for arithmetic and symbolic expressions.
The parser will input a string and output an abstract syntax tree.
Robust error handling is a priority, so extra-textual information 
is provided, both in the abstract syntax tree and error messages.

This implementation is currently written in Go, although I will probably
expand into other languages.

## Sources

Many thanks to these excellent resources.

| Article | Author |
| :---    | :---   |
| [Top Down Operator Precedence](https://tdop.github.io/) | Vaughan R. Pratt  |
| [Top Down Operator Precedence](https://crockford.com/javascript/tdop/tdop.html) | Douglas Crockford |
| [Top-Down operator precedence parsing](https://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing) | Eli Bendersky |
| [Pratt Parsers: Expression Parsing Made Easy](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/) | Bob Nystrom |

## Bugs

- `parser.Parse("7.0.1")` returns a Literal node of `7.0`. It should return an `error` instead.
- `parser.Parse("-7-")` returns an `error`, which reads `undefined NUD for EOF :1:4`. 
   Although true, this error makes little sense from a user's perspective.
