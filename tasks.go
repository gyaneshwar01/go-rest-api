package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TasksService struct {
	store Store
}

func NewTasksService(s Store) *TasksService {
	return &TasksService{store: s}
}

func (s *TasksService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", WithJWTAuth(s.handleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/{id}", WithJWTAuth(s.handleGetTask, s.store)).Methods("GET")
}

func (s *TasksService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	defer r.Body.Close()

	var task *Task

	err = json.Unmarshal(body, &task)

	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	t, err := s.store.CreateTask(task)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJson(w, http.StatusCreated, t)
}

func (s *TasksService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	t, err := s.store.GetTask(id)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJson(w, http.StatusOK, t)
}

func validateTaskPayload(task *Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}
