package workers

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ahmedennaifer/taskq/pkg"
	"github.com/google/uuid"
)

type WorkerStatus string

const (
	StatusHealthy   WorkerStatus = "healthy"
	StatusUnhealthy WorkerStatus = "unhealthy"
	StatusStarting  WorkerStatus = "starting"
)

type Worker struct {
	ID           uuid.UUID
	Addr         string
	CurrentTask  uuid.UUID
	Subscription any
	logger       *slog.Logger
	Status       WorkerStatus
	LastSeen     time.Time
	LastChecked  time.Time
	Hash         string
}

func NewWorker(port string, logger *slog.Logger) (*Worker, error) {
	logger.Info("creating new worker", "port", port)

	addr, err := getWorkerAddr()
	if err != nil {
		logger.Error("failed to get worker address", "error", err)
		return &Worker{}, err
	}

	strAddr := addr.String()
	strAddr = strAddr[:len(strAddr)-3]
	workerID := uuid.New()

	// hash will be used to register itself

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		logger.Error("secret key not found in environment")
		return &Worker{}, errors.New("secret key not found")
	}

	hash, err := pkg.Hash(workerID.String(), secret)
	if err != nil {
		logger.Error("failed to hash worker ID", "workerID", workerID, "error", err)
		return &Worker{}, errors.New("cannot hash worker id")
	}

	logger.Info("worker created successfully", "workerID", workerID, "addr", strAddr+port)

	return &Worker{
		ID:       workerID,
		Status:   StatusStarting,
		Addr:     strAddr + port,
		logger:   logger,
		Hash:     hash,
		LastSeen: time.Now(),
	}, nil
}

func getWorkerAddr() (*net.IPNet, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		err := fmt.Sprintf("error retrieving worker ip: %v\n", err)
		return nil, errors.New(err)
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet, nil
		}
	}
	return nil, errors.New("cannot resolve worker address")
}

func (w *Worker) Init() error {
	return nil
}

func (wrk *Worker) HandleHealth(w http.ResponseWriter, r *http.Request) {
	wrk.logger.Debug("health check request received")
	w.WriteHeader(http.StatusOK)
	wrk.logger.Debug("health check response sent")
}

func (wrk *Worker) Register(managerUrl string) error {
	wrk.logger.Info("registering worker with manager", "managerUrl", managerUrl, "workerID", wrk.ID)
	endpoint := managerUrl + "/worker"
	wrk.logger.Debug("registration endpoint", "endpoint", endpoint)

	resp, err := pkg.PostRequest(endpoint, wrk)
	if err != nil {
		wrk.logger.Error("registration request failed", "endpoint", endpoint, "error", err)
		return err
	}

	wrk.logger.Info("worker registered successfully", "workerID", wrk.ID, "response", string(resp))
	return nil
}
