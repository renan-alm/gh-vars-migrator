package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_MissingFile(t *testing.T) {
	err := Load("nonexistent_file_xyz.env")
	if err != nil {
		t.Fatalf("expected nil for missing file, got: %v", err)
	}
}

func TestLoad_SetsNewVars(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "# comment line\nSOURCE_ORG=my-org\nTARGET_ORG=other-org\nSOURCE_HOSTNAME=github.mycompany.com\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// t.Setenv automatically restores the previous value on cleanup and
	// satisfies errcheck because the restore is handled internally.
	for _, key := range []string{"SOURCE_ORG", "TARGET_ORG", "SOURCE_HOSTNAME"} {
		t.Setenv(key, "")
		_ = os.Unsetenv(key)
	}

	if err := Load(envPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("SOURCE_ORG"); got != "my-org" {
		t.Errorf("SOURCE_ORG = %q, want %q", got, "my-org")
	}
	if got := os.Getenv("TARGET_ORG"); got != "other-org" {
		t.Errorf("TARGET_ORG = %q, want %q", got, "other-org")
	}
	if got := os.Getenv("SOURCE_HOSTNAME"); got != "github.mycompany.com" {
		t.Errorf("SOURCE_HOSTNAME = %q, want %q", got, "github.mycompany.com")
	}
}

func TestLoad_DoesNotOverrideExistingVars(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "MY_TEST_VAR=from-file\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("MY_TEST_VAR", "from-env")

	if err := Load(envPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("MY_TEST_VAR"); got != "from-env" {
		t.Errorf("MY_TEST_VAR = %q, want %q (should not be overridden)", got, "from-env")
	}
}

func TestLoad_QuotedValues(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "DOUBLE_QUOTED=\"hello world\"\nSINGLE_QUOTED='hello world'\nUNQUOTED=hello\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"DOUBLE_QUOTED", "SINGLE_QUOTED", "UNQUOTED"} {
		t.Setenv(key, "")
		_ = os.Unsetenv(key)
	}

	if err := Load(envPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("DOUBLE_QUOTED"); got != "hello world" {
		t.Errorf("DOUBLE_QUOTED = %q, want %q", got, "hello world")
	}
	if got := os.Getenv("SINGLE_QUOTED"); got != "hello world" {
		t.Errorf("SINGLE_QUOTED = %q, want %q", got, "hello world")
	}
	if got := os.Getenv("UNQUOTED"); got != "hello" {
		t.Errorf("UNQUOTED = %q, want %q", got, "hello")
	}
}

func TestLoad_ExportPrefix(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "export MY_EXPORT_VAR=exported-value\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("MY_EXPORT_VAR", "")
	_ = os.Unsetenv("MY_EXPORT_VAR")

	if err := Load(envPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("MY_EXPORT_VAR"); got != "exported-value" {
		t.Errorf("MY_EXPORT_VAR = %q, want %q", got, "exported-value")
	}
}

func TestLoad_BlankLinesAndComments(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "\n# This is a comment\n\n  # Another comment\nBLANK_TEST_VAR=value\n\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("BLANK_TEST_VAR", "")
	_ = os.Unsetenv("BLANK_TEST_VAR")

	if err := Load(envPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("BLANK_TEST_VAR"); got != "value" {
		t.Errorf("BLANK_TEST_VAR = %q, want %q", got, "value")
	}
}

func TestLoad_InvalidLine(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "NO_EQUALS_SIGN\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Load(envPath)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantKey string
		wantVal string
		wantErr bool
	}{
		{"simple", "KEY=value", "KEY", "value", false},
		{"with spaces", "  KEY  =  value  ", "KEY", "value", false},
		{"double quoted", "KEY=\"hello world\"", "KEY", "hello world", false},
		{"single quoted", "KEY='hello world'", "KEY", "hello world", false},
		{"empty value", "KEY=", "KEY", "", false},
		{"equals in value", "KEY=a=b=c", "KEY", "a=b=c", false},
		{"no equals", "INVALID", "", "", true},
		{"empty key", "=value", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, val, err := parseLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLine(%q): err = %v, wantErr %v", tt.line, err, tt.wantErr)
				return
			}
			if key != tt.wantKey {
				t.Errorf("parseLine(%q): key = %q, want %q", tt.line, key, tt.wantKey)
			}
			if val != tt.wantVal {
				t.Errorf("parseLine(%q): val = %q, want %q", tt.line, val, tt.wantVal)
			}
		})
	}
}
