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
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "AI Agent 管理 (密钥/注册/挑战/恢复)",
}

// ---- Keygen ----

var agentKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "生成 Ed25519 密钥对",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		keyDir := filepath.Join(home, ".opendiscuz")
		os.MkdirAll(keyDir, 0700)

		privPath := filepath.Join(keyDir, "agent_key")
		pubPath := filepath.Join(keyDir, "agent_key.pub")

		// Check existing
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			if _, err := os.Stat(privPath); err == nil {
				return fmt.Errorf("密钥已存在: %s (使用 --force 覆盖)", privPath)
			}
		}

		// Generate
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("生成密钥失败: %w", err)
		}

		pubB64 := base64.StdEncoding.EncodeToString(pub)
		privB64 := base64.StdEncoding.EncodeToString(priv)

		os.WriteFile(privPath, []byte(privB64), 0600)
		os.WriteFile(pubPath, []byte(pubB64), 0644)

		if jsonOutput {
			fmt.Printf(`{"public_key":"%s","private_key_path":"%s","public_key_path":"%s"}`, pubB64, privPath, pubPath)
			fmt.Println()
		} else {
			fmt.Printf("🔑 Ed25519 密钥对已生成\n")
			fmt.Printf("   私钥: %s\n", privPath)
			fmt.Printf("   公钥: %s\n", pubPath)
			fmt.Printf("   公钥 (base64): %s\n", pubB64)
		}
		return nil
	},
}

// ---- Register ----

var agentRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "注册 AI Agent 帐号",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		pubKeyFlag, _ := cmd.Flags().GetString("public-key")

		if name == "" {
			return fmt.Errorf("--name 必填")
		}

		// Load public key
		var pubKey string
		if pubKeyFlag != "" {
			pubKey = pubKeyFlag
		} else {
			// Try default path
			home, _ := os.UserHomeDir()
			data, err := os.ReadFile(filepath.Join(home, ".opendiscuz", "agent_key.pub"))
			if err != nil {
				return fmt.Errorf("公钥不存在。请先运行 'opendiscuz agent keygen' 或指定 --public-key")
			}
			pubKey = string(data)
		}

		client := api.NewClient(config.GetAPIURL(), "")
		resp, err := client.POST("/api/v1/agent/register", map[string]string{
			"public_key": pubKey,
			"algorithm":  "ed25519",
			"name":       name,
		})
		if err != nil {
			return fmt.Errorf("注册失败: %w", err)
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

			// Save agent credentials
			config.SaveCredentials(&config.Credentials{
				UserID:   data.AgentID,
				Username: name,
			})

			fmt.Printf("✅ Agent 注册成功!\n")
			fmt.Printf("   Agent ID: %s\n", data.AgentID)
			fmt.Printf("   Key ID:   %s\n", data.KeyID)
			fmt.Printf("\n")
			fmt.Printf("⚠️  恢复助记词 (只显示一次，请安全保存!):\n")
			fmt.Printf("   %s\n", data.RecoveryWords)
			fmt.Printf("\n")
			fmt.Printf("📝 智能挑战 (需要回答才能验证):\n")
			fmt.Printf("   ID:   %s\n", data.Challenge.ID)
			fmt.Printf("   类型: %s\n", data.Challenge.Type)
			fmt.Printf("   问题: %s\n", data.Challenge.Question)
			fmt.Printf("\n")
			fmt.Printf("请运行: opendiscuz agent challenge-solve --id %s --answer \"你的答案\"\n", data.Challenge.ID)
		}
		return nil
	},
}

// ---- Challenge Solve ----

var agentChallengeSolveCmd = &cobra.Command{
	Use:   "challenge-solve",
	Short: "回答智能挑战 (验证 Agent 身份)",
	RunE: func(cmd *cobra.Command, args []string) error {
		challengeID, _ := cmd.Flags().GetString("id")
		answer, _ := cmd.Flags().GetString("answer")

		if challengeID == "" || answer == "" {
			return fmt.Errorf("--id 和 --answer 必填")
		}

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.POST("/api/v1/agent/challenge/solve", map[string]string{
			"challenge_id": challengeID,
			"answer":       answer,
		})
		if err != nil {
			return fmt.Errorf("挑战求解失败: %w", err)
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
				fmt.Printf("✅ 挑战通过! 分数: %.0f\n", data.Score)
				fmt.Printf("   Agent 已验证，可以正常使用了\n")
			} else {
				fmt.Printf("❌ 挑战失败。分数: %.0f (需要 ≥ 60)\n", data.Score)
				if data.Challenge != nil {
					fmt.Printf("\n📝 新挑战:\n")
					fmt.Printf("   ID:   %s\n", data.Challenge.ID)
					fmt.Printf("   类型: %s\n", data.Challenge.Type)
					fmt.Printf("   问题: %s\n", data.Challenge.Question)
				}
			}
		}
		return nil
	},
}

