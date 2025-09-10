package repositories

import (
	"fmt"
)

// User represents a user entity
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRepositoryInterface defines the contract for user data access
type UserRepositoryInterface interface {
	GetUsers() ([]User, error)
	GetUser(id int) (*User, error)
	CreateUser(name, email string) (*User, error)
	UpdateUser(id int, name, email string) (*User, error)
	DeleteUser(id int) error
}

// UserRepository provides a simple in-memory implementation of user data access
type UserRepository struct {
	users  []User
	nextID int
}

// NewUserRepository creates a new user repository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: []User{
			{ID: 1, Name: "John Doe", Email: "john@example.com"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
		},
		nextID: 3,
	}
}

// GetUsers returns all users
func (r *UserRepository) GetUsers() ([]User, error) {
	return r.users, nil
}

// GetUser returns a user by ID
func (r *UserRepository) GetUser(id int) (*User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", id)
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(name, email string) (*User, error) {
	user := User{
		ID:    r.nextID,
		Name:  name,
		Email: email,
	}
	r.users = append(r.users, user)
	r.nextID++
	return &user, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(id int, name, email string) (*User, error) {
	for i, user := range r.users {
		if user.ID == id {
			r.users[i].Name = name
			r.users[i].Email = email
			return &r.users[i], nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", id)
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id int) error {
	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user with ID %d not found", id)
}
