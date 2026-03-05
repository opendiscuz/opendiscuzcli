package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "个人资料管理",
}

var profileShowCmd = &cobra.Command{
	Use:   "show [username]",
	Short: "查看用户资料 (默认: 自己)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var username string
		if len(args) > 0 {
			username = args[0]
		} else {
			creds := config.LoadCredentials()
			if creds == nil {
				return fmt.Errorf("请指定用户名或先登录")
			}
			username = creds.Username
		}

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.GET("/api/v1/users/" + username)
		if err != nil {
			return err
		}

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			var user struct {
				Username    string `json:"username"`
				DisplayName string `json:"display_name"`
				Bio         string `json:"bio"`
				UserType    string `json:"user_type"`
				Verified    bool   `json:"verified"`
				Followers   int    `json:"followers_count"`
				Following   int    `json:"following_count"`
				Posts       int    `json:"posts_count"`
			}
			json.Unmarshal(resp.Data, &user)
			fmt.Printf("👤 %s (@%s)\n", user.DisplayName, user.Username)
			if user.Bio != "" {
				fmt.Printf("   %s\n", user.Bio)
			}
			fmt.Printf("   类型: %s | 帖子: %d | 关注: %d | 粉丝: %d\n",
				user.UserType, user.Posts, user.Following, user.Followers)
		}
		return nil
	},
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新个人资料",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}

		body := map[string]interface{}{}
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			body["display_name"] = name
		}
		if bio, _ := cmd.Flags().GetString("bio"); bio != "" {
			body["bio"] = bio
		}
		if avatar, _ := cmd.Flags().GetString("avatar"); avatar != "" {
			body["avatar_url"] = avatar
		}
		if banner, _ := cmd.Flags().GetString("banner"); banner != "" {
			body["banner_url"] = banner
		}
		if locale, _ := cmd.Flags().GetString("locale"); locale != "" {
			body["locale"] = locale
		}

		if len(body) == 0 {
			return fmt.Errorf("请至少指定一个要更新的字段 (--name, --bio, --avatar, --banner, --locale)")
		}

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.PUT("/api/v1/users/me", body)
		if err != nil {
			return fmt.Errorf("update profile failed: %w", err)
		}

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Println("✅ 个人资料已更新")
		}
		return nil
	},
}

var profileSetAvatarCmd = &cobra.Command{
	Use:   "set-avatar [file-path]",
	Short: "上传并设置头像",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())

		// 1. Upload image
		uploadResp, err := client.UploadFile("/api/v1/media/upload", args[0])
		if err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}

		var uploadData struct {
			URL string `json:"url"`
		}
		json.Unmarshal(uploadResp.Data, &uploadData)

		// 2. Update profile avatar
		_, err = client.PUT("/api/v1/users/me", map[string]string{
			"avatar_url": uploadData.URL,
		})
		if err != nil {
			return fmt.Errorf("set avatar failed: %w", err)
		}

		if jsonOutput {
			fmt.Printf(`{"avatar_url":"%s"}`, uploadData.URL)
			fmt.Println()
		} else {
			fmt.Printf("✅ 头像已更新: %s\n", uploadData.URL)
		}
		return nil
	},
}

func init() {
	profileUpdateCmd.Flags().String("name", "", "显示名称")
	profileUpdateCmd.Flags().String("bio", "", "个人简介")
	profileUpdateCmd.Flags().String("avatar", "", "头像 URL")
	profileUpdateCmd.Flags().String("banner", "", "背景图 URL")
	profileUpdateCmd.Flags().String("locale", "", "语言 (zh/en/ja/ko...)")

	profileCmd.AddCommand(profileShowCmd, profileUpdateCmd, profileSetAvatarCmd)
	rootCmd.AddCommand(profileCmd)
}
