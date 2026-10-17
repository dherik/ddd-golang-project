package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testConfig = `
http:
  port: 3333
database:
  host: "db"
  port: 5432
  user: "db-user"
  password: "db-password"
  name: "db-name"
jwt:
  secret: "file-secret"
broker:
  host: "broker"
  port: 5672
  user: "broker-user"
  password: "broker-password"
`

func writeTestConfig(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(testConfig), 0o600))
	return path
}

func clearConfigEnv(t *testing.T) {
	t.Helper()

	for _, key := range []string{
		"HTTP_PORT",
		"DB_HOST", "DATABASE_HOST",
		"DB_PORT", "DATABASE_PORT",
		"DB_USER", "DATABASE_USER",
		"DB_PASSWORD", "DATABASE_PASSWORD",
		"DB_NAME", "DATABASE_NAME",
		"JWT_SECRET",
		"BROKER_HOST", "RABBITMQ_HOST",
		"BROKER_PORT", "RABBITMQ_PORT",
		"BROKER_USER", "RABBITMQ_USER",
		"BROKER_PASSWORD", "RABBITMQ_PASSWORD",
	} {
		t.Setenv(key, "")
	}
}

func TestLoadConfig(t *testing.T) {
	clearConfigEnv(t)

	cfg, err := LoadConfig(writeTestConfig(t))
	require.NoError(t, err)

	assert.Equal(t, 3333, cfg.HTTP.Port)
	assert.Equal(t, "db", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "db-user", cfg.Database.User)
	assert.Equal(t, "db-password", cfg.Database.Password)
	assert.Equal(t, "db-name", cfg.Database.Name)
	assert.Equal(t, "file-secret", cfg.JWT.Secret)
	assert.Equal(t, "broker", cfg.Broker.Host)
	assert.Equal(t, 5672, cfg.Broker.Port)
	assert.Equal(t, "broker-user", cfg.Broker.User)
	assert.Equal(t, "broker-password", cfg.Broker.Password)
}

func TestLoadConfigEnvironmentOverrides(t *testing.T) {
	t.Setenv("HTTP_PORT", "4444")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DATABASE_PORT", "6543")
	t.Setenv("DB_USER", "env-user")
	t.Setenv("DB_PASSWORD", "env-password")
	t.Setenv("DB_NAME", "env-name")
	t.Setenv("JWT_SECRET", "env-secret")
	t.Setenv("BROKER_HOST", "rabbitmq")
	t.Setenv("RABBITMQ_PORT", "5673")
	t.Setenv("BROKER_USER", "env-broker-user")
	t.Setenv("BROKER_PASSWORD", "env-broker-password")

	cfg, err := LoadConfig(writeTestConfig(t))
	require.NoError(t, err)

	assert.Equal(t, 4444, cfg.HTTP.Port)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, 6543, cfg.Database.Port)
	assert.Equal(t, "env-user", cfg.Database.User)
	assert.Equal(t, "env-password", cfg.Database.Password)
	assert.Equal(t, "env-name", cfg.Database.Name)
	assert.Equal(t, "env-secret", cfg.JWT.Secret)
	assert.Equal(t, "rabbitmq", cfg.Broker.Host)
	assert.Equal(t, 5673, cfg.Broker.Port)
	assert.Equal(t, "env-broker-user", cfg.Broker.User)
	assert.Equal(t, "env-broker-password", cfg.Broker.Password)
}
