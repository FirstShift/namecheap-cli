package utils

import (
	"github.com/FirstShift/namecheapcli/pkg/config"
	"github.com/spf13/cobra"
)

func GetGlobalConfig(cmd *cobra.Command) *config.NamecheapCliConfig {
	// Get config from environment variables
	envConfig := config.GetConfig()

	// Initialize the final config with environment values
	finalConfig := &config.NamecheapCliConfig{
		ApiKey:     envConfig.ApiKey,
		UserName:   envConfig.UserName,
		ClientIp:   envConfig.ClientIp,
		UseSandbox: envConfig.UseSandbox,
	}

	// Check for CLI flags and use them if provided
	if apiKey, err := cmd.Flags().GetString("key"); err == nil && apiKey != "" {
		finalConfig.ApiKey = apiKey
	}

	if userName, err := cmd.Flags().GetString("user"); err == nil && userName != "" {
		finalConfig.UserName = userName
	}

	if clientIp, err := cmd.Flags().GetString("ip"); err == nil && clientIp != "" {
		finalConfig.ClientIp = clientIp
	}

	if useSandbox, err := cmd.Flags().GetBool("sandbox"); err == nil {
		finalConfig.UseSandbox = useSandbox
	}

	// Validate that required fields are set
	if finalConfig.ApiKey == "" {
		panic("API key is required. Set it via environment variable or --key flag.")
	}

	if finalConfig.UserName == "" {
		panic("Username is required. Set it via environment variable or --user flag.")
	}

	if finalConfig.ClientIp == "" {
		panic("Client IP is required. Set it via environment variable or --ip flag.")
	}

	return finalConfig
}
