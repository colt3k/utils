// Package print wraps tablewriter for simple CLI table output.
package print

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type Printer struct {
	header []string
	data   [][]string
}

// TablePrint renders the configured table to standard output.
func (p *Printer) TablePrint(border bool) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(border)

	table.SetHeader(p.header)

	for _, v := range p.data {
		table.Append(v)
	}
	table.Render() // Send output
}
