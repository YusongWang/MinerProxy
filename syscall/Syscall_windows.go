//go:build windows

package syscall

func IncreaseFDLimit() {
	// Windows has no file descriptor limit. Do nothing.
}
