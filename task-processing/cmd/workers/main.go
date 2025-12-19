package main

import (
	"net/http"

	"github.com/ahmedennaifer/taskq/internal"
	"github.com/ahmedennaifer/taskq/internal/workers"
)

func main() {
	logger := internal.NewLogger()
	logger.Info("starting worker application")

	w, err := workers.NewWorker(":8001", logger)
	if err != nil {
		logger.Error("failed to create worker", "error", err)
		return
	}

	err = w.Register("http://localhost:8081/api/v1")
	if err != nil {
		logger.Error("failed to register worker", "error", err)
		return
	}

	logger.Info("worker ready to accept requests", "workerID", w.ID, "addr", w.Addr)

	http.HandleFunc("GET /health", w.HandleHealth)

	logger.Info("starting http server", "addr", w.Addr)
	if err := http.ListenAndServe(w.Addr, nil); err != nil {
		logger.Error("http server failed", "error", err)
	}
}
