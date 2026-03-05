package cmd

import (
	"fmt"
	"strings"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "搜索帖子",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.GET("/api/v1/search?q=" + query)
		if err != nil {
			return err
		}
		printJSON(resp.DataJSON())
		return nil
	},
}

var trendsCmd = &cobra.Command{
	Use:   "trends",
	Short: "查看热门话题",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.GET("/api/v1/trends")
		if err != nil {
			return err
		}
		printJSON(resp.DataJSON())
		return nil
	},
}

var mediaUploadCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "上传图片/文件",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.UploadFile("/api/v1/media/upload", args[0])
		if err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}
		printJSON(resp.DataJSON())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd, trendsCmd, mediaUploadCmd)
}
