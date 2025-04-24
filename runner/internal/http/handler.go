package http

import (
	"github.com/xoesae/judge/runner/internal/config"
	"github.com/xoesae/judge/runner/internal/filesystem"
	"github.com/xoesae/judge/runner/internal/jail"
	"net/http"
	"path/filepath"
)

func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("script")
	if err != nil {
		http.Error(w, "error on upload file: "+err.Error(), 400)
		return
	}
	defer file.Close()

	cfg := config.GetConfig()

	destinationPath := filepath.Join(cfg.RootFs, "tmp", "scripts")

	_, err = filesystem.SaveScriptFile(file, header.Filename, destinationPath)
	if err != nil {
		http.Error(w, "error on save file: "+err.Error(), 500)
		return
	}

	result, err := jail.InitChildProcess(header.Filename)

	if err != nil {
		http.Error(w, "error on init child: "+err.Error(), 500)
		return
	}

	w.Write([]byte("STDOUT:\n"))
	w.Write(result.Output)
	w.Write([]byte("\nSTDERR:\n"))
	w.Write(result.Error)
}
