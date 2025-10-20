package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/skymoore/vibe-zsh/internal/client"
	"github.com/skymoore/vibe-zsh/internal/config"
	"github.com/skymoore/vibe-zsh/internal/formatter"
	"github.com/skymoore/vibe-zsh/internal/updater"
)

var appVersion string

func Execute(version string) {
	appVersion = version

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: vibe <natural language query>")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "--version", "-v":
		handleVersion()
		return
	case "--check-update-background":
		handleCheckUpdateBackground()
		return
	case "--update":
		handleUpdate()
		return
	}

	query := strings.Join(os.Args[1:], " ")

	cfg := config.Load()

	c := client.New(cfg)

	ctx := context.Background()

	resp, err := c.GenerateCommand(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	output := formatter.Format(resp, cfg.ShowExplanation, cfg.ShowWarnings)
	fmt.Print(output)

	updater.ShowUpdateNotification(appVersion)
}

func handleVersion() {
	fmt.Printf("vibe version %s\n", appVersion)
	updater.ShowUpdateNotification(appVersion)
}

func handleCheckUpdateBackground() {
	updater.CheckForUpdatesBackground(appVersion)
}

func handleUpdate() {
	if err := updater.PerformUpdate(appVersion); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		os.Exit(1)
	}
}
