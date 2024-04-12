package store

import "github.com/hsrvms/todoapp/models"

type MockStore struct {
	store []*models.User
}

func NewMockStore() *MockStore {
	mockStore := &MockStore{}
	mockStore.store = append(
		mockStore.store,
		&models.User{
			ID:       1,
			Username: "testUserLogin",
			Password: "testPassword",
		})

	return mockStore
}

// User
func (ms *MockStore) CreateUser(u *models.User) (*models.User, error) {
	return &models.User{}, nil
}
func (ms *MockStore) GetUserByID(id string) (*models.User, error) {
	return &models.User{}, nil
}
func (ms *MockStore) GetUserByUsername(username string) (*models.User, error) {
	return &models.User{}, nil
}

// Task
func (ms *MockStore) CreateTask(t *models.Task) (*models.Task, error) {
	return &models.Task{}, nil
}

func (ms *MockStore) GetAllTasks() ([]*models.Task, error) {
	return []*models.Task{}, nil
}
func (ms *MockStore) GetTaskByID(id string) (*models.Task, error) {
	return &models.Task{}, nil
}
func (ms *MockStore) UpdateTask(id string, t *models.Task) (*models.Task, error) {
	return &models.Task{}, nil
}
func (ms *MockStore) DeleteTask(id string) (*models.Task, error) {
	return &models.Task{}, nil
}
