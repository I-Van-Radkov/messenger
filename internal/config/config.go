package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP *HTTPConfig
	DB   *DBConfig
	Auth *AuthConfig
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	User     string
	Name     string
	Password string
	Host     string
	Port     int
}

type AuthConfig struct {
	JwtSecret    string
	JwtExpiresIn time.Duration
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	return cfg
}

func Load() (*Config, error) {
	var cfg Config
	var err error

	// HTTP Config
	if cfg.HTTP.Port, err = getIntEnv("API_PORT", 8080); err != nil {
		return nil, fmt.Errorf("http port: %w", err)
	}
	if cfg.HTTP.ReadTimeout, err = getDurationEnv("HTTP_READ_TIMEOUT", 15*time.Second); err != nil {
		return nil, fmt.Errorf("read timeout: %w", err)
	}
	if cfg.HTTP.WriteTimeout, err = getDurationEnv("HTTP_WRITE_TIMEOUT", 15*time.Second); err != nil {
		return nil, fmt.Errorf("write timeout: %w", err)
	}

	// DB Config
	if cfg.DB.User, err = getStringEnv("DB_USER"); err != nil {
		return nil, fmt.Errorf("db user: %w", err)
	}
	if cfg.DB.Name, err = getStringEnv("DB_NAME"); err != nil {
		return nil, fmt.Errorf("db name: %w", err)
	}
	if cfg.DB.Password, err = getStringEnv("DB_PASSWORD"); err != nil {
		return nil, fmt.Errorf("db password: %w", err)
	}
	if cfg.DB.Host, err = getStringEnv("DB_HOST"); err != nil {
		return nil, fmt.Errorf("db host: %w", err)
	}
	if cfg.DB.Port, err = getIntEnv("DB_PORT", 5432); err != nil {
		return nil, fmt.Errorf("db port: %w", err)
	}

	// Auth Config
	if cfg.Auth.JwtSecret, err = getStringEnv("JWT_SECRET"); err != nil {
		return nil, fmt.Errorf("jwt secret: %w", err)
	}
	if cfg.Auth.JwtExpiresIn, err = getDurationEnv("JWT_EXPIRES_IN", 24*time.Hour); err != nil {
		return nil, fmt.Errorf("jwt expires in: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &cfg, nil
}

func getStringEnv(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("env %s is required", name)
	}
	return val, nil
}

func getIntEnv(name string, defaultValue int) (int, error) {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(val)
}

func getDurationEnv(name string, defaultValue time.Duration) (time.Duration, error) {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue, nil
	}
	return time.ParseDuration(val)
}

func (c *Config) validate() error {
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("http port out of range")
	}
	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		return fmt.Errorf("db port out of range")
	}
	if c.HTTP.ReadTimeout < 0 {
		return fmt.Errorf("read timeout must be non-negative")
	}
	if c.HTTP.WriteTimeout < 0 {
		return fmt.Errorf("write timeout must be non-negative")
	}
	return nil
}
