package workers

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ahmedennaifer/taskq/pkg"
	"github.com/google/uuid"
)

type Worker struct {
	ID           uuid.UUID
	State        string
	Addr         string
	CurrentTask  uuid.UUID
	Subscribtion any
}

func NewWorker(port string) (*Worker, error) {
	addr, err := getWorkerAddr()
	if err != nil {
		return &Worker{}, err
	}
	strAddr := addr.String()
	strAddr = strAddr[:len(strAddr)-3]
	return &Worker{
		ID:    uuid.New(),
		State: "Starting",
		Addr:  strAddr + port,
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
	fmt.Println("healthy")
	w.WriteHeader(http.StatusOK)
}

func (w *Worker) Register(managerUrl string) error {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return errors.New("secret key not found")
	}
	hash, err := pkg.Hash(w.ID.String(), secret)
	if err != nil {
		return errors.New("cannot hash worker id")
	}
	endpoint := managerUrl + "/worker"
	fmt.Println(endpoint)
	payload := pkg.RegisterPayload{
		Hash:     hash,
		WorkerID: w.ID.String(),
	}
	resp, err := pkg.PostRequest(endpoint, payload)
	if err != nil {
		return err
	}
	fmt.Println("response: ", string(resp))
	return nil
}
