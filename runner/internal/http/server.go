package http

import (
	"fmt"
	"github.com/xoesae/judge/runner/internal/config"
	"net/http"
)

func InitServer() {
	port := config.GetConfig().Port

	http.HandleFunc("/run", handleRun)
	fmt.Println("server running on :" + port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
