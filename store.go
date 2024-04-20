package main

import "database/sql"

type Store interface {
	// Users
	CreateUser(*User) (*User, error)
	GetUserById(id string) (*User, error)
	// Tasks
	CreateTask(*Task) (*Task, error)
	GetTask(string) (*Task, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(user *User) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO users (email, firstName, lastName, password) VALUES (?, ?, ?, ?, ?)", user.Email, user.FirstName, user.LastName, user.Password)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *Storage) GetUserById(id string) (*User, error) {
	var u User

	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt FROM users WHERE id=?", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)

	return &u, err
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name, status, project_id, assigned_to) VALUES (?, ?, ?, ?)", t.Name, t.Status, t.ProjectID, t.AssignedToID)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id
	return t, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task

	err := s.db.QueryRow("SELECT id, name, status, projectID, assignedToID, createdAt FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)

	return &t, err
}
