package repositories

import (
	"fmt"
)

// Product represents a product entity
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

// ProductRepositoryInterface defines the contract for product data access
type ProductRepositoryInterface interface {
	GetProducts() ([]Product, error)
	GetProduct(id int) (*Product, error)
	GetProductsByCategory(category string) ([]Product, error)
	CreateProduct(name, description, category string, price float64) (*Product, error)
	UpdateProduct(id int, name, description, category string, price float64) (*Product, error)
	DeleteProduct(id int) error
}

// ProductRepository provides a simple in-memory implementation of product data access
type ProductRepository struct {
	products []Product
	nextID   int
}

// NewProductRepository creates a new product repository instance
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: []Product{
			{ID: 1, Name: "Laptop", Description: "High-performance laptop", Price: 1299.99, Category: "Electronics"},
			{ID: 2, Name: "Coffee Mug", Description: "Ceramic coffee mug", Price: 12.99, Category: "Kitchen"},
			{ID: 3, Name: "Book", Description: "Programming book", Price: 39.99, Category: "Books"},
		},
		nextID: 4,
	}
}

// GetProducts returns all products
func (r *ProductRepository) GetProducts() ([]Product, error) {
	return r.products, nil
}

// GetProduct returns a product by ID
func (r *ProductRepository) GetProduct(id int) (*Product, error) {
	for _, product := range r.products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, fmt.Errorf("product with ID %d not found", id)
}

// GetProductsByCategory returns products by category
func (r *ProductRepository) GetProductsByCategory(category string) ([]Product, error) {
	var results []Product
	for _, product := range r.products {
		if product.Category == category {
			results = append(results, product)
		}
	}
	return results, nil
}

// CreateProduct creates a new product
func (r *ProductRepository) CreateProduct(name, description, category string, price float64) (*Product, error) {
	product := Product{
		ID:          r.nextID,
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
	}
	r.products = append(r.products, product)
	r.nextID++
	return &product, nil
}

// UpdateProduct updates an existing product
func (r *ProductRepository) UpdateProduct(id int, name, description, category string, price float64) (*Product, error) {
	for i, product := range r.products {
		if product.ID == id {
			r.products[i].Name = name
			r.products[i].Description = description
			r.products[i].Category = category
			r.products[i].Price = price
			return &r.products[i], nil
		}
	}
	return nil, fmt.Errorf("product with ID %d not found", id)
}

// DeleteProduct deletes a product by ID
func (r *ProductRepository) DeleteProduct(id int) error {
	for i, product := range r.products {
		if product.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product with ID %d not found", id)
}
