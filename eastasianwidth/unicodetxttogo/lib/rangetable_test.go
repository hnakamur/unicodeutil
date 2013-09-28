package lib

import (
	"fmt"
	"testing"
	"unicode"
)

func TestAddRune(t *testing.T) {
	tests := []struct {
		name string
		fn   func() (t *unicode.RangeTable, err error)
		want *unicode.RangeTable
		err  error
	}{
		{
			"fromEmpty",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{}
				err = AppendRune(t, '\u0000')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0000, 1},
				},
				LatinOffset: 1,
			},
			nil,
		},
		{
			"adjustStride",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0000, 1},
					},
					LatinOffset: 1,
				}
				err = AppendRune(t, '\u0002')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0002, 2},
				},
				LatinOffset: 1,
			},
			nil,
		},
		{
			"extend",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0001, 1},
					},
					LatinOffset: 1,
				}
				err = AppendRune(t, '\u0002')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0002, 1},
				},
				LatinOffset: 1,
			},
			nil,
		},
		{
			"append",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0001, 1},
					},
					LatinOffset: 1,
				}
				err = AppendRune(t, '\u0003')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0001, 1},
					{0x0003, 0x0003, 1},
				},
				LatinOffset: 2,
			},
			nil,
		},
		{
			"error1",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0000, 1},
					},
					LatinOffset: 1,
				}
				err = AppendRune(t, '\u0000')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0000, 1},
				},
				LatinOffset: 1,
			},
			fmt.Errorf("AppendRange: lo must be greater than Hi of the last range: lo=%d", 0),
		},
	}
	for _, test := range tests {
		got, err := test.fn()
		if test.err != nil {
			compareError(t, test.name, err, test.err)
		} else {
			compareRangeTables(t, test.name, got, test.want)
		}
	}
}

func TestAddRange(t *testing.T) {
	tests := []struct {
		name string
		fn   func() (t *unicode.RangeTable, err error)
		want *unicode.RangeTable
		err  error
	}{
		{
			"fromEmpty",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{}
				err = AppendRange(t, '\u0000', '\u0001')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0001, 1},
				},
				LatinOffset: 1,
			},
			nil,
		},
		{
			"adjustStride",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0002, 2},
					},
					LatinOffset: 1,
				}
				err = AppendRange(t, '\u0003', '\u0004')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0000, 1},
					{0x0002, 0x0004, 1},
				},
				LatinOffset: 2,
			},
			nil,
		},
		{
			"adjustHi",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0004, 2},
					},
					LatinOffset: 1,
				}
				err = AppendRange(t, '\u0005', '\u0006')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0002, 2},
					{0x0004, 0x0006, 1},
				},
				LatinOffset: 2,
			},
			nil,
		},
		{
			"extend",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0001, 1},
					},
					LatinOffset: 1,
				}
				err = AppendRange(t, '\u0002', '\u0003')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0003, 1},
				},
				LatinOffset: 1,
			},
			nil,
		},
		{
			"append",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0001, 1},
					},
					LatinOffset: 1,
				}
				err = AppendRange(t, '\u0003', '\u0004')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0001, 1},
					{0x0003, 0x0004, 1},
				},
				LatinOffset: 2,
			},
			nil,
		},
		{
			"span16To32",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0xfffd, 0xfffe, 1},
					},
					LatinOffset: 0,
				}
				err = AppendRange(t, '\uffff', '\U00010000')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0xfffd, 0xffff, 1},
				},
				R32: []unicode.Range32{
					{0x10000, 0x10000, 1},
				},
				LatinOffset: 0,
			},
			nil,
		},
		{
			"error1",
			func() (t *unicode.RangeTable, err error) {
				t = &unicode.RangeTable{
					R16: []unicode.Range16{
						{0x0000, 0x0000, 1},
					},
					LatinOffset: 1,
				}
				err = AppendRange(t, '\u0000', '\u0001')
				return
			},
			&unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0000, 0x0000, 1},
				},
				LatinOffset: 1,
			},
			fmt.Errorf("AppendRange: lo must be greater than Hi of the last range: lo=%d", 0),
		},
	}
	for _, test := range tests {
		got, err := test.fn()
		if test.err != nil {
			compareError(t, test.name, err, test.err)
		} else {
			compareRangeTables(t, test.name, got, test.want)
		}
	}
}

func compareError(t *testing.T, caseName string, got, want error) {
	if got == nil || got.Error() != want.Error() {
		t.Errorf("%s\ngot:%s\nwant:%s\n", caseName, got, want)
	}
}

func compareRangeTables(t *testing.T, caseName string, got, want *unicode.RangeTable) {
	same := true
	defer func() {
		if !same {
			t.Errorf("%s\ngot:%s\nwant:%s\n", caseName, String(got), String(want))
		}
	}()
	if len(got.R16) != len(want.R16) {
		same = false
		return
	} else {
		for i := 0; i < len(got.R16); i++ {
			gotR := got.R16[i]
			wantR := want.R16[i]
			if gotR.Lo != wantR.Lo || gotR.Hi != wantR.Hi || gotR.Stride != wantR.Stride {
				same = false
				return
			}
		}
	}
	if len(got.R32) != len(want.R32) {
		same = false
		return
	} else {
		for i := 0; i < len(got.R32); i++ {
			gotR := got.R32[i]
			wantR := want.R32[i]
			if gotR.Lo != wantR.Lo || gotR.Hi != wantR.Hi || gotR.Stride != wantR.Stride {
				same = false
				return
			}
		}
	}
	same = got.LatinOffset == want.LatinOffset
	return
}
