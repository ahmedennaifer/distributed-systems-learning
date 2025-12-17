package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ahmedennaifer/taskq/internal/workers"
)

func main() {
	w, _ := workers.NewWorker(":8001")
	err := w.Register("http://localhost:8081/api/v1")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully registered worker ", w.ID.String())
	fmt.Println("addr: ", w.Addr)
	http.HandleFunc("GET /health", w.HandleHealth)
	log.Fatal(http.ListenAndServe(w.Addr, nil))
}
