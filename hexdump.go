package pretty

import (
	"fmt"
	"io"
)

func HexDump(out io.Writer, by []byte, number int, indent string) {
	n := len(by)
	rowCount := 0
	stop := (n / number) * number
	k := 0
	for i := 0; i <= stop; i += number {
		k++
		if i+number < n {
			rowCount = number
		} else {
			rowCount = min(k*number, n) % number
		}

		_, _ = fmt.Fprintf(out, indent+"%04d ", i)
		for j := 0; j < rowCount; j++ {
			_, _ = fmt.Fprintf(out, "%02x  ", by[i+j])
		}
		for j := rowCount; j < number; j++ {
			_, _ = fmt.Fprintf(out, "    ")
		}
		_, _ = fmt.Fprintf(out, "  '%s'\n", ViewString(by[i:(i+rowCount)]))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ViewString(b []byte) string {
	r := []rune(string(b))
	for i := range r {
		if r[i] < 32 || r[i] > 126 {
			r[i] = '.'
		}
	}
	return string(r)
}
