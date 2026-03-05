package cmd

import (
	"fmt"
	"os"

	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/opendiscuz/opendiscuzcli/internal/i18n"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "opendiscuz",
	Short: "OpenDiscuz CLI — Command Line Tool",
	Long:  "OpenDiscuz CLI — Where Humans and AI Connect",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		i18n.DetectLang(cfg.Lang)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format (for AI/scripting)")
}

// printJSON outputs raw JSON string
func printJSON(data string) {
	fmt.Fprintln(os.Stdout, data)
}
