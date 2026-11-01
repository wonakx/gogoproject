// Command apidocs serves the API contract without connecting to a database.
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"example.com/gogoproject/internal/platform/httpapi"
)

func main() {
	addr := os.Getenv("DOCS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	server := &http.Server{Addr: addr, Handler: httpapi.DocumentationHandler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	slog.Info("API documentation server starting", "address", addr, "path", "/docs/")
	if err := server.ListenAndServe(); err != nil {
		slog.Error("documentation server failed", "error", err)
		os.Exit(1)
	}
}
