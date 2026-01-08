package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/steipete/gogcli/internal/tracking"
	"github.com/steipete/gogcli/internal/ui"
)

type GmailTrackStatusCmd struct{}

func (c *GmailTrackStatusCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)

	cfg, err := tracking.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	path, _ := tracking.ConfigPath()
	if path != "" {
		u.Out().Printf("config_path\t%s", path)
	}

	if !cfg.IsConfigured() {
		u.Out().Printf("configured\tfalse")
		return nil
	}

	u.Out().Printf("configured\ttrue")
	u.Out().Printf("worker_url\t%s", cfg.WorkerURL)
	u.Out().Printf("admin_configured\t%t", strings.TrimSpace(cfg.AdminKey) != "")

	return nil
}
