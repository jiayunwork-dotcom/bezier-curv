package cli

import (
	"fmt"
	"io"
	"strings"
)

type table struct {
	headers []string
	rows    [][]string
	widths  []int
}

func newTable(headers ...string) *table {
	t := &table{headers: headers, widths: make([]int, len(headers))}
	for i, h := range headers {
		t.widths[i] = len(h)
	}
	return t
}

func (t *table) add(row ...string) {
	if len(row) != len(t.headers) {
		panic("table row width mismatch")
	}
	cp := make([]string, len(row))
	for i, c := range row {
		cp[i] = c
		if len(c) > t.widths[i] {
			t.widths[i] = len(c)
		}
	}
	t.rows = append(t.rows, cp)
}

func (t *table) render(w io.Writer) {
	var b strings.Builder
	sep := make([]string, len(t.widths))
	for i, wd := range t.widths {
		sep[i] = strings.Repeat("-", wd)
	}
	b.WriteString(" " + strings.Join(sep, "   ") + "\n")
	b.WriteString(" " + strings.Join(padAll(t.headers, t.widths), "   ") + "\n")
	b.WriteString(" " + strings.Join(sep, "   ") + "\n")
	for _, r := range t.rows {
		b.WriteString(" " + strings.Join(padAll(r, t.widths), "   ") + "\n")
	}
	b.WriteString(" " + strings.Join(sep, "   ") + "\n")
	io.WriteString(w, b.String())
}

func padAll(cells []string, widths []int) []string {
	out := make([]string, len(cells))
	for i, c := range cells {
		out[i] = fmt.Sprintf("%-*s", widths[i], c)
	}
	return out
}
