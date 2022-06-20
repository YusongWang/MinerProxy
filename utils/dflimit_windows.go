//go:build windows

package utils

func IncreaseFDLimit() {
	// Windows has no file descriptor limit. Do nothing.
}
