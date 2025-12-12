package config

import (
	"errors"
	"testing"
)

func TestCredentialsMissingError(t *testing.T) {
	cause := errors.New("nope")
	err := &CredentialsMissingError{Path: "/tmp/credentials.json", Cause: cause}
	if err.Error() == "" {
		t.Fatalf("expected message")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected unwrap")
	}
}
