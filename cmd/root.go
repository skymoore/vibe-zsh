package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/skymoore/vibe-zsh/internal/client"
	"github.com/skymoore/vibe-zsh/internal/config"
	"github.com/skymoore/vibe-zsh/internal/history"
	"github.com/skymoore/vibe-zsh/internal/logger"
	"github.com/skymoore/vibe-zsh/internal/progress"
	"github.com/skymoore/vibe-zsh/internal/streamer"
	"github.com/skymoore/vibe-zsh/internal/updater"
	"github.com/spf13/cobra"
)

var (
	appVersion string
	cfg        *config.Config

	apiURL               string
	apiKey               string
	model                string
	temperature          float64
	maxTokens            int
	timeout              time.Duration
	useStructuredOutput  bool
	showExplanation      bool
	showWarnings         bool
	enableCache          bool
	cacheTTL             time.Duration
	interactiveMode      bool
	maxRetries           int
	enableJSONExtraction bool
	strictValidation     bool
	debugLogs            bool
	showRetryStatus      bool
	showProgress         bool
	progressStyle        string
	streamOutput         bool
	streamDelay          time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "vibe-zsh [query]",
	Short: "Transform natural language into shell commands using AI",
	Long: `vibe-zsh is a CLI tool that converts natural language descriptions into executable shell commands.

Simply describe what you want to do in plain English, and vibe-zsh will generate the appropriate
command using AI. It works with any OpenAI-compatible API (Ollama, LM Studio, OpenAI, Claude, etc.).

Examples:
  vibe-zsh "list all docker containers"
  vibe-zsh "find all python files modified today"
  vibe-zsh "show me commits from last week"

For more information, visit: https://github.com/skymoore/vibe-zsh`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		generateCommand(query)
	},
	SilenceUsage: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of vibe-zsh",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vibe-zsh version %s\n", appVersion)
		updater.ShowUpdateNotification(appVersion)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update vibe-zsh to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		if err := updater.PerformUpdate(appVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var checkUpdateCmd = &cobra.Command{
	Use:    "check-update",
	Short:  "Check for available updates in the background",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		updater.CheckForUpdatesBackground(appVersion)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API endpoint URL (default: http://localhost:11434/v1)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API authentication key")
	rootCmd.PersistentFlags().StringVar(&model, "model", "", "Model to use (default: llama3:8b)")
	rootCmd.PersistentFlags().Float64Var(&temperature, "temperature", -1, "Generation temperature 0.0-2.0 (default: 0.2)")
	rootCmd.PersistentFlags().IntVar(&maxTokens, "max-tokens", -1, "Maximum response tokens (default: 500)")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 0, "Request timeout (default: 30s)")

	rootCmd.PersistentFlags().BoolVar(&useStructuredOutput, "structured-output", true, "Use JSON schema for responses")
	rootCmd.PersistentFlags().BoolVar(&showExplanation, "explanation", true, "Show command explanations")
	rootCmd.PersistentFlags().BoolVar(&showWarnings, "warnings", true, "Show warnings for dangerous commands")
	rootCmd.PersistentFlags().BoolVar(&interactiveMode, "interactive", false, "Confirm before executing commands")
	rootCmd.PersistentFlags().BoolVar(&enableCache, "cache", true, "Enable response caching")
	rootCmd.PersistentFlags().DurationVar(&cacheTTL, "cache-ttl", 0, "Cache lifetime (default: 24h)")

	rootCmd.PersistentFlags().IntVar(&maxRetries, "max-retries", -1, "Max retry attempts (default: 3)")
	rootCmd.PersistentFlags().BoolVar(&enableJSONExtraction, "json-extraction", true, "Extract JSON from corrupted responses")
	rootCmd.PersistentFlags().BoolVar(&strictValidation, "strict-validation", true, "Validate response structure")
	rootCmd.PersistentFlags().BoolVar(&debugLogs, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&showRetryStatus, "retry-status", true, "Show retry progress")

	rootCmd.PersistentFlags().BoolVar(&showProgress, "progress", true, "Show progress spinner")
	rootCmd.PersistentFlags().StringVar(&progressStyle, "progress-style", "", "Spinner style: dots, line, circle, bounce, arrow, runes (default: dots)")
	rootCmd.PersistentFlags().BoolVar(&streamOutput, "stream", true, "Stream output with typewriter effect")
	rootCmd.PersistentFlags().DurationVar(&streamDelay, "stream-delay", 0, "Delay between streamed words (default: 20ms)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(checkUpdateCmd)
}

func initConfig() {
	cfg = config.Load()

	if apiURL != "" {
		cfg.APIURL = apiURL
	}
	if apiKey != "" {
		cfg.APIKey = apiKey
	}
	if model != "" {
		cfg.Model = model
	}
	if temperature >= 0 {
		cfg.Temperature = temperature
	}
	if maxTokens > 0 {
		cfg.MaxTokens = maxTokens
	}
	if timeout > 0 {
		cfg.Timeout = timeout
	}
	if cacheTTL > 0 {
		cfg.CacheTTL = cacheTTL
	}
	if maxRetries >= 0 {
		cfg.MaxRetries = maxRetries
	}

	if cmd, _, err := rootCmd.Find(os.Args[1:]); err == nil && cmd.Use != "vibe-zsh [query]" {
		return
	}

	if rootCmd.PersistentFlags().Changed("structured-output") {
		cfg.UseStructuredOutput = useStructuredOutput
	}
	if rootCmd.PersistentFlags().Changed("explanation") {
		cfg.ShowExplanation = showExplanation
	}
	if rootCmd.PersistentFlags().Changed("warnings") {
		cfg.ShowWarnings = showWarnings
	}
	if rootCmd.PersistentFlags().Changed("interactive") {
		cfg.InteractiveMode = interactiveMode
	}
	if rootCmd.PersistentFlags().Changed("cache") {
		cfg.EnableCache = enableCache
	}
	if rootCmd.PersistentFlags().Changed("json-extraction") {
		cfg.EnableJSONExtraction = enableJSONExtraction
	}
	if rootCmd.PersistentFlags().Changed("strict-validation") {
		cfg.StrictValidation = strictValidation
	}
	if rootCmd.PersistentFlags().Changed("debug") {
		cfg.EnableDebugLogs = debugLogs
	}
	if rootCmd.PersistentFlags().Changed("retry-status") {
		cfg.ShowRetryStatus = showRetryStatus
	}
	if rootCmd.PersistentFlags().Changed("progress") {
		cfg.ShowProgress = showProgress
	}
	if rootCmd.PersistentFlags().Changed("progress-style") && progressStyle != "" {
		cfg.ProgressStyle = parseProgressStyle(progressStyle)
	}
	if rootCmd.PersistentFlags().Changed("stream") {
		cfg.StreamOutput = streamOutput
	}
	if rootCmd.PersistentFlags().Changed("stream-delay") {
		cfg.StreamDelay = streamDelay
	}

	logger.Init(cfg.EnableDebugLogs)
}

func parseProgressStyle(style string) progress.SpinnerStyle {
	switch strings.ToLower(style) {
	case "dots":
		return progress.StyleDots
	case "line":
		return progress.StyleLine
	case "circle":
		return progress.StyleCircle
	case "bounce":
		return progress.StyleBounce
	case "arrow":
		return progress.StyleArrow
	case "runes":
		return progress.StyleRunes
	default:
		return progress.StyleDots
	}
}

// cleanExplanation removes terminal escape codes and problematic Unicode characters
func cleanExplanation(s string) string {
	// Remove ANSI escape codes (like bracketed paste mode)
	// Pattern: ESC [ ... (any characters) ... letter
	result := ""
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1B && i+1 < len(s) && s[i+1] == '[' {
			// Start of ANSI escape sequence
			inEscape = true
			i++ // Skip the '['
			continue
		}
		if inEscape {
			// Skip until we find a letter (end of escape sequence)
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') || s[i] == '~' {
				inEscape = false
			}
			continue
		}
		result += string(s[i])
	}

	// Replace problematic Unicode characters
	replacer := strings.NewReplacer(
		"\u2026", "...", // Ellipsis → three dots
		"\u2013", "-", // En dash → hyphen
		"\u2014", "--", // Em dash → double hyphen
		"\u00a0", " ", // Non-breaking space → regular space
	)
	result = replacer.Replace(result)

	// Trim whitespace
	result = strings.TrimSpace(result)

	// Filter out garbage explanations (too many ellipses, question marks, or very short)
	if isGarbageExplanation(result) {
		return ""
	}

	return result
}