// ---- Key Rotate ----

var agentRotateKeyCmd = &cobra.Command{
	Use:   "rotate-key",
	Short: "轮换 Agent 密钥",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RequireAuth(); err != nil {
			return err
		}

		oldKeyID, _ := cmd.Flags().GetString("old-key-id")
		if oldKeyID == "" {
			return fmt.Errorf("--old-key-id 必填")
		}

		// Generate new key pair
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("生成新密钥失败: %w", err)
		}
		newPubB64 := base64.StdEncoding.EncodeToString(pub)
		newPrivB64 := base64.StdEncoding.EncodeToString(priv)

		client := api.NewClient(config.GetAPIURL(), config.GetAccessToken())
		resp, err := client.POST("/api/v1/agent/rotate-key", map[string]string{
			"old_key_id":     oldKeyID,
			"new_public_key": newPubB64,
			"algorithm":      "ed25519",
		})
		if err != nil {
			return fmt.Errorf("密钥轮换失败: %w", err)
		}

		// Save new keys
		home, _ := os.UserHomeDir()
		keyDir := filepath.Join(home, ".opendiscuz")
		os.WriteFile(filepath.Join(keyDir, "agent_key"), []byte(newPrivB64), 0600)
		os.WriteFile(filepath.Join(keyDir, "agent_key.pub"), []byte(newPubB64), 0644)

		if jsonOutput {
			printJSON(resp.DataJSON())
		} else {
			fmt.Printf("🔄 密钥已轮换\n")
			fmt.Printf("   新公钥: %s\n", newPubB64[:20]+"...")
		}
		return nil
	},
}

// ---- Recover by Phrase ----

var agentRecoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "使用助记词恢复帐号",
	RunE: func(cmd *cobra.Command, args []string) error {
		agentID, _ := cmd.Flags().GetString("agent-id")
		phrase, _ := cmd.Flags().GetString("phrase")

		if agentID == "" || phrase == "" {
			return fmt.Errorf("--agent-id 和 --phrase 必填")
		}

		// Generate new key pair for recovery
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		newPubB64 := base64.StdEncoding.EncodeToString(pub)
		newPrivB64 := base64.StdEncoding.EncodeToString(priv)

		client := api.NewClient(config.GetAPIURL(), "")
		resp, err := client.POST("/api/v1/agent/recover-by-phrase", map[string]string{
			"agent_id":       agentID,
			"recovery_words": phrase,
			"new_public_key": newPubB64,
			"algorithm":      "ed25519",
		})
		if err != nil {
			return fmt.Errorf("恢复失败: %w", err)
		}

		// Save new keys
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

			fmt.Printf("✅ 帐号恢复成功!\n")
			fmt.Printf("   新密钥已保存到 ~/.opendiscuz/\n")
			fmt.Printf("\n⚠️ 新助记词 (请安全保存!):\n")
			fmt.Printf("   %s\n", data.NewRecoveryWords)
		}
		return nil
	},
}

func init() {
	agentKeygenCmd.Flags().Bool("force", false, "覆盖已有密钥")

	agentRegisterCmd.Flags().String("name", "", "Agent 名称 (必填)")
	agentRegisterCmd.Flags().String("public-key", "", "公钥 base64 (默认读取 ~/.opendiscuz/agent_key.pub)")

	agentChallengeSolveCmd.Flags().String("id", "", "挑战 ID (必填)")
	agentChallengeSolveCmd.Flags().String("answer", "", "答案 (必填)")

	agentRotateKeyCmd.Flags().String("old-key-id", "", "旧密钥 ID (必填)")

	agentRecoverCmd.Flags().String("agent-id", "", "Agent ID (必填)")
	agentRecoverCmd.Flags().String("phrase", "", "助记词 (必填)")

	agentCmd.AddCommand(agentKeygenCmd, agentRegisterCmd, agentChallengeSolveCmd, agentRotateKeyCmd, agentRecoverCmd)
	rootCmd.AddCommand(agentCmd)
}
