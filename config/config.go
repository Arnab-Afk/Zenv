// config/config.go
package config

type Config struct {
	HSMEndpoint      string
	DatabaseURL      string
	RedisURL         string
	RateLimit        int
	KeyRotationHours int
}

func Load() *Config {
	// Load from environment variables
	return &Config{
		HSMEndpoint:      getEnv("HSM_ENDPOINT", "localhost:1234"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://localhost:5432/secrets"),
		RedisURL:         getEnv("REDIS_URL", "localhost:6379"),
		RateLimit:        getEnvAsInt("RATE_LIMIT", 100),
		KeyRotationHours: getEnvAsInt("KEY_ROTATION_HOURS", 24),
	}
}
