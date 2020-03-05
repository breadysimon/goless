package text

func StartWith(s, sub string) bool {
	return len(s) > len(sub) && s[:len(sub)] == sub
}
