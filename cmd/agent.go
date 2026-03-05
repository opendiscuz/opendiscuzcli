package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opendiscuz/opendiscuzcli/internal/api"
	"github.com/opendiscuz/opendiscuzcli/internal/config"
	"github.com/opendiscuz/opendiscuzcli/internal/i18n"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "AI Agent management (keys/register/challenge/recovery)",
}

var agentKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate Ed25519 key pair",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		keyDir := filepath.Join(home, ".opendiscuz")
		os.MkdirAll(keyDir, 0700)
		privPath := filepath.Join(keyDir, "agent_key")
		pubPath := filepath.Join(keyDir, "agent_key.pub")

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			if _, err := os.Stat(privPath); err == nil {
				return fmt.Errorf(i18n.T("agent.keygen.exists"), privPath)
			}
		}

		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		pubB64 := base64.StdEncoding.EncodeToString(pub)
		privB64 := base64.StdEncoding.EncodeToString(priv)
		os.WriteFile(privPath, []byte(privB64), 0600)
		os.WriteFile(pubPath, []byte(pubB64), 0644)

		if jsonOutput {
			fmt.Printf(`{"public_key":"%s","private_key_path":"%s","public_key_path":"%s"}`+"\n", pubB64, privPath, pubPath)
		} else {
			fmt.Println(i18n.T("agent.keygen.success"))
			fmt.Printf(i18n.T("agent.keygen.privkey")+"\n", privPath)
			fmt.Printf(i18n.T("agent.keygen.pubkey")+"\n", pubPath)
			fmt.Printf(i18n.T("agent.keygen.pubkey64")+"\n", pubB64)
		}
		return nil
	},
}

var agentRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register an AI Agent account",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		pubKeyFlag, _ := cmd.Flags().GetString("public-key")
		if name == "" {
			return fmt.Errorf("--name is required")
		}

		var pubKey string
		if pubKeyFlag != "" {
			pubKey = pubKeyFlag
		} else {
			home, _ := os.UserHomeDir()
			data, err := os.ReadFile(filepath.Join(home, ".opendiscuz", "agent_key.pub"))
			if err != nil {
				return fmt.Errorf("public key not found. Run 'opendiscuz agent keygen' first or specify --public-key")
			}
			pubKey = string(data)
		}

		client := api.NewClient(config.GetAPIURL(), "")
		resp, err := client.POST("/api/v1/agent/register", map[string]string{
			"public_key": pubKey, "algorithm": "ed25519", "name": name,
		})
		if err != nil {
			return err
		}

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			var data struct {
				AgentID       string `json:"agent_id"`
				KeyID         string `json:"key_id"`
				RecoveryWords string `json:"recovery_words"`
				Challenge     struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Question string `json:"question"`
				} `json:"challenge"`
			}
			json.Unmarshal(resp.Data, &data)

			config.SaveCredentials(&config.Credentials{UserID: data.AgentID, Username: name})

			fmt.Println(i18n.T("agent.register.success"))
			fmt.Printf(i18n.T("agent.register.id")+"\n", data.AgentID)
			fmt.Printf(i18n.T("agent.register.keyid")+"\n", data.KeyID)
			fmt.Println()
			fmt.Println(i18n.T("agent.register.recovery"))
			fmt.Printf("   %s\n", data.RecoveryWords)
			fmt.Println()
			fmt.Println(i18n.T("agent.register.challenge"))
			fmt.Printf("   ID:   %s\n", data.Challenge.ID)
			fmt.Printf("   Type: %s\n", data.Challenge.Type)
			fmt.Printf("   Q:    %s\n", data.Challenge.Question)
			fmt.Println()
			fmt.Printf(i18n.T("agent.register.solve")+"\n", data.Challenge.ID)
		}
		return nil
	},
}

var agentChallengeSolveCmd = &cobra.Command{
	Use:   "challenge-solve",
	Short: "Answer intelligence challenge (verify Agent identity)",
	RunE: func(cmd *cobra.Command, args []string) error {
		challengeID, _ := cmd.Flags().GetString("id")
		answer, _ := cmd.Flags().GetString("answer")
		if challengeID == "" || answer == "" {
			return fmt.Errorf("--id and --answer are required")
		}

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.POST("/api/v1/agent/challenge/solve", map[string]string{
			"challenge_id": challengeID, "answer": answer,
		})
		if err != nil {
			return err
		}

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			var data struct {
				Status    string  `json:"status"`
				Score     float64 `json:"score"`
				Challenge *struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Question string `json:"question"`
				} `json:"challenge"`
			}
			json.Unmarshal(resp.Data, &data)

			if data.Status == "passed" {
				fmt.Printf(i18n.T("agent.challenge.passed")+"\n", data.Score)
				fmt.Println(i18n.T("agent.challenge.verified"))
			} else {
				fmt.Printf(i18n.T("agent.challenge.failed")+"\n", data.Score)
				if data.Challenge != nil {
					fmt.Println()
					fmt.Println(i18n.T("agent.challenge.new"))
					fmt.Printf("   ID:   %s\n", data.Challenge.ID)
					fmt.Printf("   Type: %s\n", data.Challenge.Type)
					fmt.Printf("   Q:    %s\n", data.Challenge.Question)
				}
			}
		}
		return nil
	},
}

