package table

import (
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Table struct {
	headers []string
	rows    [][]string
	writer  io.Writer
	style   table.Style
}

func New(headers ...string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
		writer:  os.Stdout,
		style:   StyleLight(),
	}
}

func (t *Table) SetWriter(w io.Writer) {
	t.writer = w
}

func (t *Table) SetStyle(style table.Style) {
	t.style = style
}

func (t *Table) AddRow(values ...string) {
	t.rows = append(t.rows, values)
}

func (t *Table) Render() {
	tw := table.NewWriter()
	tw.SetOutputMirror(t.writer)
	tw.SetStyle(t.style)

	if len(t.headers) > 0 {
		headerRow := make(table.Row, len(t.headers))
		for i, h := range t.headers {
			headerRow[i] = h
		}
		tw.AppendHeader(headerRow)
	}

	for _, row := range t.rows {
		tableRow := make(table.Row, len(row))
		for i, cell := range row {
			tableRow[i] = cell
		}
		tw.AppendRow(tableRow)
	}

	tw.Render()
}

func StyleLight() table.Style {
	style := table.StyleLight
	style.Color.Header = text.Colors{text.FgHiCyan, text.Bold}
	style.Options.DrawBorder = true
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = false
	return style
}

func StyleRounded() table.Style {
	style := table.StyleRounded
	style.Color.Header = text.Colors{text.FgHiCyan, text.Bold}
	style.Options.DrawBorder = true
	style.Options.SeparateHeader = true
	return style
}

func StyleDouble() table.Style {
	style := table.StyleDouble
	style.Color.Header = text.Colors{text.FgHiCyan, text.Bold}
	style.Options.DrawBorder = true
	style.Options.SeparateHeader = true
	return style
}

func StyleColoredBright() table.Style {
	return table.StyleColoredBright
}
