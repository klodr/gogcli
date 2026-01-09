package cmd

import "testing"

func TestConfirmDestructiveForce(t *testing.T) {
	flags := &RootFlags{Force: true}
	if err := confirmDestructive(nil, flags, "delete"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestConfirmDestructiveNoInput(t *testing.T) {
	flags := &RootFlags{NoInput: true}
	err := confirmDestructive(nil, flags, "delete")
	if err == nil {
		t.Fatalf("expected error")
	}
	exitErr, ok := err.(*ExitError)
	if !ok || exitErr.Code != 2 {
		t.Fatalf("expected ExitError code 2, got %#v", err)
	}
}
