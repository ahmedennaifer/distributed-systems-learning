package main

import (
	//"fmt"
	"github.com/ahmedennaifer/taskq/internal"
	//	"github.com/ahmedennaifer/taskq/internal/publisher"
	"log"
	"net/http"
)

func main() {
	addr := ":8080"
	cache := internal.NewCache()
	server := internal.NewServer(addr, cache)
	http.HandleFunc("POST /api/v1/task", server.HandlePostTask)
	http.HandleFunc("GET /api/v1/tasks", server.HandleGetTasks)
	http.HandleFunc("GET /api/v1/task/{taskID}", server.HandleGetTaskByID)
	//if err := publisher.CreateTopics(); err != nil {
	//	fmt.Println("Cannot create topic:", err)
	//	return
	//}
	log.Fatal(http.ListenAndServe(":8081", nil))
}
