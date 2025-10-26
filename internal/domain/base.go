package domain

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"unicode"
)

// removeAccents removes accents from a given string.
func RemoveAccents(s string) (string, error) {
	// transformation chain: decompose, remove non-spacing marks, recompose
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	// apply the transformation
	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}
	return result, nil
}
