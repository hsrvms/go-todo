package server

import (
	"log"
	"net/http"

	"github.com/hsrvms/todoapp/services"
	"github.com/hsrvms/todoapp/store"
)

type APIServer struct {
	addr       string
	repository store.Store
}

func NewAPIServer(addr string, repository store.Store) *APIServer {
	return &APIServer{
		addr:       addr,
		repository: repository,
	}
}

func (s *APIServer) Start() {
	const v1Prefix = "/api/v1"
	userService := services.NewUserService(s.repository)
	taskService := services.NewTaskService(s.repository)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("health"))
	})

	userService.RegisterRoutes(mux, v1Prefix)
	taskService.RegisterRoutes(mux, v1Prefix)

	log.Println("Starting API server on", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, mux))
}
