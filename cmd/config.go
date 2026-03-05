package cmd

import (
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "CLI 配置管理",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "设置配置项 (api-url)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig()
		switch args[0] {
		case "api-url":
			cfg.APIURL = args[1]
		default:
			return fmt.Errorf("unknown config key: %s (available: api-url)", args[0])
		}
		config.SaveConfig(cfg)
		if jsonOutput {
			fmt.Printf(`{"key":"%s","value":"%s"}`, args[0], args[1])
			fmt.Println()
		} else {
			fmt.Printf("✅ %s = %s\n", args[0], args[1])
		}
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "查看当前配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig()
		if jsonOutput {
			fmt.Printf(`{"api_url":"%s"}`, cfg.APIURL)
			fmt.Println()
		} else {
			fmt.Printf("api-url: %s\n", cfg.APIURL)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd, configShowCmd)
	rootCmd.AddCommand(configCmd)
}
