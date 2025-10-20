package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/skymoore/vibe-zsh/internal/client"
	"github.com/skymoore/vibe-zsh/internal/config"
	"github.com/skymoore/vibe-zsh/internal/formatter"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: vibe <natural language query>")
		os.Exit(1)
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

	output := formatter.Format(resp, cfg.ShowExplanation)
	fmt.Print(output)
}
