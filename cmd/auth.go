package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "帐号认证 (注册/登录/退出)",
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "创建新帐号",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if username == "" || email == "" || password == "" {
			return fmt.Errorf("--username, --email, --password are all required")
		}

		client := api.NewClient(config.GetAPIURL(), "")
		resp, err := client.POST("/web/auth/register", map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		})
		if err != nil {
			return fmt.Errorf("register failed: %w", err)
		}

		// Save credentials
		var data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			User         struct {
				ID          string `json:"id"`
				Username    string `json:"username"`
				DisplayName string `json:"display_name"`
			} `json:"user"`
		}
		json.Unmarshal(resp.Data, &data)

		config.SaveCredentials(&config.Credentials{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
			UserID:       data.User.ID,
			Username:     data.User.Username,
			DisplayName:  data.User.DisplayName,
		})

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Printf("✅ 注册成功! 用户: %s (@%s)\n", data.User.DisplayName, data.User.Username)
			fmt.Printf("   Token 已保存到 ~/.opendiscuz/credentials.json\n")
		}
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录已有帐号",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if email == "" || password == "" {
			return fmt.Errorf("--email and --password are required")
		}

		client := api.NewClient(config.GetAPIURL(), "")
		resp, err := client.POST("/web/auth/login", map[string]string{
			"email":    email,
			"password": password,
		})
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		var data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			User         struct {
				ID          string `json:"id"`
				Username    string `json:"username"`
				DisplayName string `json:"display_name"`
				UserType    string `json:"user_type"`
				Verified    bool   `json:"verified"`
			} `json:"user"`
		}
		json.Unmarshal(resp.Data, &data)

		config.SaveCredentials(&config.Credentials{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
			UserID:       data.User.ID,
			Username:     data.User.Username,
			DisplayName:  data.User.DisplayName,
		})

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Printf("✅ 登录成功! 欢迎 %s (@%s)\n", data.User.DisplayName, data.User.Username)
		}
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.ClearCredentials()
		if jsonOutput {
			fmt.Println(`{"message":"logged out"}`)
		} else {
			fmt.Println("✅ 已退出登录")
		}
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "查看当前登录信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds := config.LoadCredentials()
		if creds == nil {
			return fmt.Errorf("未登录。请先运行 'opendiscuz auth login' 或设置 OPENDISCUZ_TOKEN")
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(creds, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf("用户: %s (@%s)\n", creds.DisplayName, creds.Username)
			fmt.Printf("ID:   %s\n", creds.UserID)
			fmt.Printf("API:  %s\n", config.GetAPIURL())
		}
		return nil
	},
}

func init() {
	registerCmd.Flags().String("username", "", "用户名 (必填)")
	registerCmd.Flags().String("email", "", "邮箱 (必填)")
	registerCmd.Flags().String("password", "", "密码 (必填, 最少8位)")

	loginCmd.Flags().String("email", "", "邮箱 (必填)")
	loginCmd.Flags().String("password", "", "密码 (必填)")

	authCmd.AddCommand(registerCmd, loginCmd, logoutCmd, whoamiCmd)
	rootCmd.AddCommand(authCmd)
}
