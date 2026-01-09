package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainUpdatesReadme(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	readme := filepath.Join(dir, "README.md")
	content := "# Test\n" + startMarker + "\n" + endMarker + "\n"
	if err := os.WriteFile(readme, []byte(content), 0o600); err != nil {
		t.Fatalf("write README: %v", err)
	}

	main()

	updated, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	text := string(updated)
	if !strings.Contains(text, startMarker) || !strings.Contains(text, endMarker) {
		t.Fatalf("missing markers: %q", text)
	}
	if !strings.Contains(text, "|") {
		t.Fatalf("expected markdown table: %q", text)
	}
}
