package config

import "os"

type AppConfig struct {
	Port        string
	GoogleClientID string
}

func Load() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return AppConfig{
		Port:           port,
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	}
}
