package logger

import (
	"fmt"
	"os"
)

// Color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// Info prints an info message
func Info(format string, args ...interface{}) {
	fmt.Printf(colorBlue+"ℹ "+colorReset+format+"\n", args...)
}

// Success prints a success message
func Success(format string, args ...interface{}) {
	fmt.Printf(colorGreen+"✓ "+colorReset+format+"\n", args...)
}

// Warning prints a warning message
func Warning(format string, args ...interface{}) {
	fmt.Printf(colorYellow+"⚠ "+colorReset+format+"\n", args...)
}

// Error prints an error message
func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorRed+"✗ "+colorReset+format+"\n", args...)
}

// Debug prints a debug message
func Debug(format string, args ...interface{}) {
	fmt.Printf(colorCyan+"[DEBUG] "+colorReset+format+"\n", args...)
}

// Plain prints a plain message without formatting
func Plain(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// PrintSummary prints a summary of the migration results
func PrintSummary(created, updated, skipped, errors int) {
	Plain("\n" + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	Plain("Migration Summary")
	Plain("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if created > 0 {
		Success("Created: %d", created)
	}
	if updated > 0 {
		Success("Updated: %d", updated)
	}
	if skipped > 0 {
		Warning("Skipped: %d", skipped)
	}
	if errors > 0 {
		Error("Errors: %d", errors)
	}

	total := created + updated + skipped
	Plain("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	Plain("Total processed: %d", total)
}
