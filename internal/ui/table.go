package ui

import (
	"github.com/charmbracelet/lipgloss"
	"log"
)
import "github.com/charmbracelet/lipgloss/table"

const (
	gray        = lipgloss.Color("245")
	lightGray   = lipgloss.Color("241")
	lighterGray = lipgloss.Color("253")
)

type Table struct {
	Columns      []string
	Rows         [][]string
	ColumnWidths map[string]int
	TotalWidth   int
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) SetTotalWidth(width int) *Table {
	t.TotalWidth = width

	return t
}

func (t *Table) SetColumns(columns []string) *Table {
	t.Columns = columns

	return t
}

func (t *Table) AddRow(row []string) *Table {
	t.Rows = append(t.Rows, row)

	return t
}

func (t *Table) SetColumnWidth(column string, width int) *Table {
	if t.ColumnWidths == nil {
		t.ColumnWidths = make(map[string]int)
	}
	t.ColumnWidths[column] = width

	return t
}

func (t *Table) Render() string {
	var (
		// HeaderStyle is the lipgloss style used for the table headers.
		HeaderStyle = lipgloss.NewStyle().Foreground(gray).Bold(true).Align(lipgloss.Left)
		// CellStyle is the base lipgloss style used for the table rows.
		CellStyle = lipgloss.NewStyle()
		// OddRowStyle is the lipgloss style used for odd-numbered table rows.
		OddRowStyle = CellStyle.Foreground(gray)
		// EvenRowStyle is the lipgloss style used for even-numbered table rows.
		EvenRowStyle = CellStyle.Foreground(lightGray)
		// BorderStyle is the lipgloss style used for the table border.
		BorderStyle = lipgloss.NewStyle().Foreground(lighterGray)
	)

	instance := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				style = EvenRowStyle
			default:
				style = OddRowStyle
			}

			style = style.Width(t.getColumnWidth(t.Columns[col]))

			return style
		}).
		Headers(t.Columns...).
		Rows(t.Rows...)

	return instance.Render()
}

func (t *Table) getColumnWidth(column string) int {
	// if the width is set, return it
	if width, ok := t.ColumnWidths[column]; ok {
		log.Printf("width: %d", width)
		return width
	}

	// calculate the remaining width by subtracting the total of the specified column widths
	remainingColumns := 0
	totalSpecifiedWidth := 0

	for _, col := range t.Columns {
		if width, ok := t.ColumnWidths[col]; ok {
			totalSpecifiedWidth += width
		} else {
			remainingColumns++
		}
	}

	// if no columns without specified widths, return a default width
	if remainingColumns == 0 {
		return 10 // just a default fallback
	}

	// distribute the remaining width evenly among columns without specified widths
	remainingWidth := t.TotalWidth - totalSpecifiedWidth
	return remainingWidth / remainingColumns
}
