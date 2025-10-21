package progress

import (
	"os"
)

// IsTerminal checks if the given file descriptor is a terminal
// This is used to determine if we should show progress indicators
func IsTerminal(fd uintptr) bool {
	// Check if the file is a character device (terminal)
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	// On Unix systems, terminals are character devices
	// ModeCharDevice is set for character devices like terminals
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// IsStderrTerminal checks if stderr is connected to a terminal
func IsStderrTerminal() bool {
	fileInfo, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// IsStdoutTerminal checks if stdout is connected to a terminal
func IsStdoutTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
