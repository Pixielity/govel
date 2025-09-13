package main

import (
	"fmt"
	"log"

	"govel/new/pipeline/src"
	"govel/new/pipeline/src/interfaces"
)

// Order represents an order in our system
type Order struct {
	ID       string
	UserID   int
	Items    []OrderItem
	Total    float64
	Status   string
	Metadata map[string]interface{}
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

// InventoryMiddleware checks inventory
func InventoryMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	order := passable.(*Order)
	fmt.Printf("[Inventory] Checking inventory for order %s\n", order.ID)
	
	// Simulate inventory check
	for _, item := range order.Items {
		fmt.Printf("[Inventory] Checking %d units of product %s\n", item.Quantity, item.ProductID)
		
		// Simulate insufficient inventory for product "out-of-stock"
		if item.ProductID == "out-of-stock" {
			return nil, fmt.Errorf("insufficient inventory for product %s", item.ProductID)
		}
	}
	
	fmt.Println("[Inventory] All items available")
	return next(passable)
}

// PaymentMiddleware processes payment
func PaymentMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	order := passable.(*Order)
	fmt.Printf("[Payment] Processing payment of $%.2f for order %s\n", order.Total, order.ID)
	
	// Simulate payment processing
	if order.Total > 10000 {
		return nil, fmt.Errorf("payment amount too high: $%.2f", order.Total)
	}
	
	order.Status = "paid"
	fmt.Println("[Payment] Payment processed successfully")
	return next(passable)
}

// ShippingMiddleware handles shipping
func ShippingMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	order := passable.(*Order)
	fmt.Printf("[Shipping] Arranging shipping for order %s\n", order.ID)
	
	// Add shipping information
	if order.Metadata == nil {
		order.Metadata = make(map[string]interface{})
	}
	order.Metadata["shipping"] = map[string]interface{}{
		"carrier":     "FastShip",
		"tracking_id": "FS" + order.ID,
		"estimated":   "2-3 business days",
	}
	
	fmt.Println("[Shipping] Shipping arranged")
	return next(passable)
}

// NotificationMiddleware sends notifications
func NotificationMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	order := passable.(*Order)
	fmt.Printf("[Notification] Sending notifications for order %s\n", order.ID)
	
	// Simulate sending email/SMS
	fmt.Printf("[Notification] Email sent to user %d\n", order.UserID)
	fmt.Printf("[Notification] SMS sent to user %d\n", order.UserID)
	
	return next(passable)
}

