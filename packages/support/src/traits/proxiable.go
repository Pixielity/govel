// Package traits provides a generic proxy pattern for Go structs that need self-reference capabilities.
// This solves the common issue in Go where embedded structs need to call methods on their concrete
// implementations rather than the base embedded type.
package traits

import (
	reflector "govel/support/src/reflector"
	"reflect"
	"sync"
)

// Proxiable provides self-reference capabilities for embedded structs.
//
// This solves the common Go pattern where you have a base struct (like Manager)
// embedded in concrete implementations (like HashManager), and you need the base
// struct's methods to call the concrete implementation's methods.
//
// Usage:
//
//	type BaseManager struct {
//	    Proxiable
//	    // ... other fields
//	}
//
//	type ConcreteManager struct {
//	    *BaseManager
//	    // ... concrete fields
//	}
//
//	func NewConcreteManager() *ConcreteManager {
//	    base := &BaseManager{}
//	    concrete := &ConcreteManager{BaseManager: base}
//	    base.SetProxySelf(concrete)  // Enable self-reference
//	    return concrete
//	}
type Proxiable struct {
	// self holds the reference to the concrete implementation
	self interface{}
	// mutex protects concurrent access to self
	mutex sync.RWMutex
}

// SetProxySelf sets the self-reference to the concrete implementation.
// This should be called after creating the concrete struct that embeds the base struct.
//
// Parameters:
//   - self: The concrete implementation that embeds this Proxiable
//
// Example:
//
//	manager := &ConcreteManager{BaseManager: &BaseManager{}}
//	manager.BaseManager.SetProxySelf(manager)
func (p *Proxiable) SetProxySelf(self interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.self = self
}

// GetProxySelf returns the current self-reference.
// This can be used by embedded struct methods to access the concrete implementation.
//
// Returns:
//   - interface{}: The concrete implementation, or nil if not set
//
// Example:
//
//	func (b *BaseManager) SomeMethod() {
//	    if concrete := b.GetProxySelf(); concrete != nil {
//	        // Call methods on the concrete type
//	        if cm, ok := concrete.(*ConcreteManager); ok {
//	            return cm.ConcreteMethod()
//	        }
//	    }
//	    // Fallback behavior
//	}
func (p *Proxiable) GetProxySelf() interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.self
}

// HasProxySelf checks if a self-reference has been set.
//
// Returns:
//   - bool: true if self-reference is set, false otherwise
func (p *Proxiable) HasProxySelf() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.self != nil
}

// CallOnSelf calls a method on the self-reference using reflection.
// This is a convenience method for calling methods on the concrete implementation
// when you know the method name but not the exact type.
//
// Parameters:
//   - methodName: Name of the method to call
//   - args: Arguments to pass to the method
//
// Returns:
//   - []reflect.Value: Results from the method call
//   - error: Any error that occurred during the call
//
// Example:
//
//	results, err := p.CallOnSelf("GetDefaultDriver")
//	if err == nil && len(results) > 0 {
//	    defaultDriver := results[0].String()
//	}
func (p *Proxiable) CallOnSelf(methodName string, args ...interface{}) ([]reflect.Value, error) {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return nil, &ProxyError{
			Op:      "CallOnSelf",
			Method:  methodName,
			Message: "no self-reference set",
		}
	}

	// Use our custom reflector instead of standard reflection
	if !reflector.HasMethod(self, methodName) {
		return nil, &ProxyError{
			Op:      "CallOnSelf",
			Method:  methodName,
			Message: "method not found on self-reference",
			Type:    reflector.GetTypeName(self),
		}
	}

	// Check if method is public
	if !reflector.IsMethodPublic(self, methodName) {
		return nil, &ProxyError{
			Op:      "CallOnSelf",
			Method:  methodName,
			Message: "method is not public (not exported)",
			Type:    reflector.GetTypeName(self),
		}
	}

	// Call the method using our custom reflector
	return reflector.CallMethod(self, methodName, args...)
}

// Call is a shorter alias for CallOnSelf.
// This provides a more concise interface while remaining Go-idiomatic.
// Functionally identical to CallOnSelf but with shorter naming.
//
// Parameters:
//   - methodName: Name of the method to call
//   - args: Arguments to pass to the method
//
// Returns:
//   - []reflect.Value: Results from the method call
//   - error: Any error that occurred during the call
//
// Example:
//
//	// Shorter usage
//	results, err := p.Call("SomeMethod", arg1, arg2)
//	if err == nil {
//	    // Process results
//	}
func (p *Proxiable) Call(methodName string, args ...interface{}) ([]reflect.Value, error) {
	return p.CallOnSelf(methodName, args...)
}

