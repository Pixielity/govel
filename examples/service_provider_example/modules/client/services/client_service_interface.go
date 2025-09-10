package services

import "service_provider_example/modules/client/models"

// ClientServiceInterface defines the contract for client service operations.
// This interface provides methods for managing clients, including CRUD operations,
// status management, and business logic for client interactions.
//
// All client-related business logic should be implemented through this interface
// to ensure consistency and testability across the application.
type ClientServiceInterface interface {
	// GetAllClients retrieves all clients from the system.
	//
	// Returns:
	//   []*models.Client: A slice of all client records
	//   error: Any error that occurred during retrieval
	GetAllClients() ([]*models.Client, error)

	// GetClientByID retrieves a specific client by their ID.
	//
	// Parameters:
	//   id: The unique identifier of the client to retrieve
	//
	// Returns:
	//   *models.Client: The client record if found
	//   error: Any error that occurred during retrieval
	GetClientByID(id int) (*models.Client, error)

	// CreateClient creates a new client with the specified information.
	//
	// Parameters:
	//   name: The client's full name
	//   email: The client's email address
	//   company: The client's company name
	//   phone: The client's phone number
	//
	// Returns:
	//   *models.Client: The newly created client record
	//   error: Any error that occurred during creation
	CreateClient(name, email, company, phone string) (*models.Client, error)

	// UpdateClient updates an existing client's information.
	//
	// Parameters:
	//   id: The unique identifier of the client to update
	//   name: The updated client name
	//   email: The updated client email
	//   company: The updated company name
	//   phone: The updated phone number
	//
	// Returns:
	//   *models.Client: The updated client record
	//   error: Any error that occurred during update
	UpdateClient(id int, name, email, company, phone string) (*models.Client, error)

	// DeleteClient removes a client from the system.
	//
	// Parameters:
	//   id: The unique identifier of the client to delete
	//
	// Returns:
	//   error: Any error that occurred during deletion
	DeleteClient(id int) error

	// GetClientsByStatus retrieves clients filtered by their status.
	//
	// Parameters:
	//   status: The status to filter by (active, inactive, suspended)
	//
	// Returns:
	//   []*models.Client: A slice of clients matching the specified status
	//   error: Any error that occurred during retrieval
	GetClientsByStatus(status string) ([]*models.Client, error)

	// GetClientsByCompany retrieves clients filtered by their company.
	//
	// Parameters:
	//   company: The company name to filter by
	//
	// Returns:
	//   []*models.Client: A slice of clients from the specified company
	//   error: Any error that occurred during retrieval
	GetClientsByCompany(company string) ([]*models.Client, error)

	// SearchClients searches for clients by name, email, or company.
	//
	// Parameters:
	//   query: The search query string
	//
	// Returns:
	//   []*models.Client: A slice of clients matching the search criteria
	//   error: Any error that occurred during search
	SearchClients(query string) ([]*models.Client, error)

	// UpdateClientStatus updates a client's status.
	//
	// Parameters:
	//   id: The unique identifier of the client
	//   status: The new status for the client
	//
	// Returns:
	//   *models.Client: The updated client record
	//   error: Any error that occurred during status update
	UpdateClientStatus(id int, status string) (*models.Client, error)

	// GetActiveClientsCount returns the number of active clients.
	//
	// Returns:
	//   int: The number of active clients
	//   error: Any error that occurred during counting
	GetActiveClientsCount() (int, error)

	// GetClientStatistics returns basic statistics about clients.
	//
	// Returns:
	//   map[string]int: A map containing various client statistics
	//   error: Any error that occurred during statistics gathering
	GetClientStatistics() (map[string]int, error)
}
