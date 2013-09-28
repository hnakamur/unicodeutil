package lib

import (
	"bytes"
	"fmt"
	"unicode"
)

const maxR16 = '\uFFFF'

// Add a rune to table. r must be greater than Hi of the last range.
func AppendRune(t *unicode.RangeTable, r rune) error {
	return AppendRange(t, r, r)
}

// Add a range to table. lo must be greater than Hi of the last range.
func AppendRange(t *unicode.RangeTable, lo, hi rune) error {
	if hi < lo {
		return fmt.Errorf("AppendRange: hi < lo. lo=%x, hi=%x", lo, hi)
	}

	if hi <= maxR16 {
		if err := appendRange16(t, uint16(lo), uint16(hi)); err != nil {
			return err
		}
	} else if hi <= unicode.MaxRune {
		if lo <= maxR16 {
			if err := appendRange16(t, uint16(lo), maxR16); err != nil {
				return err
			}
			if err := appendRange32(t, maxR16+1, uint32(hi)); err != nil {
				return err
			}
		} else {
			if err := appendRange32(t, uint32(lo), uint32(hi)); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("AppendRange: MaxRune < hi. lo=%x, hi=%x", lo, hi)
	}
	return nil
}

func appendRange16(t *unicode.RangeTable, lo, hi uint16) error {
	if t.R16 == nil || len(t.R16) == 0 {
		t.R16 = append(t.R16, unicode.Range16{lo, hi, 1})
		if hi <= unicode.MaxLatin1 {
			t.LatinOffset++
		}
		return nil
	}

	range_ := &t.R16[len(t.R16)-1]
	if lo <= range_.Hi {
		return fmt.Errorf("AppendRange: lo must be greater than Hi of the last range: lo=%d", lo)
	}

	if lo == hi {
		if range_.Hi+range_.Stride == lo {
			range_.Hi = lo
			if len(t.R16) >= 2 {
				prevR := &t.R16[len(t.R16)-2]
				if prevR.Stride > 1 && prevR.Hi+1 == range_.Lo {
					prevCount := (prevR.Hi-prevR.Lo)/prevR.Stride + 1
					count := (range_.Hi-range_.Lo)/range_.Stride + 1
					if prevCount <= count {
						range_.Lo--
						prevR.Hi -= prevR.Stride
						if prevR.Lo == prevR.Hi {
							prevR.Stride = 1
						}
					}
				}
			}
		} else if range_.Lo == range_.Hi {
			range_.Hi = lo
			range_.Stride = lo - range_.Lo
		} else {
			t.R16 = append(t.R16, unicode.Range16{lo, hi, 1})
			if hi <= unicode.MaxLatin1 {
				t.LatinOffset++
			}
		}
	} else {
		if range_.Stride == 1 {
			if lo-1 <= range_.Hi {
				range_.Hi = hi
			} else {
				if range_.Lo == range_.Hi {
					range_.Hi = hi
					range_.Stride = range_.Hi - range_.Lo
				} else {
					t.R16 = append(t.R16, unicode.Range16{lo, hi, 1})
					if hi <= unicode.MaxLatin1 {
						t.LatinOffset++
					}
				}
			}
		} else {
			if lo-1 <= range_.Hi {
				range_.Hi -= range_.Stride
				if range_.Lo == range_.Hi {
					range_.Stride = 1
				}
				t.R16 = append(t.R16, unicode.Range16{lo - 1, hi, 1})
				if hi <= unicode.MaxLatin1 {
					t.LatinOffset++
				}
			} else {
				t.R16 = append(t.R16, unicode.Range16{lo, hi, 1})
				if hi <= unicode.MaxLatin1 {
					t.LatinOffset++
				}
			}
		}
	}
	return nil
}

func appendRange32(t *unicode.RangeTable, lo, hi uint32) error {
	if t.R32 == nil || len(t.R32) == 0 {
		t.R32 = append(t.R32, unicode.Range32{lo, hi, 1})
		return nil
	}

	range_ := &t.R32[len(t.R32)-1]
	if lo <= range_.Hi {
		return fmt.Errorf("AppendRange: lo must be greater than Hi of the last range: lo=%d", lo)
	}

	if lo == hi {
		if range_.Hi+range_.Stride == lo {
			range_.Hi = lo
			if len(t.R32) >= 2 {
				prevR := &t.R32[len(t.R32)-2]
				if prevR.Stride > 1 && prevR.Hi+1 == range_.Lo {
					prevCount := (prevR.Hi-prevR.Lo)/prevR.Stride + 1
					count := (range_.Hi-range_.Lo)/range_.Stride + 1
					if prevCount <= count {
						range_.Lo--
						prevR.Hi -= prevR.Stride
						if prevR.Lo == prevR.Hi {
							prevR.Stride = 1
						}
					}
				}
			}
		} else if range_.Lo == range_.Hi {
			range_.Hi = lo
			range_.Stride = lo - range_.Lo
		} else {
			t.R32 = append(t.R32, unicode.Range32{lo, hi, 1})
		}
	} else {
		if range_.Stride == 1 {
			if lo-1 <= range_.Hi {
				range_.Hi = hi
			} else {
				t.R32 = append(t.R32, unicode.Range32{lo, hi, 1})
			}
		} else {
			if lo-1 <= range_.Hi {
				range_.Hi -= range_.Stride
				if range_.Lo == range_.Hi {
					range_.Stride = 1
				}
				t.R32 = append(t.R32, unicode.Range32{lo - 1, hi, 1})
			} else {
				t.R32 = append(t.R32, unicode.Range32{lo, hi, 1})
			}
		}
	}
	return nil
}

func String(t *unicode.RangeTable) string {
	var buf bytes.Buffer
	buf.WriteString("unicode.RangeTable{\n")
	buf.WriteString("\tR16: []unicode.Range16{\n")
	if t.R16 != nil {
		for _, r := range t.R16 {
			buf.WriteString(fmt.Sprintf("\t\t{0x%04x, 0x%04x, %d},\n",
				r.Lo, r.Hi, r.Stride))
		}
	}
	buf.WriteString("\t},\n")
	buf.WriteString("\tR32: []unicode.Range32{\n")
	if t.R32 != nil {
		for _, r := range t.R32 {
			buf.WriteString(fmt.Sprintf("\t\t{0x%x, 0x%x, %d},\n",
				r.Lo, r.Hi, r.Stride))
		}
	}
	buf.WriteString("\t},\n")
	buf.WriteString(fmt.Sprintf("\tLatinOffset: %d,\n", t.LatinOffset))
	buf.WriteString("}\n")
	return buf.String()
}