// isGarbageExplanation detects if an explanation is corrupted/garbage
func isGarbageExplanation(s string) bool {
	// Too short to be useful
	if len(s) < 10 {
		return true
	}

	// Count problematic characters
	ellipsisCount := strings.Count(s, "...")
	questionCount := strings.Count(s, "???")

	// Too many ellipses or question marks indicates garbage
	if ellipsisCount > 2 || questionCount > 0 {
		return true
	}

	// Check for excessive Unicode ellipsis characters (even after replacement)
	unicodeEllipsisCount := 0
	for _, r := range s {
		if r == '\u2026' {
			unicodeEllipsisCount++
		}
	}
	return unicodeEllipsisCount > 2
}

func generateCommand(query string) {
	// Setup context with cancellation for Ctrl+C handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()     // Cancel context, spinner stops via ctx.Done()
		os.Exit(130) // Standard exit code for SIGINT
	}()

	c := client.New(cfg)

	resp, err := c.GenerateCommand(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Debug: Print full JSON response as comments
	if cfg.EnableDebugLogs {
		logger.Debug("=== Full Response ===")
		logger.Debug("Command: %s", resp.Command)
		logger.Debug("Explanation count: %d", len(resp.Explanation))
		for i, exp := range resp.Explanation {
			logger.Debug("Explanation[%d]: %q", i, exp)
			logger.Debug("Explanation[%d] hex: % x", i, exp)
		}
		logger.Debug("Warning: %s", resp.Warning)
		logger.Debug("====================")
	}

	// Show explanations to stderr if enabled (user sees these while command loads into buffer)
	if cfg.ShowExplanation && len(resp.Explanation) > 0 {
		hasGarbage := false
		validExplanations := 0

		for _, line := range resp.Explanation {
			// Clean the explanation line: remove escape codes
			cleanLine := cleanExplanation(line)
			if cleanLine != "" {
				// Stream each explanation line with typewriter effect
				if cfg.StreamOutput && cfg.ShowProgress {
					fmt.Fprint(os.Stderr, "# ")
					if err := streamer.StreamWord(os.Stderr, cleanLine, cfg.StreamDelay); err != nil {
						fmt.Fprint(os.Stderr, cleanLine)
					}
					fmt.Fprintln(os.Stderr)
				} else {
					fmt.Fprintf(os.Stderr, "# %s\n", cleanLine)
				}
				validExplanations++
			} else {
				hasGarbage = true
			}
		}

		// Warn if we detected garbage explanations
		if hasGarbage || validExplanations == 0 {
			fmt.Fprintln(os.Stderr, "#")
			fmt.Fprintln(os.Stderr, "# ⚠️  Model generated incomplete explanations")
			fmt.Fprintln(os.Stderr, "# Try a different model or increase VIBE_MAX_TOKENS")
		}

		if cfg.ShowWarnings && resp.Warning != "" {
			cleanWarning := cleanExplanation(resp.Warning)
			if cleanWarning != "" {
				if cfg.StreamOutput && cfg.ShowProgress {
					fmt.Fprint(os.Stderr, "# WARNING: ")
					if err := streamer.StreamWord(os.Stderr, cleanWarning, cfg.StreamDelay); err != nil {
						fmt.Fprint(os.Stderr, cleanWarning)
					}
					fmt.Fprintln(os.Stderr)
				} else {
					fmt.Fprintf(os.Stderr, "# WARNING: %s\n", cleanWarning)
				}
			}
		}
	}

	// Ensure all stderr output is flushed before writing to stdout
	os.Stderr.Sync()

	// Output only the command to stdout (this is what ZSH captures for the buffer)
	fmt.Print(resp.Command)

	// Save to history if enabled
	if cfg.EnableHistory {
		h, err := history.New(cfg.CacheDir, cfg.HistorySize)
		if err == nil {
			// Ignore errors when saving to history - don't fail the command
			_ = h.Add(query, resp.Command)
		}
	}

	updater.ShowUpdateNotification(appVersion)
}

func Execute(version string) {
	appVersion = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
