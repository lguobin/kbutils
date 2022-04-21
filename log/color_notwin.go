//go:build !windows
// +build !windows

package log

// IsSupportColor IsSupportColor
func IsSupportColor() bool {
	return supportColor
}
