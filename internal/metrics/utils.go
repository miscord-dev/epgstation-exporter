package metrics

// DerefIntPtr returns the value of the given pointer to int, or defaultVal if nil.
func derefIntPtr(p *int, defaultVal int) int {
	if p == nil {
		return defaultVal
	}
	return *p
}
