package util

// NewFalse is a silly hack to get a false bool pointer
func NewFalse() *bool {
	b := false
	return &b
}
