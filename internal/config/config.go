package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	GRPCPort string
	DSN      string
}

func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		DSN:      getEnv("CONCERTO_DSN", ""),
	}
	if cfg.DSN == "" {
		log.Fatal("CONCERTO_DSN environment variable not set")
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
