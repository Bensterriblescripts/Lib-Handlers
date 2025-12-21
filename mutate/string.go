package mutate

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Capitalise a UTF-8 String
func Capitalise(s string) string {
	if s == "" {
		return s
	}

	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size == 1 { // Invalid UTF-8
		return s
	}

	ur := unicode.ToUpper(r)
	if ur == r {
		return s
	}

	var b strings.Builder
	b.Grow(len(s) + 3)
	b.WriteRune(ur)
	b.WriteString(s[size:])
	return b.String()
}
