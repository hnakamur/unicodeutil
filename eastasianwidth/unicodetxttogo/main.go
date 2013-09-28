package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"github.com/hnakamur/unicodeutil/eastasianwidth/unicodetxttogo/lib"
)

// This program generates eastasianwidth/table.go.
func main() {
	fullwidthTable := &unicode.RangeTable{}
	scanner := lib.NewEastAsianWidthTextScanner()
	for {
		lo, hi, category := lib.NextEntry(scanner)
		if category == "" {
			break
		}
		if !lib.IsFullwidthCategory(category) {
			continue
		}

		if hi == 0 {
			lib.AppendRune(fullwidthTable, lo)
		} else {
			lib.AppendRange(fullwidthTable, lo, hi)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	file, err := os.Create("../table.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	out := bufio.NewWriter(file)
	fmt.Fprintf(out, format, "Fullwidth",
		lib.String(fullwidthTable))
	out.Flush()
}

const format = `// Generated by running
//   ./unicodetxttogo
// DO NOT EDIT

package eastasianwidth

import (
	"unicode"
)

var %s = &%s`