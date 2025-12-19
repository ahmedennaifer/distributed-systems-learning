package main

import (
	"fmt"
	"net/http"

	"github.com/ahmedennaifer/taskq/internal"
	"github.com/ahmedennaifer/taskq/internal/workers"
)

func main() {
	logger := internal.NewLogger("worker.log")
	fmt.Println("Starting worker application...")
	logger.Info("starting worker application")

	w, err := workers.NewWorker(":8001", logger)
	if err != nil {
		fmt.Printf("ERROR: Failed to create worker: %v\n", err)
		logger.Error("failed to create worker", "error", err)
		return
	}

	fmt.Printf("Worker created: ID=%s, Addr=%s\n", w.ID, w.Addr)

	err = w.Register("http://localhost:8081/api/v1")
	if err != nil {
		fmt.Printf("ERROR: Failed to register worker: %v\n", err)
		logger.Error("failed to register worker", "error", err)
		return
	}

	fmt.Printf("Worker registered successfully: %s\n", w.ID)
	logger.Info("worker ready to accept requests", "workerID", w.ID, "addr", w.Addr)

	http.HandleFunc("GET /health", w.HandleHealth)

	fmt.Printf("Starting HTTP server on %s\n", w.Addr)
	logger.Info("starting http server", "addr", w.Addr)
	if err := http.ListenAndServe(w.Addr, nil); err != nil {
		fmt.Printf("ERROR: HTTP server failed: %v\n", err)
		logger.Error("http server failed", "error", err)
	}
}
