package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/steipete/gogcli/internal/config"
	"github.com/steipete/gogcli/internal/outfmt"
	"github.com/steipete/gogcli/internal/ui"
)

type AuthKeyringCmd struct {
	Set AuthKeyringSetCmd `cmd:"" name:"set" help:"Set keyring backend (writes config.json)"`
}

type AuthKeyringSetCmd struct {
	Backend string `arg:"" name:"backend" help:"Keyring backend: auto|keychain|file"`
}

func (c *AuthKeyringSetCmd) Run(ctx context.Context) error {
	u := ui.FromContext(ctx)

	backend := strings.ToLower(strings.TrimSpace(c.Backend))
	switch backend {
	case "default":
		backend = "auto"
	case "auto", "keychain", "file":
	default:
		return usagef("invalid backend: %q (expected auto, keychain, or file)", c.Backend)
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}
	cfg.KeyringBackend = backend
	if err := config.WriteConfig(cfg); err != nil {
		return err
	}

	path, _ := config.ConfigPath()

	// Env var wins; warn so it doesn't look "broken".
	if v := strings.TrimSpace(os.Getenv("GOG_KEYRING_BACKEND")); v != "" &&
		u != nil &&
		!outfmt.IsJSON(ctx) &&
		!outfmt.IsPlain(ctx) {
		u.Err().Printf("NOTE: GOG_KEYRING_BACKEND=%s overrides config.json", v)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"written":         true,
			"path":            path,
			"keyring_backend": backend,
		})
	}

	if u == nil {
		return nil
	}

	u.Out().Printf("written\ttrue")
	u.Out().Printf("path\t%s", path)
	u.Out().Printf("keyring_backend\t%s", backend)
	return nil
}
