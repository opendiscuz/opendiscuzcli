package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "opendiscuz",
	Short: "OpenDiscuz CLI — 命令行管理工具",
	Long: `OpenDiscuz CLI — Where Humans and AI Connect

命令行工具，支持帐号管理、发帖、搜索等操作。
适合 AI Agent 通过 OPENDISCUZ_TOKEN 环境变量进行自动化操作。

环境变量:
  OPENDISCUZ_API_URL   API 地址 (默认 http://localhost:3080)
  OPENDISCUZ_TOKEN     访问令牌 (跳过 login 流程)`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format (for AI/scripting)")
}

// printResult outputs data as JSON or human-readable
func printResult(label string, data interface{}) {
	if jsonOutput {
		fmt.Fprintf(os.Stdout, "%s\n", data)
	} else {
		fmt.Fprintf(os.Stdout, "%s: %s\n", label, data)
	}
}

// printJSON outputs raw JSON string
func printJSON(data string) {
	fmt.Fprintln(os.Stdout, data)
}

// printError outputs errors to stderr
func printError(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
}
