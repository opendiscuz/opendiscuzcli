package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/opendiscuz/opendiscuzcli/internal/i18n"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Profile management",
}

var profileShowCmd = &cobra.Command{
	Use:   "show [username]",
	Short: "View user profile (default: self)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var username string
		if len(args) > 0 {
			username = args[0]
		} else {
			creds := config.LoadCredentials()
			if creds == nil {
				return fmt.Errorf("please specify a username or login first")
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
				Followers   int    `json:"followers_count"`
				Following   int    `json:"following_count"`
				Posts       int    `json:"posts_count"`
			}
			json.Unmarshal(resp.Data, &user)
			fmt.Printf("👤 %s (@%s)\n", user.DisplayName, user.Username)
			if user.Bio != "" {
				fmt.Printf("   %s\n", user.Bio)
			}
			fmt.Printf(i18n.T("profile.show.type")+"\n",
				user.UserType, user.Posts, user.Following, user.Followers)
		}
		return nil
	},
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["display_name"] = v
		}
		if v, _ := cmd.Flags().GetString("bio"); v != "" {
			body["bio"] = v
		}
		if v, _ := cmd.Flags().GetString("avatar"); v != "" {
			body["avatar_url"] = v
		}
		if v, _ := cmd.Flags().GetString("banner"); v != "" {
			body["banner_url"] = v
		}
		if v, _ := cmd.Flags().GetString("locale"); v != "" {
			body["locale"] = v
		}
		if len(body) == 0 {
			return fmt.Errorf(i18n.T("profile.update.empty"))
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.PUT("/api/v1/users/me", body)
		if err != nil {
			return err
		}
		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Println(i18n.T("profile.update.success"))
		}
		return nil
	},
}

var profileSetAvatarCmd = &cobra.Command{
	Use:   "set-avatar [file-path]",
	Short: "Upload and set avatar",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		uploadResp, err := client.UploadFile("/api/v1/media/upload", args[0])
		if err != nil {
			return err
		}
		var uploadData struct {
			URL string `json:"url"`
		}
		json.Unmarshal(uploadResp.Data, &uploadData)
		if _, err = client.PUT("/api/v1/users/me", map[string]string{"avatar_url": uploadData.URL}); err != nil {
			return err
		}
		if jsonOutput {
			fmt.Printf(`{"avatar_url":"%s"}`+"\n", uploadData.URL)
		} else {
			fmt.Printf(i18n.T("profile.avatar.success")+"\n", uploadData.URL)
		}
		return nil
	},
}

func init() {
	profileUpdateCmd.Flags().String("name", "", "Display name")
	profileUpdateCmd.Flags().String("bio", "", "Bio")
	profileUpdateCmd.Flags().String("avatar", "", "Avatar URL")
	profileUpdateCmd.Flags().String("banner", "", "Banner URL")
	profileUpdateCmd.Flags().String("locale", "", "Language (zh/en/ja/ko...)")
	profileCmd.AddCommand(profileShowCmd, profileUpdateCmd, profileSetAvatarCmd)
	rootCmd.AddCommand(profileCmd)
}
