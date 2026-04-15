package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"fyne.io/systray"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Appid    string    `yaml:"appid"`
	FilePath string    `yaml:"filepath"`
	WaitTime int       `yaml:"waittime"`
	Command  *[]string `yaml:"command"`
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
	logFile, err := os.OpenFile("remote-restarter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger := slog.New(slog.NewTextHandler(logFile, nil))
		slog.SetDefault(logger)
	}

	slog.Info("Application starting")

	mux := http.NewServeMux()

	data, err := os.ReadFile("RestartConfigs.yml")
	if err != nil {
		slog.Error("reading config", "error", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		slog.Error("parsing yaml", "error", err)
	}

	mux.HandleFunc("GET /", IndexHandler(&cfg))
	mux.HandleFunc("GET /health", HealthHandler)
	mux.HandleFunc("GET /restart/{appid}", RestartHandler(&cfg))

	server := http.Server{
		Addr:    ":" + cfg.ServerConfig.Port,
		Handler: mux,
	}

	go func() {
		slog.Info("Server is starting", "port", cfg.ServerConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server Start Failed", "error", err)
		}
	}()

	systray.Run(func() { onReady(&cfg) }, onExit)
}

func onReady(cfg *Config) {
	systray.SetTitle("Restarter")
	systray.SetTooltip("Remote Restarter Server")

	mPort := systray.AddMenuItem(fmt.Sprintf("Server is running at localhost:%s", cfg.ServerConfig.Port), "")
	mPort.Disable()

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}
