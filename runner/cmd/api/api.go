package main

import (
	"github.com/xoesae/judge/runner/internal/config"
	"github.com/xoesae/judge/runner/internal/filesystem"
	"github.com/xoesae/judge/runner/internal/http"
	"github.com/xoesae/judge/runner/internal/jail"
	"os"
)

func main() {
	config.LoadConfig()

	// if is child process, exec jail
	if len(os.Args) > 2 && os.Args[1] == "child" {
		jail.RunIsolated(os.Args[2])
		return
	}

	// upload dir must exist
	filesystem.MakeUploadDir()

	http.InitServer()
}