// GetSelfType returns the reflect.Type of the self-reference.
// This is useful for type checking and debugging.
//
// Returns:
//   - reflect.Type: The type of the self-reference, or nil if not set
func (p *Proxiable) GetSelfType() reflect.Type {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.self == nil {
		return nil
	}
	return reflect.TypeOf(p.self)
}

// GetSelfTypeName returns the fully qualified type name of the self-reference.
// This uses our custom reflector for consistent type naming.
//
// Returns:
//   - string: The type name, or "nil" if not set
func (p *Proxiable) GetSelfTypeName() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.self == nil {
		return "nil"
	}
	return reflector.GetTypeName(p.self)
}

// FindMethodOnSelf checks if a method exists on the self-reference.
// This is useful for conditional method calling.
//
// Parameters:
//   - methodName: Name of the method to find
//
// Returns:
//   - reflect.Method: The method if found
//   - bool: true if method exists, false otherwise
func (p *Proxiable) FindMethodOnSelf(methodName string) (reflect.Method, bool) {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return reflect.Method{}, false
	}

	// Use our custom reflector for method lookup
	return reflector.GetMethodByName(self, methodName)
}

// HasMethodOnSelf checks if a method exists on the self-reference using our custom reflector.
// This provides better integration with our reflection utilities.
//
// Parameters:
//   - methodName: Name of the method to check
//
// Returns:
//   - bool: true if method exists, false otherwise
func (p *Proxiable) HasMethodOnSelf(methodName string) bool {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return false
	}

	return reflector.HasMethod(self, methodName)
}

// IsMethodPublicOnSelf checks if a method on the self-reference is public (exported).
// This uses our custom reflector for consistent checking.
//
// Parameters:
//   - methodName: Name of the method to check
//
// Returns:
//   - bool: true if method exists and is public, false otherwise
func (p *Proxiable) IsMethodPublicOnSelf(methodName string) bool {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return false
	}

	return reflector.HasMethod(self, methodName) && reflector.IsMethodPublic(self, methodName)
}

// GetMethodInfoOnSelf returns detailed information about a method on the self-reference.
// This uses our custom reflector to provide comprehensive method details.
//
// Parameters:
//   - methodName: Name of the method to inspect
//
// Returns:
//   - *reflector.MethodInfo: Detailed method information
//   - error: Any error that occurred during inspection
func (p *Proxiable) GetMethodInfoOnSelf(methodName string) (*reflector.MethodInfo, error) {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return nil, &ProxyError{
			Op:      "GetMethodInfoOnSelf",
			Method:  methodName,
			Message: "no self-reference set",
		}
	}

	return reflector.GetMethodInfo(self, methodName)
}

// GetReflectionResultOnSelf returns comprehensive reflection information about the self-reference.
// This uses our custom reflector to provide cached and detailed reflection data.
//
// Returns:
//   - *reflector.ReflectionResult: Comprehensive reflection information
func (p *Proxiable) GetReflectionResultOnSelf() *reflector.ReflectionResult {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return &reflector.ReflectionResult{IsValid: false}
	}

	return reflector.GetReflectionResult(self)
}

// GetAllMethodsOnSelf returns all methods available on the self-reference.
// This uses our custom reflector for consistent method enumeration.
//
// Returns:
//   - []reflect.Method: All methods on the self-reference
func (p *Proxiable) GetAllMethodsOnSelf() []reflect.Method {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return nil
	}

	return reflector.GetAllMethods(self)
}

// GetPublicMethodsOnSelf returns only the public (exported) methods on the self-reference.
// This uses our custom reflector for filtering exported methods.
//
// Returns:
//   - []reflect.Method: Public methods on the self-reference
func (p *Proxiable) GetPublicMethodsOnSelf() []reflect.Method {
	p.mutex.RLock()
	self := p.self
	p.mutex.RUnlock()

	if self == nil {
		return nil
	}

	return reflector.GetPublicMethods(self)
}

// ProxyError represents errors that can occur during proxy operations.
type ProxyError struct {
	Op      string // Operation that failed
	Method  string // Method name (if applicable)
	Type    string // Type name (if applicable)
	Message string // Error message
}

// Error implements the error interface.
func (e *ProxyError) Error() string {
	msg := e.Op
	if e.Method != "" {
		msg += "(" + e.Method + ")"
	}
	if e.Type != "" {
		msg += " on " + e.Type
	}
	msg += ": " + e.Message
	return msg
}
