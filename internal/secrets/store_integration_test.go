package secrets

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/99designs/keyring"

	"github.com/steipete/gogcli/internal/config"
)

func setupKeyringEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg"))
	t.Setenv("GOG_KEYRING_BACKEND", "file")
	t.Setenv("GOG_KEYRING_PASSWORD", "testpass")
}

func TestSetAndGetSecret_FileBackend(t *testing.T) {
	setupKeyringEnv(t)

	if err := SetSecret("test/key", []byte("value")); err != nil {
		t.Fatalf("SetSecret: %v", err)
	}
	val, err := GetSecret("test/key")
	if err != nil {
		t.Fatalf("GetSecret: %v", err)
	}
	if string(val) != "value" {
		t.Fatalf("unexpected value: %q", val)
	}
}

func TestKeyringStore_TokenRoundTrip(t *testing.T) {
	ring := keyring.NewArrayKeyring(nil)
	store := &KeyringStore{ring: ring}

	tok := Token{RefreshToken: "rt", Services: []string{"gmail"}, Scopes: []string{"s"}, CreatedAt: time.Now()}
	if err := store.SetToken("a@b.com", tok); err != nil {
		t.Fatalf("SetToken: %v", err)
	}
	got, err := store.GetToken("a@b.com")
	if err != nil {
		t.Fatalf("GetToken: %v", err)
	}
	if got.RefreshToken != "rt" {
		t.Fatalf("unexpected token: %#v", got)
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if len(keys) == 0 {
		t.Fatalf("expected keys")
	}
}

func TestEnsureKeyringDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg"))

	if _, err := config.EnsureKeyringDir(); err != nil {
		t.Fatalf("EnsureKeyringDir: %v", err)
	}
}
