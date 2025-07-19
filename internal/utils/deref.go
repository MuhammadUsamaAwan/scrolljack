package utils

func DerefStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

func DerefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func DerefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
