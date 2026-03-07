package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	IMAP struct {
		Host    string `mapstructure:"host"`
		User    string `mapstructure:"user"`
		Pass    string `mapstructure:"pass"`
		Subject string `mapstructure:"subject"`
	} `mapstructure:"imap"`

	Auth struct {
		SPAPublicKey        string `mapstructure:"spa_public_key"`
		VerifierPrivateKey string `mapstructure:"verifier_private_key"`
	} `mapstructure:"auth"`

	App struct {
		SPACallbackURL string        `mapstructure:"spa_callback_url"`
		PollInterval   time.Duration `mapstructure:"poll_interval"`
	} `mapstructure:"app"`

	Puller struct {
		StartDate   string `mapstructure:"start_date"`
		StateFile   string `mapstructure:"state_file"`
		MCPCommand  string `mapstructure:"mcp_command"`
		MCPModel    string `mapstructure:"mcp_model"`
		TempDir     string `mapstructure:"temp_dir"`
	} `mapstructure:"puller"`
}

func Load() (*Config, error) {
	// Load .env if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file found but failed to load: %v", err)
		}
	}

	v := viper.New()

	// Configuration file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Set defaults
	v.SetDefault("app.poll_interval", "30s")

	// Environment variable settings
	v.SetEnvPrefix("MAILADK")
	v.AutomaticEnv()

	// Handle nested keys for environment variables
	// e.g., MAILADK_IMAP_HOST -> imap.host
	v.BindEnv("imap.host", "MAILADK_IMAP_HOST")
	v.BindEnv("imap.user", "MAILADK_IMAP_USER")
	v.BindEnv("imap.pass", "MAILADK_IMAP_PASS")
	v.BindEnv("imap.subject", "MAILADK_IMAP_SUBJECT")
	v.BindEnv("auth.spa_public_key", "MAILADK_SPA_PUBLIC_KEY")
	v.BindEnv("auth.verifier_private_key", "MAILADK_VERIFIER_PRIVATE_KEY")
	v.BindEnv("app.spa_callback_url", "MAILADK_SPA_CALLBACK_URL")
	v.BindEnv("app.poll_interval", "MAILADK_POLL_INTERVAL")
	v.BindEnv("puller.start_date", "MAILADK_PULLER_START_DATE")
	v.BindEnv("puller.state_file", "MAILADK_PULLER_STATE_FILE")
	v.BindEnv("puller.mcp_command", "MAILADK_PULLER_MCP_COMMAND")
	v.BindEnv("puller.mcp_model", "MAILADK_PULLER_MCP_MODEL")
	v.BindEnv("puller.temp_dir", "MAILADK_PULLER_TEMP_DIR")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found is fine, we'll fall back to env vars
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
