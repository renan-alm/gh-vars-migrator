package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// captureOutput captures stdout/stderr output for testing
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

// TestInfo tests the Info logging function
func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		Info("test message")
	})

	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "ℹ") {
		t.Errorf("Expected output to contain info icon, got: %s", output)
	}
}

// TestSuccess tests the Success logging function
func TestSuccess(t *testing.T) {
	output := captureOutput(func() {
		Success("success message")
	})

	if !strings.Contains(output, "success message") {
		t.Errorf("Expected output to contain 'success message', got: %s", output)
	}
	if !strings.Contains(output, "✓") {
		t.Errorf("Expected output to contain success icon, got: %s", output)
	}
}

// TestWarning tests the Warning logging function
func TestWarning(t *testing.T) {
	output := captureOutput(func() {
		Warning("warning message")
	})

	if !strings.Contains(output, "warning message") {
		t.Errorf("Expected output to contain 'warning message', got: %s", output)
	}
	if !strings.Contains(output, "⚠") {
		t.Errorf("Expected output to contain warning icon, got: %s", output)
	}
}

// TestError tests the Error logging function
func TestError(t *testing.T) {
	// Capture stderr instead of stdout
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Error("error message")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Errorf("Expected output to contain 'error message', got: %s", output)
	}
	if !strings.Contains(output, "✗") {
		t.Errorf("Expected output to contain error icon, got: %s", output)
	}
}

// TestDebug tests the Debug logging function
func TestDebug(t *testing.T) {
	output := captureOutput(func() {
		Debug("debug message")
	})

	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected output to contain 'debug message', got: %s", output)
	}
	if !strings.Contains(output, "[DEBUG]") {
		t.Errorf("Expected output to contain [DEBUG], got: %s", output)
	}
}

// TestPlain tests the Plain logging function
func TestPlain(t *testing.T) {
	output := captureOutput(func() {
		Plain("plain message")
	})

	if !strings.Contains(output, "plain message") {
		t.Errorf("Expected output to contain 'plain message', got: %s", output)
	}
	// Plain should not have icons or special formatting
	if strings.Contains(output, "✓") || strings.Contains(output, "✗") || strings.Contains(output, "ℹ") {
		t.Errorf("Plain output should not contain icons, got: %s", output)
	}
}

// TestPrintSummary tests the PrintSummary function
func TestPrintSummary(t *testing.T) {
	tests := []struct {
		name    string
		created int
		updated int
		skipped int
		errors  int
	}{
		{
			name:    "all counts",
			created: 5,
			updated: 3,
			skipped: 2,
			errors:  1,
		},
		{
			name:    "no errors",
			created: 10,
			updated: 0,
			skipped: 0,
			errors:  0,
		},
		{
			name:    "only errors",
			created: 0,
			updated: 0,
			skipped: 0,
			errors:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture both stdout and stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			PrintSummary(tt.created, tt.updated, tt.skipped, tt.errors)

			wOut.Close()
			wErr.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			var bufOut bytes.Buffer
			_, _ = bufOut.ReadFrom(rOut)
			var bufErr bytes.Buffer
			_, _ = bufErr.ReadFrom(rErr)

			// Combine stdout and stderr
			output := bufOut.String() + bufErr.String()

			if !strings.Contains(output, "Migration Summary") {
				t.Errorf("Expected output to contain 'Migration Summary', got: %s", output)
			}

			if tt.created > 0 && !strings.Contains(output, "Created:") {
				t.Errorf("Expected output to contain 'Created:', got: %s", output)
			}

			if tt.updated > 0 && !strings.Contains(output, "Updated:") {
				t.Errorf("Expected output to contain 'Updated:', got: %s", output)
			}

			if tt.skipped > 0 && !strings.Contains(output, "Skipped:") {
				t.Errorf("Expected output to contain 'Skipped:', got: %s", output)
			}

			if tt.errors > 0 && !strings.Contains(output, "Errors:") {
				t.Errorf("Expected output to contain 'Errors:', got: %s", output)
			}

			if !strings.Contains(output, "Total processed:") {
				t.Errorf("Expected output to contain 'Total processed:', got: %s", output)
			}
		})
	}
}

// TestFormattingWithArguments tests that formatting with arguments works
func TestFormattingWithArguments(t *testing.T) {
	output := captureOutput(func() {
		Info("Test with %s and %d", "string", 42)
	})

	if !strings.Contains(output, "Test with string and 42") {
		t.Errorf("Expected formatted output, got: %s", output)
	}
}
