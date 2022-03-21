package ext

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
