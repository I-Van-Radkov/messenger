package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
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

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	return cfg
}

func Load() (*Config, error) {
	cfg := &Config{}

	writeTimeoutDur, err := getDurationEnv("HTTP_WRITE_TIMEOUT")
	if err != nil {
		return nil, fmt.Errorf("write timeout: %w", err)
	}

	readTimeoutDur, err := getDurationEnv("HTTP_READ_TIMEOUT")
	if err != nil {
		return nil, fmt.Errorf("read timeout: %w", err)
	}

	httpPort, err := getIntEnv("API_PORT")
	if err != nil {
		return nil, fmt.Errorf("http port: %w", err)
	}

	cfg.HTTP = HTTPConfig{
		Port:         httpPort,
		WriteTimeout: writeTimeoutDur,
		ReadTimeout:  readTimeoutDur,
	}

	dbPort, err := getIntEnv("DB_PORT")
	if err != nil {
		return nil, fmt.Errorf("db port: %w", err)
	}

	cfg.DB = DBConfig{
		User:     getRequiredEnv("DB_USER"),
		Name:     getRequiredEnv("DB_NAME"),
		Password: getRequiredEnv("DB_PASSWORD"),
		Host:     getRequiredEnv("HOST"),
		Port:     dbPort,
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return cfg, nil
}

func getDurationEnv(name string) (time.Duration, error) {
	val := os.Getenv(name)
	if val == "" {
		return 0, fmt.Errorf("env %v is required", name)
	}

	return time.ParseDuration(val)
}

func getRequiredEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		panic(fmt.Sprintf("env %s is required", name))
	}

	return val
}

func getIntEnv(name string) (int, error) {
	val := os.Getenv(name)
	if val == "" {
		return 0, fmt.Errorf("env %v is required", name)
	}

	return strconv.Atoi(val)
}

func (c *Config) validate() error {
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("http port out of range")
	}

	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		return fmt.Errorf("db port out of range")
	}

	if c.HTTP.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if c.HTTP.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	return nil
}
