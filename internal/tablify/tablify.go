package tablify

import (
	"fmt"
	"io"
	"strings"

	"github.com/troydai/http-over-uds/internal/summary"
)

const _columns = ",Count,Error,p99,p95,p50,Status"

func GetLines(data []*summary.Series) []string {
	var content [][]string
	content = append(content, strings.Split(_columns, ","))
	for _, s := range data {
		content = append(content, s.PresentData())
	}

	cow := make([]int, len(content[0])) // column widths
	for _, row := range content {
		for idx, co := range row {
			if len(co) > cow[idx] {
				cow[idx] = len(co)
			}
		}
	}

	var template string
	for _, w := range cow {
		template += fmt.Sprintf("%%%ds  ", w)
	}

	var retval []string
	for _, s := range content {
		line := fmt.Sprintf(template, toAny(s)...)
		retval = append(retval, line)
	}

	return retval
}

func Print(w io.Writer, data []*summary.Series) {
	fmt.Fprintln(w)
	for _, l := range GetLines(data) {
		fmt.Fprintln(w, l)
	}
	fmt.Fprintln(w)
}

func toAny(s []string) []interface{} {
	any := make([]interface{}, len(s))
	for idx := range s {
		any[idx] = s[idx]
	}

	return any
}
