package cmd

import (
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/opendiscuz/opendiscuzcli/internal/i18n"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "CLI configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set config value (api-url, lang)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig()
		switch args[0] {
		case "api-url":
			cfg.APIURL = args[1]
		case "lang":
			cfg.Lang = args[1]
			i18n.SetLang(args[1])
		default:
			return fmt.Errorf("unknown config key: %s (available: api-url, lang)", args[0])
		}
		config.SaveConfig(cfg)
		if jsonOutput {
			fmt.Printf(`{"key":"%s","value":"%s"}`+"\n", args[0], args[1])
		} else {
			fmt.Printf("✅ %s = %s\n", args[0], args[1])
		}
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "View current config",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig()
		lang := cfg.Lang
		if lang == "" {
			lang = "en"
		}
		if jsonOutput {
			fmt.Printf(`{"api_url":"%s","lang":"%s"}`+"\n", cfg.APIURL, lang)
		} else {
			fmt.Printf("api-url: %s\n", cfg.APIURL)
			fmt.Printf("lang:    %s\n", lang)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd, configShowCmd)
	rootCmd.AddCommand(configCmd)
}
