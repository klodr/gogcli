package secrets

import (
	"strings"
	"testing"
)

func TestFileKeyringPasswordFuncFrom_UsesEnvPassword(t *testing.T) {
	prompt := fileKeyringPasswordFuncFrom("pw", false)
	got, err := prompt("ignored")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got != "pw" {
		t.Fatalf("expected password, got: %q", got)
	}
}

func TestFileKeyringPasswordFuncFrom_NoTTYErrors(t *testing.T) {
	prompt := fileKeyringPasswordFuncFrom("", false)
	if _, err := prompt("ignored"); err == nil {
		t.Fatalf("expected error")
	} else if !strings.Contains(err.Error(), keyringPasswordEnv) {
		t.Fatalf("expected error mentioning %s, got: %v", keyringPasswordEnv, err)
	}
}
