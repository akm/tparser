package runes

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsNumericHead(r rune) bool {
	return IsDigit(r) || r == '-' // || r == '+' // Don't allow + for now
}

func IsUnaryOp(r rune) bool {
	return r == '-'
}
