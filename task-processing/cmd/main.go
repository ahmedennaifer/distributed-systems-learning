package main

import (
	"github.com/ahmedennaifer/taskq/internal"
	"log"
	"net/http"
)

func main() {
	addr := ":8080"
	cache := internal.NewCache()
	server, err := internal.NewServer(addr, cache)
	if err != nil {
		log.Fatal(err)
		return
	}
	http.HandleFunc("POST /api/v1/task", server.HandlePostTask)
	http.HandleFunc("GET /api/v1/tasks", server.HandleGetTasks)
	http.HandleFunc("GET /api/v1/task/{taskID}", server.HandleGetTaskByID)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
