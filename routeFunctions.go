package main

import (
	_ "embed"
	"html/template"
	"net/http"
	"os/exec"
	"time"
)

//go:embed index.tmpl.html
var indexHTML string

var indexTmpl = template.Must(template.New("index").Parse(indexHTML))

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func IndexHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTmpl.Execute(w, cfg); err != nil {
			http.Error(w, "template error", 500)
		}
	}
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	appid := r.PathValue("appid")
	stopPs := exec.Command("Stop-Process", "-Name", appid)
	stopPs.Run()

	time.Sleep(2 * time.Second)

	startPs := exec.Command("Start-Process", "-FilePath", appid)
	startPs.Run()
	w.Write([]byte("Restart Complete!"))
}
