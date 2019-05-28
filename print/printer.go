package print

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type Printer struct {
	header []string
	data   [][]string
}

func (p *Printer) TablePrint(border bool) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(border)

	table.SetHeader(p.header)

	for _, v := range p.data {
		table.Append(v)
	}
	table.Render() // Send output
}
