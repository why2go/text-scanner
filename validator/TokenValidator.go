package validator

import (
	"fmt"
	"unicode"

	scanner "gitee.com/piecat/text-scanner"
)

type Validator struct{}

func (v *Validator) IsValidToken(rtext []rune, match scanner.Match) bool {
	lcp := rtext[match.S]
	rcp := rtext[match.E-1]
	needBeToken := !v.isEastAsianCodePoint(lcp) && !v.isEastAsianCodePoint(rcp)
	if !needBeToken {
		return true
	}
	return v.isToken(rtext, match)
}

func (v *Validator) isEastAsianCodePoint(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana)
}

func (v *Validator) isToken(rtext []rune, match scanner.Match) bool {
	fmt.Printf("rtext[match.S:match.E]: %s\n", string(rtext[match.S:match.E]))
	l := match.S
	r := match.E - 1
	var lmatched bool
	var rmatched bool
	if l <= 0 || v.notTokenCodePoint(rtext[l-1]) {
		lmatched = true
	}
	if r >= len(rtext)-1 || v.notTokenCodePoint(rtext[r+1]) {
		rmatched = true
	}
	return lmatched && rmatched
}

func (v *Validator) notTokenCodePoint(r rune) bool {
	return unicode.IsSpace(r) ||
		unicode.IsPunct(r) ||
		// unicode.IsControl(r) ||
		unicode.IsSymbol(r) ||
		unicode.IsMark(r) ||
		v.isEastAsianCodePoint(r)
}
