package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Server struct {
	addr      string
	taskCache *Cache
}

func NewServer(addr string, cache *Cache) *Server {
	return &Server{
		addr:      addr,
		taskCache: cache,
	}
}

func (s *Server) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.taskCache.Db); err != nil {
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}
}

func (s *Server) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")
	if taskID == "" {
		http.Error(w, fmt.Sprintf("error: taskID cannot be nul\n"), 401)
		return
	}
	parsed, err := uuid.Parse(taskID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: malformed task id: %v\n", taskID), 401)
		return
	}
	task, err := s.taskCache.Get(parsed)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v\n", err), 401)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, fmt.Sprintf("error: cannot encode response payload: %v", err), 500)
		return
	}
}

func (s *Server) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	task := NewTask()
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}
	if err := task.ValidateFields(); err != nil {
		http.Error(w, fmt.Sprintf("%v\n", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := s.taskCache.Add(*task); err != nil {
		http.Error(w, fmt.Sprintf("%v\n", err), 500)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]string{"task added with success": task.ID.String()}); err != nil {
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}
}
