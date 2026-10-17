package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type HTTP struct {
	Port int `mapstructure:"port"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type JWT struct {
	Secret string `mapstructure:"secret"`
}

type Broker struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type AppConfig struct {
	HTTP     HTTP     `mapstructure:"http"`
	Database Database `mapstructure:"database"`
	JWT      JWT      `mapstructure:"jwt"`
	Broker   Broker   `mapstructure:"broker"`
}

func LoadConfig(configFilePath string) (*AppConfig, error) {
	v := viper.New()
	v.SetConfigFile(configFilePath)

	// Defaults keep the application runnable when a config file omits a section.
	v.SetDefault("http.port", 3333)
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "pguser")
	v.SetDefault("database.password", "pgpassword")
	v.SetDefault("database.name", "dddtasks")
	v.SetDefault("jwt.secret", "")
	v.SetDefault("broker.host", "localhost")
	v.SetDefault("broker.port", 5672)
	v.SetDefault("broker.user", "guest")
	v.SetDefault("broker.password", "guest")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind the environment variable names already used by docker-compose and the
	// Makefile (DB_HOST, BROKER_HOST), plus their newer, prefixed equivalents.
	bindings := []struct {
		key  string
		envs []string
	}{
		{"http.port", []string{"HTTP_PORT"}},
		{"database.host", []string{"DB_HOST", "DATABASE_HOST"}},
		{"database.port", []string{"DB_PORT", "DATABASE_PORT"}},
		{"database.user", []string{"DB_USER", "DATABASE_USER"}},
		{"database.password", []string{"DB_PASSWORD", "DATABASE_PASSWORD"}},
		{"database.name", []string{"DB_NAME", "DATABASE_NAME"}},
		{"jwt.secret", []string{"JWT_SECRET"}},
		{"broker.host", []string{"BROKER_HOST", "RABBITMQ_HOST"}},
		{"broker.port", []string{"BROKER_PORT", "RABBITMQ_PORT"}},
		{"broker.user", []string{"BROKER_USER", "RABBITMQ_USER"}},
		{"broker.password", []string{"BROKER_PASSWORD", "RABBITMQ_PASSWORD"}},
	}

	for _, binding := range bindings {
		args := append([]string{binding.key}, binding.envs...)
		if err := v.BindEnv(args...); err != nil {
			return nil, fmt.Errorf("failed to bind environment variables for %s: %w", binding.key, err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
