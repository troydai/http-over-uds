package tablify

import (
	"fmt"
	"io"
	"strings"

	"github.com/troydai/http-over-uds/internal/summary"
)

const _columns = ",Count,Success,Error,p99,p95,p50,Status"

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

func Print(w io.Writer, data []*summary.Series, digest, hasTotal bool) {
	lines := GetLines(data)
	header := lines[0]
	lines = lines[1:]

	var digestEffect bool

	boundary := len(data)
	if hasTotal {
		boundary--
	}
	if digest && boundary > 10 {
		boundary = 10
		digestEffect = true
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, header)

	for i := 0; i < boundary; i++ {
		fmt.Fprintln(w, lines[i])
	}

	if digestEffect {
		fmt.Fprintln(w, "... omit other samples beyond the first 10 ...")
	}

	if hasTotal {
		fmt.Fprintln(w, lines[len(lines)-1])
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
