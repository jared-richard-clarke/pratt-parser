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

func (p *printer) write(s string) {
	p.output.WriteString(s)
}

func (p *printer) writepad(s string) {
	p.output.WriteString(p.padding)
	p.output.WriteString(s)
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
	close := "}" + newline
	nl := func(i int) string { return fmt.Sprintf("Line:   %d%s", i, newline) }
	nc := func(i int) string { return fmt.Sprintf("Column: %d%s", i, newline) }

	switch n := (*n).(type) {
	case Number:
		label := "Number{" + newline
		value := fmt.Sprintf("Value:  %g%s", n.Value, newline)
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label)
		p.indent()
		p.writepad(value)
		p.writepad(line)
		p.writepad(column)
		p.outdent()
		p.writepad(close)
	case Symbol:
		label := "Symbol{" + newline
		value := fmt.Sprintf("Value:  %s%s", n.Value, newline)
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label)
		p.indent()
		p.writepad(value)
		p.writepad(line)
		p.writepad(column)
		p.outdent()
		p.writepad(close)
	case Unary:
		label := "Unary{" + newline
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label)
		p.indent()
		p.writepad(op)
		p.writepad("X: ")
		p.format(&n.X)
		p.writepad(line)
		p.writepad(column)
		p.outdent()
		p.writepad(close)
	case Binary:
		label := "Binary{" + newline
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label)
		p.indent()
		p.writepad(op)
		p.writepad("X: ")
		p.format(&n.X)
		p.writepad("Y: ")
		p.format(&n.Y)
		p.writepad(line)
		p.writepad(column)
		p.outdent()
		p.writepad(close)
	case ImpliedBinary:
		label := "ImpliedBinary{" + newline
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
		label := "Call{" + newline
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label)
		p.indent()
		p.writepad("Callee: ")
		p.format(&n.Callee)
		p.writepad("Args: [" + newline)
		p.indent()
		if n.Args == nil {
			p.writepad("<nil>" + newline)
		} else {
			for _, arg := range n.Args {
				p.writepad("") // pad each argument
				p.format(&arg)
			}
		}
		p.outdent()
		p.writepad("]" + newline)
		p.writepad(line)
		p.writepad(column)
		p.outdent()
		p.writepad(close)
	default:
		p.write("Empty{}")
	}
}

// Inputs a pointer to a Node and outputs a formatted string.
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
