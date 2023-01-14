package parser

import (
	"fmt"
	"strings"
)

// Under heavy construction.

const newline = "\n"

type printer struct {
	spacer  string
	level   int
	padding string
	output  strings.Builder
}

func (p *printer) write(x string) { p.output.WriteString(x) }

func (p *printer) writepad(xs ...string) {
	for _, x := range xs {
		p.output.WriteString(p.padding)
		p.output.WriteString(x)
	}
}

func (p *printer) indent() {
	p.level += 1
	p.padding = strings.Repeat(p.spacer, p.level)
}
func (p *printer) outdent() {
	p.level -= 1
	p.padding = strings.Repeat(p.spacer, p.level)
}

func (p *printer) print() string { return p.output.String() }

func (p *printer) format(n *Node) {
	nl := func(s string) string { return s + newline }
	li := func(i int) string { return fmt.Sprintf("Line:   %d%s", i, newline) }
	co := func(i int) string { return fmt.Sprintf("Column: %d%s", i, newline) }
	close := nl("}")

	switch n := (*n).(type) {
	case Number:
		label := nl("Number{")
		value := fmt.Sprintf("Value:  %g%s", n.Value, newline)
		line := li(n.Line)
		column := co(n.Column)

		p.write(label)
		p.indent()
		p.writepad(value, line, column)
		p.outdent()
		p.writepad(close)
	case Symbol:
		label := nl("Symbol{")
		value := fmt.Sprintf("Value:  %s%s", n.Value, newline)
		line := li(n.Line)
		column := co(n.Column)

		p.write(label)
		p.indent()
		p.writepad(value, line, column)
		p.outdent()
		p.writepad(close)
	case Unary:
		label := nl("Unary{")
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)
		line := li(n.Line)
		column := co(n.Column)

		p.write(label)
		p.indent()
		p.writepad(op)
		p.writepad("X: ")
		p.format(&n.X)
		p.writepad(line, column)
		p.outdent()
		p.writepad(close)
	case Binary:
		label := nl("Binary{")
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)
		line := li(n.Line)
		column := co(n.Column)

		p.write(label)
		p.indent()
		p.writepad(op)
		p.writepad("X: ")
		p.format(&n.X)
		p.writepad("Y: ")
		p.format(&n.Y)
		p.writepad(line, column)
		p.outdent()
		p.writepad(close)
	case ImpliedBinary:
		label := nl("ImpliedBinary{")
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)

		p.write(label)
		p.indent()
		p.writepad(op)
		p.writepad("X: ")
		p.format(&n.X)
		p.writepad("Y: ")
		p.format(&n.Y)
		p.outdent()
		p.writepad(close)
	case Call:
		label := nl("Call{")
		line := li(n.Line)
		column := co(n.Column)

		p.write(label)
		p.indent()
		p.writepad("Callee: ")
		p.format(&n.Callee)
		p.writepad(nl("Args: ["))
		p.indent()
		if len(n.Args) == 0 {
			p.writepad(nl("[]"))
		} else {
			for _, arg := range n.Args {
				p.writepad("") // pad each argument
				p.format(&arg)
			}
		}
		p.outdent()
		p.writepad(nl("]"))
		p.writepad(line, column)
		p.outdent()
		p.writepad(close)
	default:
		p.write("Empty{}")
	}
}

// Inputs a pointer to a Node and outputs a formatted string of that Node.
func PrettyPrint(n *Node) string {
	var b strings.Builder
	p := printer{
		spacer:  strings.Repeat(" ", 4),
		level:   0,
		padding: "",
		output:  b,
	}
	p.format(n)
	return p.print()
}
