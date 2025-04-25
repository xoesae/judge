package http

import (
	"encoding/json"
	"github.com/xoesae/judge/runner/internal/config"
	"github.com/xoesae/judge/runner/internal/filesystem"
	"github.com/xoesae/judge/runner/internal/jail"
	"net/http"
	"path/filepath"
)

type ProcessResultResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

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

	response := ProcessResultResponse{
		Output: string(result.Output),
		Error:  string(result.Error),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
