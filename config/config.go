package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Env           string
	Host          string
	Port          string
	DatabaseURI   string
	JWTSecret     string
	JWTExpiration time.Duration
}

var server *ServerConfig

// Load the env vars into ServerConfig
// If the env is development it uses godotenv to load otherwise
// the env variables are directly loaded
func Load() {
	env := getEnvOrDefault("SERVER_ENV", "development")
	log.Printf("Running on %s environment", env)

	// If no environment is set assume its dev environment
	if env == "development" {
		// Load .env file if dev environment
		if err := godotenv.Load(); err != nil {
			log.Println("error loading .env file")
			log.Fatal(err)
		}
	}
	server = &ServerConfig{
		Env:           env,
		Host:          getEnvOrDefault("HOST", "localhost"),
		Port:          getEnvOrDefault("PORT", "8080"),
		DatabaseURI:   getEnvOrDefault("DATABASE_URI", ""),
		JWTSecret:     getEnvOrDefault("JWT_SECRET", "secret:)"),
		JWTExpiration: time.Hour * 24, // TODO change this to something better
	}
}

// Wrapper around os.LookupEnv() that returns the default value if not the environment var is not set
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// Generic wrapper that checks if the config variable is nil
// If it is then it exits using log.Fatal otherwise returns the config variable
func getConfig[T any](config *T) *T {
	if config == nil {
		log.Fatal("Accessing config before loading")
	}
	return config
}

// Accessor for getting the server config exits the program if called before loading the config
// Otherwise it just returns the config
func Server() *ServerConfig {
	return getConfig(server)
}
