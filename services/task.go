package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/hsrvms/todoapp/auth"
	"github.com/hsrvms/todoapp/models"
	"github.com/hsrvms/todoapp/store"
	"github.com/hsrvms/todoapp/utils"
)

var ErrTitleRequired = errors.New("title is required")

type TaskService struct {
	store store.Store
}

func NewTaskService(store store.Store) *TaskService {
	return &TaskService{store: store}
}

// # POST /tasks:
//
// Payload:
//
//	{
//	 "title": "Learn Golang",
//	 "description": "Learning process of Golang",
//	}
//
// Response:
//
//	{
//	 "id": 1,
//	 "title": "Learn Golang",
//	 "description": "Learning process of Golang",
//	 "status": false,
//	 "created_at": "2024-04-12 18:02:27.924693",
//	}
//
// # GET /tasks:
//
// Response:
//
//	[
//	 {
//		"id": 1,
//		"title": "Learn Golang",
//		"description": "Learning process of Golang",
//		"status": false,
//		"created_at": "2024-04-12 18:02:27.924693",
//	 },
//	]
//
// # GET /tasks/{id}:
//
// Response:
//
//	{
//	 "id": 1,
//	 "title": "Learn Golang",
//	 "description": "Learning process of Golang",
//	 "status": false,
//	 "created_at": "2024-04-12 18:02:27.924693",
//	}
//
// # PUT /tasks/{id}:
//
// Payload:
//
//	{
//	 "title": "Learn Golang +",
//	 "description": "Learning process of Golang",
//	 "status": false,
//	}
//
// Response:
//
//	{
//	 "id": 1,
//	 "title": "Learn Golang +",
//	 "description": "Learning process of Golang",
//	 "status": false,
//	 "created_at": "2024-04-12 18:02:27.924693",
//	}
//
// # DELETE /tasks/{id}:
//
// Response:
//
//	{
//	 "id": 1,
//	 "title": "Learn Golang",
//	 "description": "Learning process of Golang",
//	 "status": false,
//	 "created_at": "2024-04-12 18:02:27.924693",
//	}
func (s *TaskService) RegisterRoutes(mux *http.ServeMux, prefix string) {
	endpointCreate := generateEndpoint("POST", prefix, "/tasks")
	endpointGetAll := generateEndpoint("GET", prefix, "/tasks")
	endpointGetByID := generateEndpoint("GET", prefix, "/tasks/{id}")
	endpointUpdate := generateEndpoint("PUT", prefix, "/tasks/{id}")
	endpointDelete := generateEndpoint("DELETE", prefix, "/tasks/{id}")

	mux.HandleFunc(endpointCreate, auth.WithJWTAuth(s.handleTaskCreate, s.store))
	mux.HandleFunc(endpointGetAll, auth.WithJWTAuth(s.handleTaskGetAll, s.store))
	mux.HandleFunc(endpointGetByID, auth.WithJWTAuth(s.handleTaskGetByID, s.store))
	mux.HandleFunc(endpointUpdate, auth.WithJWTAuth(s.handleTaskUpdate, s.store))
	mux.HandleFunc(endpointDelete, auth.WithJWTAuth(s.handleTaskDelete, s.store))
}

func (s *TaskService) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := decodeJSON(r, &task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validateTaskPayload(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := s.store.CreateTask(&task)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	if createdTask == nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdTask)
}

func (s *TaskService) handleTaskGetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.GetAllTasks()
	if err != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks)
}

func (s *TaskService) handleTaskGetByID(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	if taskID == "" {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := s.store.GetTaskByID(taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving task", http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}

func (s *TaskService) handleTaskUpdate(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	if taskID == "" {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := decodeJSON(r, &task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validateTaskPayload(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := s.store.UpdateTask(taskID, &task)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	if updatedTask == nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedTask)
}

func (s *TaskService) handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	if taskID == "" {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	deletedTask, err := s.store.DeleteTask(taskID)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	if deletedTask == nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, deletedTask)
}

func validateTaskPayload(task *models.Task) error {
	if task.Title == "" {
		return ErrTitleRequired
	}

	return nil
}
