package main

import (
	"fmt"
	"net/http"

	"github.com/ahmedennaifer/taskq/internal"
)

func main() {
	addr := ":8081"
	fmt.Println("Starting server application...")

	cache := internal.NewCache()
	server, err := internal.NewServer(addr, cache)
	if err != nil {
		fmt.Printf("ERROR: Failed to create server: %v\n", err)
		return
	}

	fmt.Println("Server initialized successfully")
	fmt.Println("Registering HTTP handlers...")

	http.HandleFunc("POST /api/v1/task", server.HandlePostTask)
	http.HandleFunc("GET /api/v1/tasks", server.HandleGetTasks)
	http.HandleFunc("GET /api/v1/task/{taskID}", server.HandleGetTaskByID)
	http.HandleFunc("GET /api/v1/workers", server.HandleListWorkers)
	http.HandleFunc("POST /api/v1/worker", server.HandleRegisterWorker)

	fmt.Printf("Starting HTTP server on %s\n", server.Addr)
	fmt.Println("Server is ready to accept requests")

	server.RunHealthCheck()
	if err := http.ListenAndServe(server.Addr, nil); err != nil {
		fmt.Printf("ERROR: HTTP server failed: %v\n", err)
	}
}
