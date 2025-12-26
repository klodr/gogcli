package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfirmDestructive(t *testing.T) {
	cmd := &cobra.Command{}

	// Force: always ok.
	if err := confirmDestructive(cmd, &rootFlags{Force: true}, "delete thing"); err != nil {
		t.Fatalf("force: %v", err)
	}

	// NoInput: should refuse without force.
	err := confirmDestructive(cmd, &rootFlags{NoInput: true}, "delete thing")
	if err == nil || !strings.Contains(err.Error(), "refusing") {
		t.Fatalf("expected refusal, got: %v", err)
	}
}
