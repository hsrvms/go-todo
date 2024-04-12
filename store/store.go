package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/hsrvms/todoapp/models"
)

type Store interface {
	// Users
	CreateUser(u *models.User) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)

	// Task
	CreateTask(t *models.Task) (*models.Task, error)
	GetAllTasks() ([]*models.Task, error)
	GetTaskByID(id string) (*models.Task, error)
	UpdateTask(id string, t *models.Task) (*models.Task, error)
	DeleteTask(id string) (*models.Task, error)
}
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser creates a new user in the repository.
func (r *Repository) CreateUser(u *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(u.Username, u.Password).Scan(&u.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetUserByID retrieves a user by their ID from the repository.
func (r *Repository) GetUserByID(id string) (*models.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}

	user := &models.User{}
	query := `
		SELECT id, username, created_at 
		FROM users 
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username is empty")
	}

	user := &models.User{}
	query := `
		SELECT id, username, password
		FROM users
		WHERE username = $1
	`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) CreateTask(t *models.Task) (*models.Task, error) {
	if t == nil {
		return nil, fmt.Errorf("task is nil")
	}

	query := `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(t.Title, t.Description, t.Status).Scan(&t.ID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *Repository) GetAllTasks() ([]*models.Task, error) {
	query := `
		SELECT * FROM tasks
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetTaskByID(id string) (*models.Task, error) {
	if id == "" {
		return nil, errors.New("task ID cannot be empty")
	}

	query := "SELECT * FROM tasks WHERE id = $1"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	task := &models.Task{}
	err = stmt.QueryRow(id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return task, nil
}
func (r *Repository) UpdateTask(id string, t *models.Task) (*models.Task, error) {
	query := `
		UPDATE tasks SET
		title = $1,
		description = $2,
		status = $3
		WHERE id = $4
		RETURNING id, title, description, status, created_at
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	task := &models.Task{}
	err = stmt.QueryRow(t.Title, t.Description, t.Status, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return task, nil
}
func (r *Repository) DeleteTask(taskID string) (*models.Task, error) {
	query := `
		DELETE FROM tasks
		WHERE id = $1
		RETURNING id, title
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	task := &models.Task{}
	err = stmt.QueryRow(taskID).Scan(
		&task.ID,
		&task.Title,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return task, err
}
