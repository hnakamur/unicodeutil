package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func IsFullwidthCategory(category string) bool {
	return category == "W" || category == "F" || category == "A"
}

// Returns a new scanner to read from EastAsianWidthText
func NewEastAsianWidthTextScanner() *bufio.Scanner {
	return bufio.NewScanner(bytes.NewReader([]byte(EastAsianWidthText)))
}

// Returns the next entry, lo: range start rune, hi: range end rune,
// category: category of the entry ("N", "A", "H", "W", "F", or "Na").
// Returns the empty string "" as category when no more entries exist.
func NextEntry(scanner *bufio.Scanner) (lo rune, hi rune, category string) {
	line := nextDataLine(scanner)
	if line == "" {
		return
	}

	i := strings.Index(line, "..")
	j := strings.LastIndex(line, ";")
	if i == -1 {
		lo = readHex(line[:j])
		hi = 0
		category = line[j+1:]
	} else {
		lo = readHex(line[:i])
		hi = readHex(line[i+2 : j])
		category = line[j+1:]
	}
	return
}

func readHex(s string) (r rune) {
	n, err := fmt.Sscanf(s, "%X", &r)
	if err != nil {
		panic(err)
	} else if n != 1 {
		panic(fmt.Errorf("expected hex string, but got %s", s))
	}
	return
}

// Returns the next data line or an empty string when eof
func nextDataLine(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		line := trimComment(scanner.Text())
		if line != "" {
			return line
		}
	}
	return ""
}

func trimComment(line string) string {
	index := strings.Index(line, "#")
	if index != -1 {
		line = line[0:index]
	}
	return strings.TrimSpace(line)
}
