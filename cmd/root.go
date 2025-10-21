package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/skymoore/vibe-zsh/internal/client"
	"github.com/skymoore/vibe-zsh/internal/config"
	"github.com/skymoore/vibe-zsh/internal/formatter"
	"github.com/skymoore/vibe-zsh/internal/logger"
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
)

var rootCmd = &cobra.Command{
	Use:   "vibe [query]",
	Short: "Transform natural language into shell commands using AI",
	Long: `Vibe is a CLI tool that converts natural language descriptions into executable shell commands.

Simply describe what you want to do in plain English, and vibe will generate the appropriate
command using AI. It works with any OpenAI-compatible API (Ollama, OpenAI, Claude, etc.).

Examples:
  vibe "list all docker containers"
  vibe "find all python files modified today"
  vibe "show me commits from last week"

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
	Short: "Print the version number of vibe",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vibe version %s\n", appVersion)
		updater.ShowUpdateNotification(appVersion)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update vibe to the latest version",
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

	if cmd, _, err := rootCmd.Find(os.Args[1:]); err == nil && cmd.Use != "vibe [query]" {
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

	logger.Init(cfg.EnableDebugLogs)
}

func generateCommand(query string) {
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

func Execute(version string) {
	appVersion = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
