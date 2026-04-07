package main

import (
	"fmt"
	"net/http"
	"os"

	"fyne.io/systray"
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
func (c *Config) GetApp(name string) (*AppConfig, error) {
	for i := range c.Apps {
		if c.Apps[i].Appid == name {
			return &c.Apps[i], nil
		}
	}
	return nil, fmt.Errorf("App Not Found In configuration")
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
	mux.HandleFunc("GET /restart/{appid}", RestartHandler(&cfg))

	server := http.Server{
		Addr:    ":" + cfg.ServerConfig.Port,
		Handler: mux,
	}

	fmt.Printf("Server is Starting at %s\n", cfg.ServerConfig.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server Start Failed: %v\n", err)
		}
	}()

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Restarter")
	systray.SetTooltip("Remote Restarter Server")
	
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}
