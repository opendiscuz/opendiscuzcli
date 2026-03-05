package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/spf13/cobra"
)

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "时间线 (热门/关注)",
}

var trendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "查看热门帖子",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.GET(fmt.Sprintf("/api/v1/timeline/trending?limit=%d", limit))
		if err != nil {
			return err
		}

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			var posts []struct {
				ID      string `json:"id"`
				Content string `json:"content"`
				Author  struct {
					Username    string `json:"username"`
					DisplayName string `json:"display_name"`
				} `json:"author"`
				LikesCount   int `json:"likes_count"`
				RepliesCount int `json:"replies_count"`
			}
			json.Unmarshal(resp.Data, &posts)
			for _, p := range posts {
				fmt.Printf("📝 @%s: %s\n", p.Author.Username, truncate(p.Content, 80))
				fmt.Printf("   ❤️ %d  💬 %d  [%s]\n\n", p.LikesCount, p.RepliesCount, p.ID[:8])
			}
			if len(posts) == 0 {
				fmt.Println("暂无热门帖子")
			}
		}
		return nil
	},
}

var homeCmd = &cobra.Command{
	Use:   "home",
	Short: "查看关注者的帖子",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		limit, _ := cmd.Flags().GetInt("limit")
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.GET(fmt.Sprintf("/api/v1/timeline/home?limit=%d", limit))
		if err != nil {
			return err
		}
		printJSON(resp.DataJSON())
		return nil
	},
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}

func init() {
	trendingCmd.Flags().Int("limit", 20, "返回数量")
	homeCmd.Flags().Int("limit", 20, "返回数量")

	timelineCmd.AddCommand(trendingCmd, homeCmd)
	rootCmd.AddCommand(timelineCmd)
}
