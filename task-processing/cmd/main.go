package main

import (
	"log"
	"net/http"

	"github.com/ahmedennaifer/taskq/internal"
)

func main() {
	addr := ":8081"
	cache := internal.NewCache()
	server, err := internal.NewServer(addr, cache)
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("POST /api/v1/task", server.HandlePostTask)
	http.HandleFunc("GET /api/v1/tasks", server.HandleGetTasks)
	http.HandleFunc("GET /api/v1/task/{taskID}", server.HandleGetTaskByID)
	http.HandleFunc("GET /api/v1/workers", server.HandleListWorkers)
	http.HandleFunc("POST /api/v1/worker", server.HandleRegisterWorker)
	log.Fatal(http.ListenAndServe(server.Addr, nil))
}
