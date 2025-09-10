package services

import (
	"fmt"
	"strings"

	"service_provider_example/modules/product/repositories"
)

// ProductServiceInterface defines the contract for product business logic
type ProductServiceInterface interface {
	GetAllProducts() ([]repositories.Product, error)
	GetProductByID(id int) (*repositories.Product, error)
	GetProductsByCategory(category string) ([]repositories.Product, error)
	CreateProduct(name, description, category string, price float64) (*repositories.Product, error)
	UpdateProduct(id int, name, description, category string, price float64) (*repositories.Product, error)
	DeleteProduct(id int) error
	SearchProducts(query string) ([]repositories.Product, error)
	GetCategories() ([]string, error)
}

// ProductService provides business logic for product operations
type ProductService struct {
	repository repositories.ProductRepositoryInterface
}

// NewProductService creates a new product service instance
func NewProductService(repository repositories.ProductRepositoryInterface) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

// GetAllProducts returns all products
func (s *ProductService) GetAllProducts() ([]repositories.Product, error) {
	return s.repository.GetProducts()
}

// GetProductByID returns a product by ID
func (s *ProductService) GetProductByID(id int) (*repositories.Product, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid product ID: %d", id)
	}
	return s.repository.GetProduct(id)
}

// GetProductsByCategory returns products by category
func (s *ProductService) GetProductsByCategory(category string) ([]repositories.Product, error) {
	if strings.TrimSpace(category) == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}
	return s.repository.GetProductsByCategory(strings.TrimSpace(category))
}

// CreateProduct creates a new product with validation
func (s *ProductService) CreateProduct(name, description, category string, price float64) (*repositories.Product, error) {
	// Validate input
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if strings.TrimSpace(description) == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}
	if strings.TrimSpace(category) == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}

	return s.repository.CreateProduct(
		strings.TrimSpace(name),
		strings.TrimSpace(description),
		strings.TrimSpace(category),
		price,
	)
}

// UpdateProduct updates an existing product with validation
func (s *ProductService) UpdateProduct(id int, name, description, category string, price float64) (*repositories.Product, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid product ID: %d", id)
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if strings.TrimSpace(description) == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}
	if strings.TrimSpace(category) == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}

	// Check if product exists
	_, err := s.repository.GetProduct(id)
	if err != nil {
		return nil, err
	}

	return s.repository.UpdateProduct(
		id,
		strings.TrimSpace(name),
		strings.TrimSpace(description),
		strings.TrimSpace(category),
		price,
	)
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid product ID: %d", id)
	}

	// Check if product exists
	_, err := s.repository.GetProduct(id)
	if err != nil {
		return err
	}

	return s.repository.DeleteProduct(id)
}

// SearchProducts searches for products by name, description, or category
func (s *ProductService) SearchProducts(query string) ([]repositories.Product, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return s.repository.GetProducts()
	}

	products, err := s.repository.GetProducts()
	if err != nil {
		return nil, err
	}

	var results []repositories.Product
	for _, product := range products {
		if strings.Contains(strings.ToLower(product.Name), query) ||
			strings.Contains(strings.ToLower(product.Description), query) ||
			strings.Contains(strings.ToLower(product.Category), query) {
			results = append(results, product)
		}
	}

	return results, nil
}

// GetCategories returns all unique product categories
func (s *ProductService) GetCategories() ([]string, error) {
	products, err := s.repository.GetProducts()
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]bool)
	for _, product := range products {
		categoryMap[product.Category] = true
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}

	return categories, nil
}
