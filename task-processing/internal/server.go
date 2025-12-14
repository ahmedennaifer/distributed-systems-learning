package internal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ahmedennaifer/taskq/internal/publisher"
	"github.com/google/uuid"
)

type Server struct {
	addr        string
	taskCache   *Cache
	kafkaClient *publisher.KafkaClient
	logger      *slog.Logger
}

func NewServer(addr string, cache *Cache) (*Server, error) {
	logger := NewLogger()
	logger.Info("initializing server", "addr", addr)

	pubClient, err := publisher.Init()
	if err != nil {
		logger.Error("failed to init kafkaClient", "error", err)
		fmt.Println("error: cannot init kafkaClient: ", err)
		return &Server{}, err
	}
	logger.Info("server initialized successfully")
	return &Server{
		addr:        addr,
		taskCache:   cache,
		kafkaClient: pubClient,
		logger:      logger,
	}, nil
}

func (s *Server) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("handling get tasks request")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.taskCache.Db); err != nil {
		s.logger.Error("failed to encode tasks", "error", err)
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}

	s.logger.Info("successfully returned all tasks", "count", len(s.taskCache.Db))
}

func (s *Server) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")
	s.logger.Debug("handling get task by ID request", "taskID", taskID)

	if taskID == "" {
		s.logger.Warn("taskID is empty")
		http.Error(w, fmt.Sprintf("error: taskID cannot be nul\n"), 400)
		return
	}
	parsed, err := uuid.Parse(taskID)
	if err != nil {
		s.logger.Warn("malformed task ID", "taskID", taskID, "error", err)
		http.Error(w, fmt.Sprintf("error: malformed task id: %v\n", taskID), 400)
		return
	}
	task, err := s.taskCache.Get(parsed)
	if err != nil {
		s.logger.Error("failed to get task from cache", "taskID", parsed, "error", err)
		http.Error(w, fmt.Sprintf("error: %v\n", err), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		s.logger.Error("failed to encode task response", "taskID", parsed, "error", err)
		http.Error(w, fmt.Sprintf("error: cannot encode response payload: %v", err), 500)
		return
	}

	s.logger.Info("successfully returned task", "taskID", parsed)
}

func (s *Server) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("handling post task request")

	task := NewTask()
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		s.logger.Error("failed to decode request body", "error", err)
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}

	s.logger.Info("received task", "taskID", task.ID)

	if err := task.ValidateFields(); err != nil {
		s.logger.Warn("task validation failed", "taskID", task.ID, "error", err)
		http.Error(w, fmt.Sprintf("%v\n", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := s.taskCache.Add(*task); err != nil {
		s.logger.Error("failed to add task to cache", "taskID", task.ID, "error", err)
		http.Error(w, fmt.Sprintf("%v\n", err), 500)
		return
	}

	s.logger.Info("task added to cache successfully", "taskID", task.ID)

	if err := json.NewEncoder(w).Encode(map[string]string{"task added with success": task.ID.String()}); err != nil {
		s.logger.Error("failed to encode response", "taskID", task.ID, "error", err)
		http.Error(w, fmt.Sprintf("error: cannot decode payload: %v", err), 500)
		return
	}
	taskBytes, err := task.Marshall()
	if err != nil {
		s.logger.Error("failed to marshal task", "taskID", task.ID, "error", err)
		http.Error(w, fmt.Sprintf("error: cannot marshal task: %v", err), 500)
		return
	}
	taskIDBytes, err := task.ID.MarshalBinary()
	if err != nil {
		s.logger.Error("failed to marshal task", "taskID", task.ID, "error", err)
		http.Error(w, fmt.Sprintf("error: cannot marshal task: %v", err), 500)
		return
	}
	if err := s.kafkaClient.Publish("tasks", taskIDBytes, taskBytes); err != nil {
		s.logger.Error("failed to send task to workers", "taskID", task.ID, "error", err)
		http.Error(w, fmt.Sprintf("error: failed to send task to workers: %v", err), 503)
		return
	}
	w.WriteHeader(202)
	s.logger.Info("task accepted", "taskID", task.ID)
}
