package table

import (
	"io"
	"os"

	"tracker/internal/ui"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Table struct {
	headers      []string
	rows         [][]string
	writer       io.Writer
	style        table.Style
	columnWidths map[int]int
}

func New(headers ...string) *Table {
	return &Table{
		headers:      headers,
		rows:         make([][]string, 0),
		writer:       os.Stdout,
		style:        StyleLight(),
		columnWidths: make(map[int]int),
	}
}

func (t *Table) SetWriter(w io.Writer) {
	t.writer = w
}

func (t *Table) SetStyle(style table.Style) {
	t.style = style
}

func (t *Table) SetColumnWidths(widths map[int]int) {
	t.columnWidths = widths
}

func (t *Table) AddRow(values ...string) {
	t.rows = append(t.rows, values)
}

func (t *Table) Render() {
	tw := table.NewWriter()
	tw.SetOutputMirror(t.writer)
	tw.SetStyle(t.style)

	if len(t.columnWidths) > 0 {
		columnConfigs := make([]table.ColumnConfig, 0, len(t.columnWidths))
		for colNum, width := range t.columnWidths {
			columnConfigs = append(columnConfigs, table.ColumnConfig{
				Number:   colNum + 1,
				WidthMax: width,
				WidthMin: 0,
			})
		}
		tw.SetColumnConfigs(columnConfigs)
	}

	if len(t.headers) > 0 {
		headerRow := make(table.Row, len(t.headers))
		for i, h := range t.headers {
			if ui.ShouldUseColors() {
				headerRow[i] = color.New(color.FgHiCyan, color.Bold).Sprint(h)
			} else {
				headerRow[i] = h
			}
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
	style.Options.DrawBorder = true
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = false
	style.Format.Header = text.FormatUpper
	style.Format.Row = text.FormatDefault
	style.Format.Footer = text.FormatUpper

	style.Color.Header = text.Colors{}

	return style
}

func StyleRounded() table.Style {
	style := table.StyleRounded
	style.Options.DrawBorder = true
	style.Options.SeparateHeader = true
	style.Color.Header = text.Colors{}
	return style
}

func StyleDouble() table.Style {
	style := table.StyleDouble
	style.Options.DrawBorder = true
	style.Options.SeparateHeader = true
	style.Color.Header = text.Colors{}
	return style
}

func StyleColoredBright() table.Style {
	return table.StyleColoredBright
}
