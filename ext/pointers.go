package ext

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func BoolPtr(v bool) *bool {
	return &v
}
