package cmd

import (
	"fmt"
	"os"

	"github.com/skymoore/vibe-zsh/internal/history"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show interactive query history",
	Long: `Display an interactive menu of previous queries.

Note: This command outputs to stdout. For shell integration, use:
  - 'vh' command (inserts into buffer)
  - Ctrl+X H keybinding (inserts into buffer)

Direct usage of 'vibe-zsh history' is mainly for scripting.`,
	Run: func(cmd *cobra.Command, args []string) {
		showInteractiveHistory()
	},
}

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List query history in plain text",
	Run: func(cmd *cobra.Command, args []string) {
		listHistory()
	},
}

var historyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all query history",
	Run: func(cmd *cobra.Command, args []string) {
		clearHistory()
	},
}

func init() {
	historyCmd.AddCommand(historyListCmd)
	historyCmd.AddCommand(historyClearCmd)
	rootCmd.AddCommand(historyCmd)
}

func showInteractiveHistory() {
	if !cfg.EnableHistory {
		fmt.Fprintln(os.Stderr, "History is disabled. Set VIBE_ENABLE_HISTORY=true to enable.")
		os.Exit(1)
	}

	h, err := history.New(cfg.CacheDir, cfg.HistorySize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing history: %v\n", err)
		os.Exit(1)
	}

	entries, err := h.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading history: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "No history entries found.")
		fmt.Fprintln(os.Stderr, "Start using vibe-zsh to build your query history!")
		os.Exit(0)
	}

	command, err := history.ShowInteractive(entries)
	if err != nil {
		// If not in TTY, fall back to plain list
		if err.Error() == "not in a TTY, use 'vibe-zsh history list' instead" {
			listHistory()
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if command != "" {
		// Output the selected command to stdout
		fmt.Print(command)
	}
}

func listHistory() {
	if !cfg.EnableHistory {
		fmt.Fprintln(os.Stderr, "History is disabled. Set VIBE_ENABLE_HISTORY=true to enable.")
		os.Exit(1)
	}

	h, err := history.New(cfg.CacheDir, cfg.HistorySize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing history: %v\n", err)
		os.Exit(1)
	}

	entries, err := h.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading history: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(history.FormatPlainList(entries))
}

func clearHistory() {
	if !cfg.EnableHistory {
		fmt.Fprintln(os.Stderr, "History is disabled. Set VIBE_ENABLE_HISTORY=true to enable.")
		os.Exit(1)
	}

	h, err := history.New(cfg.CacheDir, cfg.HistorySize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing history: %v\n", err)
		os.Exit(1)
	}

	if err := h.Clear(); err != nil {
		fmt.Fprintf(os.Stderr, "Error clearing history: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("History cleared successfully.")
}
