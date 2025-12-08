package main

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Appid    string `yaml:"appid"`
	FilePath string `yaml:"filepath"`
	WaitTime int    `yaml:"waittime"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Apps         []AppConfig  `yaml:"apps"`
	ServerConfig ServerConfig `yaml:"serverconfig"`
}

// Helper to get app by name
func (c *Config) GetApp(name string) *AppConfig {
	for i := range c.Apps {
		if c.Apps[i].Appid == name {
			return &c.Apps[i]
		}
	}
	return nil
}

func main() {
	mux := http.NewServeMux()

	data, err := os.ReadFile("RestartConfigs.yml")
	if err != nil {
		fmt.Println("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println("parsing yaml: %w", err)
	}

	mux.HandleFunc("GET /", IndexHandler(&cfg))
	mux.HandleFunc("GET /health", HealthHandler)
	mux.HandleFunc("GET /restart/{appid}", restartHandler)

	server := http.Server{
		Addr:    ":" + cfg.ServerConfig.Port,
		Handler: mux,
	}
	server.ListenAndServe()

}
