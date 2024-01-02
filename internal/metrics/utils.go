package metrics

import "github.com/miscord-dev/epgstation-exporter/api/epgstation"

// DerefIntPtr returns the value of the given pointer to int, or defaultVal if nil.
func derefIntPtr(p *int, defaultVal int) int {
	if p == nil {
		return defaultVal
	}
	return *p
}

// IntPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}

// getReserveTypePtr returns a pointer to the given epgstation.GetReserveType.
func getReserveTypePtr(rt epgstation.GetReserveType) *epgstation.GetReserveType {
	return &rt
}
