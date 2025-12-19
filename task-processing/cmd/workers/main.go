package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/ahmedennaifer/taskq/internal"
	"github.com/ahmedennaifer/taskq/internal/workers"
)

func getRandomPort() int {
	rand.Seed(time.Now().UnixNano())
	excludedPorts := map[int]bool{8000: true, 8080: true, 8081: true}

	for {
		port := 8000 + rand.Intn(1000)
		if !excludedPorts[port] {
			return port
		}
	}
}

func main() {
	logger := internal.NewLogger("worker.log")
	fmt.Println("Starting worker application...")
	logger.Info("starting worker application")

	port := getRandomPort()
	addr := fmt.Sprintf(":%d", port)

	w, err := workers.NewWorker(addr, logger)
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
