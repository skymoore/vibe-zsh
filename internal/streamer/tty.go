package streamer

import (
	"os"
)

// IsStdoutTerminal checks if stdout is connected to a terminal
// This is used to determine if we should stream output or print instantly
func IsStdoutTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// ShouldStream determines if output should be streamed based on TTY detection
// Returns false if stdout is not a terminal (e.g., piped to file or another command)
func ShouldStream() bool {
	return IsStdoutTerminal()
}
