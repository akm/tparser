package runes

import "unicode"

func IsWord(r rune) bool {
	return IsWordHead(r) || IsDigit(r)
}

func IsWordHead(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}
