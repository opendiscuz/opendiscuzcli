package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/opendiscuz/opendiscuzcli/internal/i18n"
	"github.com/spf13/cobra"
)

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Timeline (trending/home)",
}

var trendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "View trending posts",
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
				fmt.Println(i18n.T("timeline.empty"))
			}
		}
		return nil
	},
}

var homeCmd = &cobra.Command{
	Use:   "home",
	Short: "View posts from followed users",
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
	trendingCmd.Flags().Int("limit", 20, "Number of results")
	homeCmd.Flags().Int("limit", 20, "Number of results")
	timelineCmd.AddCommand(trendingCmd, homeCmd)
	rootCmd.AddCommand(timelineCmd)
}
