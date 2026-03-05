package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var configDir string

func init() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".opendiscuz")
}

// Config holds CLI configuration
type Config struct {
	APIURL string `json:"api_url"`
}

// Credentials holds auth tokens
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	DisplayName  string `json:"display_name"`
}

// LoadConfig reads config from ~/.opendiscuz/config.json
func LoadConfig() *Config {
	cfg := &Config{APIURL: "http://localhost:3080"}
	data, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, cfg)
	return cfg
}

// SaveConfig writes config to ~/.opendiscuz/config.json
func SaveConfig(cfg *Config) error {
	os.MkdirAll(configDir, 0700)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(filepath.Join(configDir, "config.json"), data, 0600)
}

// LoadCredentials reads credentials from ~/.opendiscuz/credentials.json
func LoadCredentials() *Credentials {
	data, err := os.ReadFile(filepath.Join(configDir, "credentials.json"))
	if err != nil {
		return nil
	}
	creds := &Credentials{}
	json.Unmarshal(data, creds)
	return creds
}

// SaveCredentials writes credentials to ~/.opendiscuz/credentials.json
func SaveCredentials(creds *Credentials) error {
	os.MkdirAll(configDir, 0700)
	data, _ := json.MarshalIndent(creds, "", "  ")
	return os.WriteFile(filepath.Join(configDir, "credentials.json"), data, 0600)
}

// ClearCredentials removes stored credentials
func ClearCredentials() error {
	return os.Remove(filepath.Join(configDir, "credentials.json"))
}

// GetAPIURL returns API URL from env or config
func GetAPIURL() string {
	if url := os.Getenv("OPENDISCUZ_API_URL"); url != "" {
		return url
	}
	return LoadConfig().APIURL
}

// GetAccessToken returns token from env or stored credentials
func GetAccessToken() string {
	if token := os.Getenv("OPENDISCUZ_TOKEN"); token != "" {
		return token
	}
	creds := LoadCredentials()
	if creds != nil {
		return creds.AccessToken
	}
	return ""
}

// RequireAuth checks for auth and returns a helpful error if missing
func RequireAuth() error {
	if GetAccessToken() == "" {
		return fmt.Errorf("not authenticated. Run 'opendiscuz auth login' or set OPENDISCUZ_TOKEN env var")
	}
	return nil
}