// ExampleHubNamedPipelines demonstrates using Hub with different named pipelines for different order types
func ExampleHubNamedPipelines() {
	fmt.Println("=== Hub Named Pipelines Example ===")
	
	// Create hub
	hub := pipeline.NewHub(nil)
	
	// Define "standard" order pipeline
	hub.Pipeline("standard", func(p interfaces.PipelineInterface, passable interface{}) interface{} {
		fmt.Println("\n--- Processing Standard Order ---")
		
		result, err := p.
			Send(passable).
			Through([]interface{}{InventoryMiddleware, PaymentMiddleware, ShippingMiddleware}).
			Then(func(passable interface{}) interface{} {
				order := passable.(*Order)
				order.Status = "completed"
				fmt.Printf("[Standard] Order %s completed successfully\n", order.ID)
				return order
			})
		
		if err != nil {
			fmt.Printf("[Standard] Order failed: %v\n", err)
			return nil
		}
		
		return result
	})
	
	// Define "express" order pipeline (includes priority processing and notifications)
	hub.Pipeline("express", func(p interfaces.PipelineInterface, passable interface{}) interface{} {
		fmt.Println("\n--- Processing Express Order ---")
		
		result, err := p.
			Send(passable).
			Through([]interface{}{InventoryMiddleware, PaymentMiddleware, ShippingMiddleware, NotificationMiddleware}).
			Then(func(passable interface{}) interface{} {
				order := passable.(*Order)
				order.Status = "express-completed"
				// Upgrade shipping for express orders
				if shipping, ok := order.Metadata["shipping"].(map[string]interface{}); ok {
					shipping["carrier"] = "ExpressShip"
					shipping["estimated"] = "Next business day"
				}
				fmt.Printf("[Express] Express order %s completed with priority shipping\n", order.ID)
				return order
			})
		
		if err != nil {
			fmt.Printf("[Express] Express order failed: %v\n", err)
			return nil
		}
		
		return result
	})
	
	// Define "digital" order pipeline (no shipping, immediate delivery)
	hub.Pipeline("digital", func(p interfaces.PipelineInterface, passable interface{}) interface{} {
		fmt.Println("\n--- Processing Digital Order ---")
		
		result, err := p.
			Send(passable).
			Through([]interface{}{PaymentMiddleware, NotificationMiddleware}). // No inventory or shipping for digital
			Then(func(passable interface{}) interface{} {
				order := passable.(*Order)
				order.Status = "digital-delivered"
				if order.Metadata == nil {
					order.Metadata = make(map[string]interface{})
				}
				order.Metadata["delivery"] = map[string]interface{}{
					"type":         "instant",
					"download_url": "https://downloads.example.com/" + order.ID,
					"expires":      "30 days",
				}
				fmt.Printf("[Digital] Digital order %s delivered instantly\n", order.ID)
				return order
			})
		
		if err != nil {
			fmt.Printf("[Digital] Digital order failed: %v\n", err)
			return nil
		}
		
		return result
	})
	
	// Test standard order
	standardOrder := &Order{
		ID:     "ORD001",
		UserID: 123,
		Items: []OrderItem{
			{ProductID: "BOOK001", Quantity: 2, Price: 15.99},
			{ProductID: "PEN001", Quantity: 5, Price: 2.50},
		},
		Total: 44.48,
	}
	
	result1, err1 := hub.Pipe(standardOrder, "standard")
	if err1 != nil {
		log.Printf("Standard order failed: %v", err1)
	} else {
		if order, ok := result1.(*Order); ok {
			fmt.Printf("Standard order result: %s (Status: %s)\n", order.ID, order.Status)
		}
	}
	
	// Test express order
	expressOrder := &Order{
		ID:     "ORD002",
		UserID: 456,
		Items: []OrderItem{
			{ProductID: "LAPTOP001", Quantity: 1, Price: 999.99},
		},
		Total: 999.99,
	}
	
	result2, err2 := hub.Pipe(expressOrder, "express")
	if err2 != nil {
		log.Printf("Express order failed: %v", err2)
	} else {
		if order, ok := result2.(*Order); ok {
			fmt.Printf("Express order result: %s (Status: %s)\n", order.ID, order.Status)
			if shipping := order.Metadata["shipping"]; shipping != nil {
				fmt.Printf("Shipping: %+v\n", shipping)
			}
		}
	}
	
	// Test digital order
	digitalOrder := &Order{
		ID:     "ORD003",
		UserID: 789,
		Items: []OrderItem{
			{ProductID: "EBOOK001", Quantity: 1, Price: 9.99},
			{ProductID: "SOFTWARE001", Quantity: 1, Price: 49.99},
		},
		Total: 59.98,
	}
	
	result3, err3 := hub.Pipe(digitalOrder, "digital")
	if err3 != nil {
		log.Printf("Digital order failed: %v", err3)
	} else {
		if order, ok := result3.(*Order); ok {
			fmt.Printf("Digital order result: %s (Status: %s)\n", order.ID, order.Status)
			if delivery := order.Metadata["delivery"]; delivery != nil {
				fmt.Printf("Delivery: %+v\n", delivery)
			}
		}
	}
	
	// Test error case - insufficient inventory
	fmt.Println("\n--- Testing Error Case ---")
	errorOrder := &Order{
		ID:     "ORD004",
		UserID: 999,
		Items: []OrderItem{
			{ProductID: "out-of-stock", Quantity: 1, Price: 19.99},
		},
		Total: 19.99,
	}
	
	_, err4 := hub.Pipe(errorOrder, "standard")
	if err4 != nil {
		fmt.Printf("Expected error occurred: %v\n", err4)
	}
}

func main() {
	ExampleHubNamedPipelines()
}