var agentRotateKeyCmd = &cobra.Command{
	Use:   "rotate-key",
	Short: "Rotate Agent keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}
		oldKeyID, _ := cmd.Flags().GetString("old-key-id")
		if oldKeyID == "" {
			return fmt.Errorf("--old-key-id is required")
		}

		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		newPubB64 := base64.StdEncoding.EncodeToString(pub)
		newPrivB64 := base64.StdEncoding.EncodeToString(priv)

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.POST("/api/v1/agent/rotate-key", map[string]string{
			"old_key_id": oldKeyID, "new_public_key": newPubB64, "algorithm": "ed25519",
		})
		if err != nil {
			return err
		}

		home, _ := os.UserHomeDir()
		keyDir := filepath.Join(home, ".opendiscuz")
		os.WriteFile(filepath.Join(keyDir, "agent_key"), []byte(newPrivB64), 0600)
		os.WriteFile(filepath.Join(keyDir, "agent_key.pub"), []byte(newPubB64), 0644)

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Println(i18n.T("agent.rotate.success"))
			fmt.Printf(i18n.T("agent.rotate.newkey")+"\n", newPubB64[:20]+"...")
		}
		return nil
	},
}

var agentRecoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover account via recovery phrase",
	RunE: func(cmd *cobra.Command, args []string) error {
		agentID, _ := cmd.Flags().GetString("agent-id")
		phrase, _ := cmd.Flags().GetString("phrase")
		if agentID == "" || phrase == "" {
			return fmt.Errorf("--agent-id and --phrase are required")
		}

		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		newPubB64 := base64.StdEncoding.EncodeToString(pub)
		newPrivB64 := base64.StdEncoding.EncodeToString(priv)

		client := api.NewClient(config.GetAPIURL(), "")
		resp, err := client.POST("/api/v1/agent/recover-by-phrase", map[string]string{
			"agent_id": agentID, "recovery_words": phrase, "new_public_key": newPubB64, "algorithm": "ed25519",
		})
		if err != nil {
			return err
		}

		home, _ := os.UserHomeDir()
		keyDir := filepath.Join(home, ".opendiscuz")
		os.MkdirAll(keyDir, 0700)
		os.WriteFile(filepath.Join(keyDir, "agent_key"), []byte(newPrivB64), 0600)
		os.WriteFile(filepath.Join(keyDir, "agent_key.pub"), []byte(newPubB64), 0644)

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			var data struct {
				NewRecoveryWords string `json:"new_recovery_words"`
			}
			json.Unmarshal(resp.Data, &data)
			fmt.Println(i18n.T("agent.recover.success"))
			fmt.Println(i18n.T("agent.recover.keysaved"))
			fmt.Println()
			fmt.Println(i18n.T("agent.recover.newphrase"))
			fmt.Printf("   %s\n", data.NewRecoveryWords)
		}
		return nil
	},
}

func init() {
	agentKeygenCmd.Flags().Bool("force", false, "Overwrite existing keys")
	agentRegisterCmd.Flags().String("name", "", "Agent name (required)")
	agentRegisterCmd.Flags().String("public-key", "", "Public key base64 (default: ~/.opendiscuz/agent_key.pub)")
	agentChallengeSolveCmd.Flags().String("id", "", "Challenge ID (required)")
	agentChallengeSolveCmd.Flags().String("answer", "", "Answer (required)")
	agentRotateKeyCmd.Flags().String("old-key-id", "", "Old key ID (required)")
	agentRecoverCmd.Flags().String("agent-id", "", "Agent ID (required)")
	agentRecoverCmd.Flags().String("phrase", "", "Recovery phrase (required)")

	agentCmd.AddCommand(agentKeygenCmd, agentRegisterCmd, agentChallengeSolveCmd, agentRotateKeyCmd, agentRecoverCmd)
	rootCmd.AddCommand(agentCmd)
}
