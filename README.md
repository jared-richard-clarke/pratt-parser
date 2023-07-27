# Pratt Parser

I built a Pratt Parser for arithmetic and symbolic expressions.
The parser inputs a string and outputs an abstract syntax tree.
Robust error handling is a priority, so extra-textual information 
is provided, both in the abstract syntax tree and in error messaging.

This implementation is currently written in Go, although I might
expand into other languages.

## Sources

Many thanks to these excellent resources.

| Article | Author |
| :---    | :---   |
| [Top Down Operator Precedence](https://tdop.github.io/) | Vaughan R. Pratt  |
| [Top Down Operator Precedence](https://crockford.com/javascript/tdop/tdop.html) | Douglas Crockford |
| [Top-Down operator precedence parsing](https://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing) | Eli Bendersky |
| [Pratt Parsers: Expression Parsing Made Easy](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/) | Bob Nystrom |

## Example

```go
// === imports omitted ===
func main() {
    text := "sum(7, 11x)"
    node, _ := parser.Parse(text)
    fmt.Print(parser.PrettyPrint(&node))
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
