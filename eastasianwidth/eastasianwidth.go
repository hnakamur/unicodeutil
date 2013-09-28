package eastasianwidth

import (
	"unicode"
)

// Returns whether r is fullwidth or not.
// It returns true if the category in East Asian Width is in one of
// "W", "F" or "A".
// See:
// - UAX #11: East Asian Width http://www.unicode.org/reports/tr11/
// - 東アジアの文字幅 - Wikipedia http://ja.wikipedia.org/wiki/%E6%9D%B1%E3%82%A2%E3%82%B8%E3%82%A2%E3%81%AE%E6%96%87%E5%AD%97%E5%B9%85
func IsFullwidth(r rune) bool {
	return unicode.Is(Fullwidth, r)
}
