package utils

func StringLengthBetween(str string, min int, max int) bool {
	return len(str) >= min && len(str) <= max
}
