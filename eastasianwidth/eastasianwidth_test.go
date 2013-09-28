package eastasianwidth

import (
	"testing"
	"unicode"

	"github.com/hnakamur/unicodeutil/eastasianwidth/unicodetxttogo/lib"
)

func TestIsFullwidthSamles(t *testing.T) {
	if IsFullwidth('a') {
		t.Errorf("'a' must be halfwidth")
	}
	if !IsFullwidth('あ') {
		t.Errorf("'あ' must be fullwidth")
	}
}

func TestIsFullwidthAll(t *testing.T) {
	m := buildFullwidthMap()
	for r := '\u0000'; r <= unicode.MaxRune; r++ {
		wanted := m[r]
		got := IsFullwidth(r)
		if got != wanted {
			t.Errorf("TestIsFullwidthAll: got=%v, wanted=%v, rune=%x",
				got, wanted, r)
		}
	}
}

func buildFullwidthMap() map[rune]bool {
	m := make(map[rune]bool)
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
			m[lo] = true
		} else {
			for r := lo; r <= hi; r++ {
				m[r] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return m
}
