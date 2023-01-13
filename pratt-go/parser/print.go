package parser

import (
	"fmt"
	"strings"
)

// Under heavy construction

const newline = "\n"

type printer struct {
	spacer string
	level  int
	output strings.Builder
}

func (p *printer) write(s string, pad bool) {
	if pad {
		p.output.WriteString(strings.Repeat(p.spacer, p.level))
	}
	p.output.WriteString(s)
}

func (p *printer) indent()  { p.level += 1 }
func (p *printer) outdent() { p.level -= 1 }

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

		p.write(label, false)
		p.indent()
		p.write(value, true)
		p.write(line, true)
		p.write(column, true)
		p.outdent()
		p.write(close, true)
	case Symbol:
		label := "Symbol{" + newline
		value := fmt.Sprintf("Value:  %s%s", n.Value, newline)
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label, false)
		p.indent()
		p.write(value, true)
		p.write(line, true)
		p.write(column, true)
		p.outdent()
		p.write(close, true)
	case Unary:
		label := "Unary{" + newline
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label, false)
		p.indent()
		p.write(op, true)

		p.write("X: ", true)
		p.format(&n.X)

		p.write(line, true)
		p.outdent()
		p.write(column, true)
		p.write(close, true)
	case Binary:
		label := "Binary{" + newline
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label, false)
		p.indent()
		p.write(op, true)

		p.write("X: ", true)
		p.format(&n.X)

		p.write("Y: ", true)
		p.format(&n.Y)

		p.write(line, true)
		p.write(column, true)
		p.outdent()
		p.write(close, true)
	case ImpliedBinary:
		label := "ImpliedBinary{" + newline
		op := fmt.Sprintf("Op: %q%s", n.Op, newline)

		p.write(label, false)
		p.indent()
		p.write(op, true)

		p.write("X: ", true)
		p.format(&n.X)

		p.write("Y: ", true)
		p.format(&n.Y)

		p.outdent()
		p.write(close, true)
	case Call:
		label := "Call{" + newline
		line := nl(n.Line)
		column := nc(n.Column)

		p.write(label, false)
		p.indent()

		p.write("Callee: ", true)
		p.format(&n.Callee)

		p.write("Args: ["+newline, true)
		p.indent()
		if n.Args == nil {
			p.write("<nil>"+newline, true)
		} else {
			for _, arg := range n.Args {
				p.format(&arg)
			}
		}
		p.outdent()
		p.write("]"+newline, true)

		p.write(line, true)
		p.write(column, true)
		p.outdent()
		p.write(close, true)
	default:
		p.write("Empty{}", false)
	}
}

func print(n *Node) string {
	var b strings.Builder
	p := printer{
		spacer: "    ",
		level:  0,
		output: b,
	}
	p.format(n)
	return p.print()
}
