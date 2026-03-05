package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/opendiscuz/opendiscuzcli/internal/i18n"
	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post operations (create/reply/like/bookmark)",
}

var postCreateCmd = &cobra.Command{
	Use:   "create [content]",
	Short: "Create a new post",
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
			return err
		}
		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			var data struct {
				ID string `json:"id"`
			}
			json.Unmarshal(resp.Data, &data)
			fmt.Printf(i18n.T("post.create.success")+"\n", data.ID)
		}
		return nil
	},
}

var postGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get post details",
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
	Short: "Reply to a post",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.POST("/api/v1/posts/"+args[0]+"/replies", map[string]string{
			"content": strings.Join(args[1:], " "),
		})
		if err != nil {
			return err
		}
		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Println(i18n.T("post.reply.success"))
		}
		return nil
	},
}

var postLikeCmd = &cobra.Command{
	Use: "like [post-id]", Short: "Like a post", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		if _, err := client.POST("/api/v1/posts/"+args[0]+"/like", nil); err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"liked"}`)
		} else {
			fmt.Println(i18n.T("post.like.success"))
		}
		return nil
	},
}

var postUnlikeCmd = &cobra.Command{
	Use: "unlike [post-id]", Short: "Unlike a post", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		if _, err := client.DELETE("/api/v1/posts/" + args[0] + "/like"); err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"unliked"}`)
		} else {
			fmt.Println(i18n.T("post.unlike.success"))
		}
		return nil
	},
}

var postBookmarkCmd = &cobra.Command{
	Use: "bookmark [post-id]", Short: "Bookmark a post", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		if _, err := client.POST("/api/v1/posts/"+args[0]+"/bookmark", nil); err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"bookmarked"}`)
		} else {
			fmt.Println(i18n.T("post.bookmark.success"))
		}
		return nil
	},
}

var postDeleteCmd = &cobra.Command{
	Use: "delete [post-id]", Short: "Delete a post", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		if _, err := client.DELETE("/api/v1/posts/" + args[0]); err != nil {
			return err
		}
		if jsonOutput {
			fmt.Println(`{"status":"deleted"}`)
		} else {
			fmt.Println(i18n.T("post.delete.success"))
		}
		return nil
	},
}

func init() {
	postCreateCmd.Flags().StringSlice("images", nil, "Image URLs")
	postCmd.AddCommand(postCreateCmd, postGetCmd, postReplyCmd, postLikeCmd, postUnlikeCmd, postBookmarkCmd, postDeleteCmd)
	rootCmd.AddCommand(postCmd)
}
