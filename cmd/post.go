package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "帖子操作 (发帖/回复/点赞/收藏)",
}

var postCreateCmd = &cobra.Command{
	Use:   "create [content]",
	Short: "发新帖",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}

		content := strings.Join(args, " ")
		images, _ := cmd.Flags().GetStringSlice("images")

		body := map[string]interface{}{"content": content}
		if len(images) > 0 {
			body["images"] = images
		}

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.POST("/api/v1/posts", body)
		if err != nil {
			return fmt.Errorf("create post failed: %w", err)
		}

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			var data struct {
				ID string `json:"id"`
			}
			json.Unmarshal(resp.Data, &data)
			fmt.Printf("✅ 帖子已发布 (ID: %s)\n", data.ID)
		}
		return nil
	},
}

var postGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "获取帖子详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.GET("/api/v1/posts/" + args[0])
		if err != nil {
			return err
		}
		printJSON(resp.DataJSON())
		return nil
	},
}

var postReplyCmd = &cobra.Command{
	Use:   "reply [post-id] [content]",
	Short: "回复帖子",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}

		postID := args[0]
		content := strings.Join(args[1:], " ")

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.POST("/api/v1/posts/"+postID+"/replies", map[string]string{
			"content": content,
		})
		if err != nil {
			return fmt.Errorf("reply failed: %w", err)
		}

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Printf("✅ 回复已发送\n")
		}
		return nil
	},
}

var postLikeCmd = &cobra.Command{
	Use:   "like [post-id]",
	Short: "点赞帖子",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		_, err := client.POST("/api/v1/posts/"+args[0]+"/like", nil)
		if err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"liked"}`)
		} else {
			fmt.Println("❤️ 已点赞")
		}
		return nil
	},
}

var postUnlikeCmd = &cobra.Command{
	Use:   "unlike [post-id]",
	Short: "取消点赞",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		_, err := client.DELETE("/api/v1/posts/" + args[0] + "/like")
		if err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"unliked"}`)
		} else {
			fmt.Println("💔 已取消点赞")
		}
		return nil
	},
}

var postBookmarkCmd = &cobra.Command{
	Use:   "bookmark [post-id]",
	Short: "收藏帖子",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		_, err := client.POST("/api/v1/posts/"+args[0]+"/bookmark", nil)
		if err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"bookmarked"}`)
		} else {
			fmt.Println("🔖 已收藏")
		}
		return nil
	},
}

var postDeleteCmd = &cobra.Command{
	Use:   "delete [post-id]",
	Short: "删除帖子",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		_, err := client.DELETE("/api/v1/posts/" + args[0])
		if err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"deleted"}`)
		} else {
			fmt.Println("🗑️ 帖子已删除")
		}
		return nil
	},
}

func init() {
	postCreateCmd.Flags().StringSlice("images", nil, "图片 URL 列表")

	postCmd.AddCommand(postCreateCmd, postGetCmd, postReplyCmd, postLikeCmd, postUnlikeCmd, postBookmarkCmd, postDeleteCmd)
	rootCmd.AddCommand(postCmd)
}
