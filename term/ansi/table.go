package ansi

import (
	"strings"
)

type Table struct {
	rows [][]string
	maxW []int

	Min int
	Pad int
}

func (t *Table) Reset() {
	t.rows = t.rows[:0]
	t.maxW = t.maxW[:0]
}

func (t *Table) AddLine(s string) {
	t.rows = append(t.rows, []string{s})
}

func (t *Table) AddRow(cols ...string) {
	t.rows = append(t.rows, cols)
	if len(t.maxW) == 0 {
		t.maxW = make([]int, len(cols))
	}

	for i, col := range cols {
		if len(col)+t.Pad > t.maxW[i] {
			t.maxW[i] = len(col) + t.Pad
		}
		if t.maxW[i] < t.Min+t.Pad {
			t.maxW[i] = t.Min + t.Pad
		}
	}
}

func (t *Table) String() string {
	var b strings.Builder

	for _, row := range t.rows {
		if len(row) == 1 {
			b.WriteString(row[0] + "\r\n")
			continue
		}
		for i, col := range row {
			if i > 0 {
				b.WriteRune(' ')
			}
			b.WriteString(col)
			b.WriteString(strings.Repeat(" ", t.maxW[i]-len(col)))
		}
		b.WriteString("\r\n")
	}

	return b.String()
}
