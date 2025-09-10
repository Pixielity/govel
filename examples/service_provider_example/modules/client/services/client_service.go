package services

import (
	"fmt"
	"service_provider_example/modules/client/models"
	"strings"
	"sync"
	"time"
)

// ClientService provides the concrete implementation of ClientServiceInterface.
// This service manages client data using an in-memory storage system for demonstration purposes.
//
// In a real-world application, this would typically interact with a database
// or external API to persist and retrieve client information.
type ClientService struct {
	// clients holds the in-memory storage for client records
	clients map[int]*models.Client

	// nextID tracks the next available ID for new clients
	nextID int

	// mutex provides thread-safe access to the client data
	mutex sync.RWMutex

	// logger interface for logging operations
	logger interface{}
}

// NewClientService creates a new client service instance with initial test data.
//
// Parameters:
//   logger: Logger interface for logging service operations
//
// Returns:
//   *ClientService: A new client service instance with sample data
func NewClientService(logger interface{}) *ClientService {
	service := &ClientService{
		clients: make(map[int]*models.Client),
		nextID:  1,
		logger:  logger,
	}

	// Initialize with sample client data
	service.initializeSampleData()

	return service
}

// initializeSampleData populates the service with initial client records for testing.
func (cs *ClientService) initializeSampleData() {
	statuses := models.GetClientStatuses()
	now := time.Now()

	sampleClients := []*models.Client{
		{
			ID:             1,
			Name:           "Acme Corporation",
			Email:          "contact@acme.com",
			Company:        "Acme Corp",
			Phone:          "+1-555-0101",
			Status:         statuses.Active,
			RegisteredAt:   now.AddDate(0, -6, 0), // 6 months ago
			LastActivityAt: now.AddDate(0, 0, -2), // 2 days ago
		},
		{
			ID:             2,
			Name:           "TechStart Inc",
			Email:          "hello@techstart.io",
			Company:        "TechStart",
			Phone:          "+1-555-0202",
			Status:         statuses.Active,
			RegisteredAt:   now.AddDate(0, -3, 0), // 3 months ago
			LastActivityAt: now.AddDate(0, 0, -1), // 1 day ago
		},
		{
			ID:             3,
			Name:           "Global Solutions",
			Email:          "info@globalsolutions.com",
			Company:        "Global Solutions Ltd",
			Phone:          "+1-555-0303",
			Status:         statuses.Inactive,
			RegisteredAt:   now.AddDate(-1, 0, 0), // 1 year ago
			LastActivityAt: now.AddDate(0, -2, 0), // 2 months ago
		},
	}

	for _, client := range sampleClients {
		cs.clients[client.ID] = client
	}
	cs.nextID = 4 // Set next ID after sample data
}

// GetAllClients retrieves all clients from the system.
func (cs *ClientService) GetAllClients() ([]*models.Client, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	clients := make([]*models.Client, 0, len(cs.clients))
	for _, client := range cs.clients {
		clients = append(clients, client)
	}

	return clients, nil
}

// GetClientByID retrieves a specific client by their ID.
func (cs *ClientService) GetClientByID(id int) (*models.Client, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	client, exists := cs.clients[id]
	if !exists {
		return nil, fmt.Errorf("client with ID %d not found", id)
	}

	return client, nil
}

// CreateClient creates a new client with the specified information.
func (cs *ClientService) CreateClient(name, email, company, phone string) (*models.Client, error) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	statuses := models.GetClientStatuses()
	now := time.Now()

	client := &models.Client{
		ID:             cs.nextID,
		Name:           name,
		Email:          email,
		Company:        company,
		Phone:          phone,
		Status:         statuses.Active, // New clients start as active
		RegisteredAt:   now,
		LastActivityAt: now,
	}

	cs.clients[cs.nextID] = client
	cs.nextID++

	return client, nil
}

// UpdateClient updates an existing client's information.
func (cs *ClientService) UpdateClient(id int, name, email, company, phone string) (*models.Client, error) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	client, exists := cs.clients[id]
	if !exists {
		return nil, fmt.Errorf("client with ID %d not found", id)
	}

	// Update fields
	client.Name = name
	client.Email = email
	client.Company = company
	client.Phone = phone
	client.LastActivityAt = time.Now()

	return client, nil
}

// DeleteClient removes a client from the system.
func (cs *ClientService) DeleteClient(id int) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if _, exists := cs.clients[id]; !exists {
		return fmt.Errorf("client with ID %d not found", id)
	}

	delete(cs.clients, id)
	return nil
}

// GetClientsByStatus retrieves clients filtered by their status.
func (cs *ClientService) GetClientsByStatus(status string) ([]*models.Client, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	var filteredClients []*models.Client
	for _, client := range cs.clients {
		if client.Status == status {
			filteredClients = append(filteredClients, client)
		}
	}

	return filteredClients, nil
}

// GetClientsByCompany retrieves clients filtered by their company.
func (cs *ClientService) GetClientsByCompany(company string) ([]*models.Client, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	var filteredClients []*models.Client
	for _, client := range cs.clients {
		if strings.EqualFold(client.Company, company) {
			filteredClients = append(filteredClients, client)
		}
	}

	return filteredClients, nil
}

// SearchClients searches for clients by name, email, or company.
func (cs *ClientService) SearchClients(query string) ([]*models.Client, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	query = strings.ToLower(query)
	var matchingClients []*models.Client

	for _, client := range cs.clients {
		if strings.Contains(strings.ToLower(client.Name), query) ||
			strings.Contains(strings.ToLower(client.Email), query) ||
			strings.Contains(strings.ToLower(client.Company), query) {
			matchingClients = append(matchingClients, client)
		}
	}

	return matchingClients, nil
}

// UpdateClientStatus updates a client's status.
func (cs *ClientService) UpdateClientStatus(id int, status string) (*models.Client, error) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	client, exists := cs.clients[id]
	if !exists {
		return nil, fmt.Errorf("client with ID %d not found", id)
	}

	client.Status = status
	client.LastActivityAt = time.Now()

	return client, nil
}

// GetActiveClientsCount returns the number of active clients.
func (cs *ClientService) GetActiveClientsCount() (int, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	statuses := models.GetClientStatuses()
	count := 0

	for _, client := range cs.clients {
		if client.Status == statuses.Active {
			count++
		}
	}

	return count, nil
}

// GetClientStatistics returns basic statistics about clients.
func (cs *ClientService) GetClientStatistics() (map[string]int, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	statuses := models.GetClientStatuses()
	stats := map[string]int{
		"total":     len(cs.clients),
		"active":    0,
		"inactive":  0,
		"suspended": 0,
	}

	for _, client := range cs.clients {
		switch client.Status {
		case statuses.Active:
			stats["active"]++
		case statuses.Inactive:
			stats["inactive"]++
		case statuses.Suspended:
			stats["suspended"]++
		}
	}

	return stats, nil
}

// Compile-time interface compliance check
var _ ClientServiceInterface = (*ClientService)(nil)
