package services

import (
	"fmt"
	"strings"

	"service_provider_example/modules/user/repositories"
)

// UserServiceInterface defines the contract for user business logic
type UserServiceInterface interface {
	GetAllUsers() ([]repositories.User, error)
	GetUserByID(id int) (*repositories.User, error)
	CreateUser(name, email string) (*repositories.User, error)
	UpdateUser(id int, name, email string) (*repositories.User, error)
	DeleteUser(id int) error
	SearchUsers(query string) ([]repositories.User, error)
}

// UserService provides business logic for user operations
type UserService struct {
	repository repositories.UserRepositoryInterface
}

// NewUserService creates a new user service instance
func NewUserService(repository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		repository: repository,
	}
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() ([]repositories.User, error) {
	return s.repository.GetUsers()
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id int) (*repositories.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}
	return s.repository.GetUser(id)
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(name, email string) (*repositories.User, error) {
	// Validate input
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email format")
	}

	// Check for duplicate email
	users, err := s.repository.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to check for duplicate email: %w", err)
	}

	for _, user := range users {
		if strings.EqualFold(user.Email, email) {
			return nil, fmt.Errorf("user with email %s already exists", email)
		}
	}

	return s.repository.CreateUser(strings.TrimSpace(name), strings.TrimSpace(email))
}

// UpdateUser updates an existing user with validation
func (s *UserService) UpdateUser(id int, name, email string) (*repositories.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email format")
	}

	// Check if user exists
	_, err := s.repository.GetUser(id)
	if err != nil {
		return nil, err
	}

	return s.repository.UpdateUser(id, strings.TrimSpace(name), strings.TrimSpace(email))
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID: %d", id)
	}

	// Check if user exists
	_, err := s.repository.GetUser(id)
	if err != nil {
		return err
	}

	return s.repository.DeleteUser(id)
}

// SearchUsers searches for users by name or email
func (s *UserService) SearchUsers(query string) ([]repositories.User, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return s.repository.GetUsers()
	}

	users, err := s.repository.GetUsers()
	if err != nil {
		return nil, err
	}

	var results []repositories.User
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), query) ||
			strings.Contains(strings.ToLower(user.Email), query) {
			results = append(results, user)
		}
	}

	return results, nil
}
