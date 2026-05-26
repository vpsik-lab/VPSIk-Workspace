package cliui

import (
	"fmt"
	"strings"
)

type Column struct {
	Header string
	Width  int
}

type Table struct {
	columns []Column
	rows    [][]string
}

func NewTable(columns []Column) *Table {
	return &Table{columns: columns}
}

func (t *Table) AddRow(row ...string) {
	t.rows = append(t.rows, row)
}

func (t *Table) Render() string {
	var b strings.Builder

	// Calculate column widths
	widths := make([]int, len(t.columns))
	for i, col := range t.columns {
		w := len(col.Header)
		for _, row := range t.rows {
			if i < len(row) && len(row[i]) > w {
				w = len(row[i])
			}
		}
		if col.Width > 0 && w < col.Width {
			w = col.Width
		}
		widths[i] = w + 2
	}

	// Header
	sep := "╭"
	for i, w := range widths {
		sep += strings.Repeat("─", w)
		sep += "┬"
	}
	sep = sep[:len(sep)-1] + "╮\n"
	b.WriteString(sep)

	line := "│"
	for i, col := range t.columns {
		line += " " + BoldText(fmt.Sprintf("%-*s", widths[i]-2, col.Header))
		line += " │"
	}
	line += "\n"
	b.WriteString(line)

	// Separator
	sep = "├"
	for i, w := range widths {
		sep += strings.Repeat("─", w)
		sep += "┼"
	}
	sep = sep[:len(sep)-1] + "┤\n"
	b.WriteString(sep)

	// Rows
	for _, row := range t.rows {
		line := "│"
		for i, w := range widths {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			line += " " + fmt.Sprintf("%-*s", w-2, val)
			line += " │"
		}
		line += "\n"
		b.WriteString(line)
	}

	// Footer
	sep = "╰"
	for i, w := range widths {
		sep += strings.Repeat("─", w)
		sep += "┴"
	}
	sep = sep[:len(sep)-1] + "╯\n"
	b.WriteString(sep)

	return b.String()
}

func (t *Table) Print() {
	fmt.Print(t.Render())
}
