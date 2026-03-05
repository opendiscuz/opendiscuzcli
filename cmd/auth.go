package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/opendiscuz/opendiscuzcli/internal/i18n"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Account authentication (register/login/logout)",
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new account",
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
			fmt.Printf(i18n.T("auth.register.success")+"\n", data.User.DisplayName, data.User.Username)
			fmt.Println(i18n.T("auth.register.saved"))
		}
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to existing account",
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
			fmt.Printf(i18n.T("auth.login.success")+"\n", data.User.DisplayName, data.User.Username)
		}
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.ClearCredentials()
		if jsonOutput {
			fmt.Println(`{"message":"logged out"}`)
		} else {
			fmt.Println(i18n.T("auth.logout.success"))
		}
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "View current login info",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds := config.LoadCredentials()
		if creds == nil {
			return fmt.Errorf(i18n.T("auth.notlogged"))
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(creds, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf(i18n.T("auth.whoami.user")+"\n", creds.DisplayName, creds.Username)
			fmt.Printf(i18n.T("auth.whoami.id")+"\n", creds.UserID)
			fmt.Printf(i18n.T("auth.whoami.api")+"\n", config.GetAPIURL())
		}
		return nil
	},
}

func init() {
	registerCmd.Flags().String("username", "", "Username (required)")
	registerCmd.Flags().String("email", "", "Email (required)")
	registerCmd.Flags().String("password", "", "Password (required, min 8 chars)")

	loginCmd.Flags().String("email", "", "Email (required)")
	loginCmd.Flags().String("password", "", "Password (required)")

	authCmd.AddCommand(registerCmd, loginCmd, logoutCmd, whoamiCmd)
	rootCmd.AddCommand(authCmd)
}
