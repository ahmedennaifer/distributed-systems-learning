package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ahmedennaifer/taskq/internal"
	"github.com/google/uuid"
)

var cache = NewCache()

type Cache struct {
	Db map[uuid.UUID]internal.Task
}

func NewCache() *Cache {
	return &Cache{
		Db: make(map[uuid.UUID]internal.Task, 0),
	}
}

func (c *Cache) Add(task internal.Task) error {
	for _, ctask := range c.Db {
		if ctask.Name == task.Name {
			fmt.Printf("error: cannot add task %v, it already exists\n", task.Name)
			return fmt.Errorf("error: cannot add task %v, it already exists", task.Name)
		}
	}
	c.Db[task.Id] = task
	return nil
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	task := internal.NewTask()
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		fmt.Printf("error decoding payload: %v\n", err)
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}
	if err := task.Validate(); err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, fmt.Sprintf("%v\n", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := cache.Add(*task); err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, fmt.Sprintf("%v\n", err), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		fmt.Printf("error encoding payload: %v\n", err)
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cache.Db); err != nil {
		fmt.Printf("error encoding payload: %v\n", err)
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}
}

func main() {
	http.HandleFunc("POST /api/v1/task", handlePost)
	http.HandleFunc("GET /api/v1/task", handleGet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
