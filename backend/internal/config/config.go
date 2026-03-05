package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
	Auth       AuthConfig       `mapstructure:"auth"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpiryHours int    `mapstructure:"expiry_hours"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

type AuthConfig struct {
	Username     string `mapstructure:"username"`
	PasswordHash string `mapstructure:"password_hash"`
}

func (c *Config) Validate() error {
	if len(c.Encryption.Key) != 32 {
		return fmt.Errorf("encryption key must be exactly 32 bytes, got %d", len(c.Encryption.Key))
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters, got %d", len(c.JWT.Secret))
	}
	if c.Auth.Username == "" {
		return fmt.Errorf("auth username cannot be empty")
	}
	if c.Auth.PasswordHash == "" {
		return fmt.Errorf("auth password hash cannot be empty")
	}
	return nil
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("KM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
