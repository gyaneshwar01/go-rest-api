package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr  string
	store Store
}

func NewAPIServer(addr string, store Store) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// Registering our services
	userService := NewUserService(s.store)
	userService.RegisterRoutes(subRouter)

	tasksService := NewTasksService(s.store)
	tasksService.RegisterRoutes(subRouter)

	log.Println("starting the API server at", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subRouter))
}
