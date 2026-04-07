package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os/exec"
	"time"
)

//go:embed index.tmpl.html
var indexHTML string

// Compile the homepage
var indexTmpl = template.Must(template.New("index").Parse(indexHTML))

// Returns 200 OK if server is running
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Health Endpoint Requested")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// Return the Restartable Applications
func IndexHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Index Endpoint Requested")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTmpl.Execute(w, cfg); err != nil {
			http.Error(w, "template error", 500)
		}
	}
}

func RestartHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get appid from url
		appid := r.PathValue("appid")
		slog.Info("Restart Request Received", "appid", appid)

		// Find the app
		appcfg, err := cfg.GetApp(appid)
		if err != nil {
			errmsg := fmt.Sprintf("Cant find the app configuration for %s : %s", appid, err.Error())
			http.Error(w, errmsg, 500)
		}

		// Stop the app
		stopCmd := fmt.Sprintf("(Stop-Process -Name %v)", appcfg.Appid)
		stopPs := exec.Command("powershell", "-Command", stopCmd)
		if err := stopPs.Run(); err != nil {
			errmsg := fmt.Sprintf("Stop-Process Request Failed for %s :\n %s", appcfg.Appid, err.Error())
			http.Error(w, errmsg, 500)
		}

		// Wait for it to close
		time.Sleep(time.Duration(appcfg.WaitTime) * time.Second)

		// Start the app again
		startPs := exec.Command(appcfg.FilePath)
		if err := startPs.Start(); err != nil {
			errmsg := fmt.Sprintf("Starting Process Failed for %s :\n %s", appcfg.Appid, err.Error())
			http.Error(w, errmsg, 500)
		}

		// Application started checking its status
		time.Sleep(1 * time.Second)
		if startPs.ProcessState != nil && startPs.ProcessState.Exited() {
			http.Error(w, "App Exited After Starting", 500)
		}

		// Return OK
		if _, err = fmt.Fprintf(w, "App %s Restarted Successfully", appcfg.Appid); err != nil {
			slog.Error("Restart Successful but response could not be written", "error", err)
		}
	}
}
