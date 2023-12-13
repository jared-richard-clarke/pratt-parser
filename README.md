# Pratt Parser

Contained within in this repository are two experimental, **top-down operator precedence** parsers based
on ideas pioneered by computer scientist **Vaughan R. Pratt**. One is implemented in **Go**, the other in **JavaScript**.

## Sources

I could not have built these parsers on my own. Many thanks to these excellent resources.

| Article | Author |
| :---    | :---   |
| [Top Down Operator Precedence](https://tdop.github.io/) | Vaughan R. Pratt  |
| [Top Down Operator Precedence](https://crockford.com/javascript/tdop/tdop.html) | Douglas Crockford |
| [Top-Down operator precedence parsing](https://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing) | Eli Bendersky |
| [Pratt Parsers: Expression Parsing Made Easy](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/) | Bob Nystrom |

## JavaScript Implementation

JavaScript library scans, parses, and evaluates strings as arithmetic expressions.
The parser does not build an abstract syntax tree. Instead it evaluates the expression as it parses.
It returns a two part array — "[string, null]" if successful, "[null, [token]]" if unsuccessful.
"string" is an evaluated arithmetic expression, and "[token]" is an array of error tokens
both locating and describing errors within the input string.

The library contains both a complementary big number library and a formatter. The big number library provides extended
precision for evaluating floating point numbers within the parsers internal machinery. The formatter is an optional
part of the parser API. If applied, it transforms the parser output into a user-friendly string.

### Example

```JavaScript
import { parse, format } from "parser.js";

// === output ===
parse("0.1 + 0.2"); // ->
// big number: ["0.3", null],
// otherwise:  [0.30000000000000004, null]

// === formatted output ===
format(parse, "1 ÷ 0"); // ->
// 1 ÷ 0
//   ^
// 1. Cannot divide by zero.
```

## Go Implementation

Go library provides a parser for arithmetic and symbolic expressions. The parser inputs a string and outputs
an abstract syntax tree. Error handling is robust. Extra-textual information  is provided, both in the abstract syntax tree
should the parser succeed and in the error messaging should the parser fail. A complementary formatter transforms
the AST into a formatted string.

### Example

```go
// === imports omitted ===
func main() {
    text := "sum(7, 11x)"
    node, _ := parser.Parse(text)
    fmt.Print(parser.Format(&node))
}
// === standard output ===
// Call{
//     Callee: Symbol{
//         Value:  "sum"
//         Line:   1
//         Column: 1
//     }
//     Args: [
//         Number{
//             Value:  7
//             Line:   1
//             Column: 5
//         }
//         Implied Binary{
//             Op: "*"
//             X: Number{
//                    Value:  11
//                    Line:   1
//                    Column: 8
//             }
//             Y: Symbol{
//                    Value:  "x"
//                    Line:   1
//                    Column: 10
//             }
//         }
//     ]
//     Line:   1
//     Column: 4
// }
```
