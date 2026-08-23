package cli

import (
	"fmt"
	"io"
	"strings"
)

// table 是简单的文本表格渲染器。
type table struct {
	headers []string
	rows    [][]string
	widths  []int
}

// newTable 创建带表头的表格。
func newTable(headers ...string) *table {
	t := &table{headers: headers, widths: make([]int, len(headers))}
	for i, h := range headers {
		t.widths[i] = len(h)
	}
	return t
}

// add 追加一行；单元格数量须与表头一致。
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

// render 输出表格（无对齐填充）。
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

// padAll 把每列左对齐填充到固定宽度。
func padAll(cells []string, widths []int) []string {
	out := make([]string, len(cells))
	for i, c := range cells {
		out[i] = fmt.Sprintf("%-*s", widths[i], c)
	}
	return out
}
