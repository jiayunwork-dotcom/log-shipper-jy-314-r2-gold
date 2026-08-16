// Package report 渲染检索命中的可读结果。
package report

import (
	"fmt"
	"io"

	"log-shipper/internal/index"
)

// Render 将命中列表输出为可读文本。
func Render(hits []index.Hit, w io.Writer) {
	if len(hits) == 0 {
		fmt.Fprintln(w, "no matches")
		return
	}
	for _, h := range hits {
		fmt.Fprintf(w, "%s:%d: %s\n", h.File, h.Line, h.Text)
	}
}
