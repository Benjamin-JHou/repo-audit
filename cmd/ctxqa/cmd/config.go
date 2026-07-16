package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ctxqa/ctxqa/pkg/config"
	"github.com/ctxqa/ctxqa/pkg/renderer"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [command]",
	Short: "Manage configuration",
	Long:  "Manage ctxqa configuration including API key, provider, model, and audit preferences.",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration interactively",
	Long:  "Run an interactive wizard to configure API key, provider, model, and other settings.",
	Run:   runConfigInit,
}

var configSetProviderCmd = &cobra.Command{
	Use:   "set-provider [openai|anthropic]",
	Short: "Set API provider",
	Run:   runConfigSetProvider,
}

var configSetModelCmd = &cobra.Command{
	Use:   "set-model [model-name]",
	Short: "Set model name",
	Run:   runConfigSetModel,
}

var configSetAPIKeyCmd = &cobra.Command{
	Use:   "set-api-key",
	Short: "Set API key interactively",
	Run:   runConfigSetAPIKey,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run:   runConfigShow,
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Run:   runConfigReset,
}

var configExcludeAddCmd = &cobra.Command{
	Use:   "exclude-add [patterns...]",
	Short: "Add file exclusion patterns",
	Run:   runConfigExcludeAdd,
}

var configExcludeRemoveCmd = &cobra.Command{
	Use:   "exclude-remove [patterns...]",
	Short: "Remove file exclusion patterns",
	Run:   runConfigExcludeRemove,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configSetProviderCmd)
	configCmd.AddCommand(configSetModelCmd)
	configCmd.AddCommand(configSetAPIKeyCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configExcludeAddCmd)
	configCmd.AddCommand(configExcludeRemoveCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) {
	renderer.PrintWelcome(version)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== ctxqa Configuration Wizard ===")
	fmt.Println()

	fmt.Print("Enter API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	fmt.Print("Select provider (openai/anthropic) [openai]: ")
	line, _ := reader.ReadString('\n')
	provider := strings.TrimSpace(line)
	if provider == "" {
		provider = "openai"
	}

	fmt.Printf("Enter Base URL [%s]: ", getDefaultBaseURL(provider))
	line, _ = reader.ReadString('\n')
	baseURL := strings.TrimSpace(line)
	if baseURL == "" {
		baseURL = getDefaultBaseURL(provider)
	}

	fmt.Printf("Enter model [%s]: ", getDefaultModel(provider))
	line, _ = reader.ReadString('\n')
	model := strings.TrimSpace(line)
	if model == "" {
		model = getDefaultModel(provider)
	}

	cfg := &config.Config{
		Provider:       provider,
		APIKey:         apiKey,
		BaseURL:        baseURL,
		Model:          model,
		MaxFileSize:    config.DefaultMaxFileSize,
		BatchSize:      config.DefaultBatchSize,
		MaxBatchChars:  config.DefaultMaxBatchChars,
		TimeoutSeconds: config.DefaultTimeout,
		RetryCount:     config.DefaultRetryCount,
	}

	if err := cfg.Save(""); err != nil {
		renderer.PrintError("Failed to save config: " + err.Error())
		os.Exit(1)
	}

	fmt.Println()
	renderer.PrintSuccess("Configuration saved successfully!")
}

func runConfigSetProvider(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load("")
	if len(args) > 0 {
		cfg.Provider = args[0]
	} else {
		fmt.Print("Enter provider (openai/anthropic): ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		cfg.Provider = strings.TrimSpace(line)
	}
	cfg.BaseURL = getDefaultBaseURL(cfg.Provider)
	cfg.Model = getDefaultModel(cfg.Provider)
	cfg.Save("")
	renderer.PrintSuccess("Provider set to " + cfg.Provider)
}

func runConfigSetModel(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load("")
	if len(args) > 0 {
		cfg.Model = args[0]
	} else {
		fmt.Print("Enter model name: ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		cfg.Model = strings.TrimSpace(line)
	}
	cfg.Save("")
	renderer.PrintSuccess("Model set to " + cfg.Model)
}

func runConfigSetAPIKey(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load("")
	if len(args) > 0 {
		cfg.APIKey = args[0]
	} else {
		fmt.Print("Enter API Key: ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		cfg.APIKey = strings.TrimSpace(line)
	}
	cfg.Save("")
	renderer.PrintSuccess("API Key updated")
}

func runConfigShow(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load("")
	style := lipgloss.NewStyle().PaddingLeft(2)

	fmt.Println(style.Render("Provider:") + " " + cfg.Provider)
	fmt.Println(style.Render("Base URL:") + " " + cfg.BaseURL)
	fmt.Println(style.Render("Model:") + " " + cfg.Model)
	fmt.Println(style.Render("API Key:") + " " + maskKey(cfg.APIKey))
	fmt.Println(style.Render("Max File Size:") + " " + fmt.Sprintf("%d bytes", cfg.MaxFileSize))
	fmt.Println(style.Render("Batch Size:") + " " + fmt.Sprintf("%d files", cfg.BatchSize))
	fmt.Println(style.Render("Max Batch Chars:") + " " + fmt.Sprintf("%d", cfg.MaxBatchChars))
	fmt.Println(style.Render("Timeout:") + " " + fmt.Sprintf("%ds", cfg.TimeoutSeconds))
	fmt.Println(style.Render("Retry Count:") + " " + fmt.Sprintf("%d", cfg.RetryCount))
}

func runConfigReset(cmd *cobra.Command, args []string) {
	cfg := &config.Config{
		Provider:       config.DefaultProvider,
		BaseURL:        config.DefaultBaseURL,
		Model:          config.DefaultModel,
		MaxFileSize:    config.DefaultMaxFileSize,
		BatchSize:      config.DefaultBatchSize,
		MaxBatchChars:  config.DefaultMaxBatchChars,
		TimeoutSeconds: config.DefaultTimeout,
		RetryCount:     config.DefaultRetryCount,
	}
	cfg.Save("")
	renderer.PrintSuccess("Configuration reset to defaults")
}

func runConfigExcludeAdd(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load("")
	cfg.Context.Exclude = append(cfg.Context.Exclude, args...)
	cfg.Save("")
	renderer.PrintSuccess("Added " + fmt.Sprintf("%d", len(args)) + " exclusion pattern(s)")
}

func runConfigExcludeRemove(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load("")
	for _, pattern := range args {
		for i, p := range cfg.Context.Exclude {
			if p == pattern {
				cfg.Context.Exclude = append(cfg.Context.Exclude[:i], cfg.Context.Exclude[i+1:]...)
				break
			}
		}
	}
	cfg.Save("")
	renderer.PrintSuccess("Removed " + fmt.Sprintf("%d", len(args)) + " exclusion pattern(s)")
}

func getDefaultBaseURL(provider string) string {
	switch provider {
	case "anthropic":
		return "https://api.anthropic.com/v1"
	default:
		return "https://api.openai.com/v1"
	}
}

func getDefaultModel(provider string) string {
	switch provider {
	case "anthropic":
		return "claude-sonnet-4-20250514"
	default:
		return "gpt-4o"
	}
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
