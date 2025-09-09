package mocks

import (
	containerInterfaces "govel/packages/container/interfaces"
)

/**
 * MockContainer provides a mock implementation of ContainerInterface for testing.
 * This mock allows tests to verify dependency injection behavior without actual service resolution complexity.
 */
type MockContainer struct {
	// Service Bindings Storage
	Bindings map[string]interface{}

	// Singleton Instances Storage
	Singletons map[string]interface{}

	// Singleton Bindings (tracks which services should be singletons)
	SingletonBindings map[string]bool

	// Mock Control Flags
	ShouldFailBind      bool
	ShouldFailMake      bool
	ShouldFailSingleton bool

	// Operation History
	BindHistory      []BindOperation
	MakeHistory      []MakeOperation
	SingletonHistory []SingletonOperation
	ForgetHistory    []string
	FlushHistory     []string // timestamps or reasons
}

/**
 * BindOperation represents a bind operation for testing verification
 */
type BindOperation struct {
	Abstract string
	Concrete interface{}
	Success  bool
}

/**
 * MakeOperation represents a make operation for testing verification
 */
type MakeOperation struct {
	Abstract string
	Result   interface{}
	Success  bool
	Error    error
}

/**
 * SingletonOperation represents a singleton operation for testing verification
 */
type SingletonOperation struct {
	Abstract string
	Concrete interface{}
	Success  bool
}

/**
 * NewMockContainer creates a new mock container with default values
 */
func NewMockContainer() *MockContainer {
	return &MockContainer{
		Bindings:          make(map[string]interface{}),
		Singletons:        make(map[string]interface{}),
		SingletonBindings: make(map[string]bool),
		BindHistory:       make([]BindOperation, 0),
		MakeHistory:       make([]MakeOperation, 0),
		SingletonHistory:  make([]SingletonOperation, 0),
		ForgetHistory:     make([]string, 0),
		FlushHistory:      make([]string, 0),
	}
}

// ContainerInterface Implementation

func (m *MockContainer) Bind(abstract string, concrete interface{}) error {
	operation := BindOperation{
		Abstract: abstract,
		Concrete: concrete,
		Success:  !m.ShouldFailBind,
	}
	m.BindHistory = append(m.BindHistory, operation)

	if m.ShouldFailBind {
		return &MockContainerError{Message: "mock bind failure", Abstract: abstract}
	}

	m.Bindings[abstract] = concrete
	return nil
}

func (m *MockContainer) Singleton(abstract string, concrete interface{}) error {
	operation := SingletonOperation{
		Abstract: abstract,
		Concrete: concrete,
		Success:  !m.ShouldFailSingleton,
	}
	m.SingletonHistory = append(m.SingletonHistory, operation)

	if m.ShouldFailSingleton {
		return &MockContainerError{Message: "mock singleton failure", Abstract: abstract}
	}

	m.Bindings[abstract] = concrete
	m.SingletonBindings[abstract] = true
	return nil
}

func (m *MockContainer) Make(abstract string) (interface{}, error) {
	var result interface{}
	var err error
	success := !m.ShouldFailMake

	if m.ShouldFailMake {
		err = &MockContainerError{Message: "mock make failure", Abstract: abstract}
	} else {
		result, err = m.makeInternal(abstract)
		success = err == nil
	}

	operation := MakeOperation{
		Abstract: abstract,
		Result:   result,
		Success:  success,
		Error:    err,
	}
	m.MakeHistory = append(m.MakeHistory, operation)

	return result, err
}

func (m *MockContainer) makeInternal(abstract string) (interface{}, error) {
	// Check if it's a singleton that's already instantiated
	if m.SingletonBindings[abstract] {
		if instance, exists := m.Singletons[abstract]; exists {
			return instance, nil
		}
	}

	// Get the binding
	binding, exists := m.Bindings[abstract]
	if !exists {
		return nil, &MockContainerError{Message: "binding not found", Abstract: abstract}
	}

	// Resolve the binding
	var instance interface{}

	// If it's a function, call it
	if fn, ok := binding.(func() interface{}); ok {
		instance = fn()
	} else if fn, ok := binding.(func() (interface{}, error)); ok {
		var err error
		instance, err = fn()
		if err != nil {
			return nil, err
		}
	} else {
		// Return the binding as-is
		instance = binding
	}

	// If it's a singleton, store the instance
	if m.SingletonBindings[abstract] {
		m.Singletons[abstract] = instance
	}

	return instance, nil
}

func (m *MockContainer) IsBound(abstract string) bool {
	_, exists := m.Bindings[abstract]
	return exists
}

func (m *MockContainer) Forget(abstract string) {
	m.ForgetHistory = append(m.ForgetHistory, abstract)

	delete(m.Bindings, abstract)
	delete(m.Singletons, abstract)
	delete(m.SingletonBindings, abstract)
}

func (m *MockContainer) FlushContainer() {
	m.FlushHistory = append(m.FlushHistory, "flush_all")

	m.Bindings = make(map[string]interface{})
	m.Singletons = make(map[string]interface{})
	m.SingletonBindings = make(map[string]bool)
}

// Mock-specific helper methods

/**
 * GetBindings returns all current bindings
 */
func (m *MockContainer) GetBindings() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m.Bindings {
		result[k] = v
	}
	return result
}

/**
 * GetSingletons returns all current singleton instances
 */
func (m *MockContainer) GetSingletons() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m.Singletons {
		result[k] = v
	}
	return result
}

/**
 * GetSingletonBindings returns all singleton binding flags
 */
func (m *MockContainer) GetSingletonBindings() map[string]bool {
	result := make(map[string]bool)
	for k, v := range m.SingletonBindings {
		result[k] = v
	}
	return result
}

/**
 * GetBindHistory returns the history of bind operations
 */
func (m *MockContainer) GetBindHistory() []BindOperation {
	return m.BindHistory
}

/**
 * GetMakeHistory returns the history of make operations
 */
func (m *MockContainer) GetMakeHistory() []MakeOperation {
	return m.MakeHistory
}

/**
 * GetSingletonHistory returns the history of singleton operations
 */
func (m *MockContainer) GetSingletonHistory() []SingletonOperation {
	return m.SingletonHistory
}

/**
 * GetForgetHistory returns the history of forget operations
 */
func (m *MockContainer) GetForgetHistory() []string {
	return m.ForgetHistory
}

/**
 * GetFlushHistory returns the history of flush operations
 */
func (m *MockContainer) GetFlushHistory() []string {
	return m.FlushHistory
}

/**
 * SetFailureMode sets whether various operations should fail
 */
func (m *MockContainer) SetFailureMode(bind, make, singleton bool) {
	m.ShouldFailBind = bind
	m.ShouldFailMake = make
	m.ShouldFailSingleton = singleton
}

/**
 * ClearHistory clears all operation history
 */
func (m *MockContainer) ClearHistory() {
	m.BindHistory = make([]BindOperation, 0)
	m.MakeHistory = make([]MakeOperation, 0)
	m.SingletonHistory = make([]SingletonOperation, 0)
	m.ForgetHistory = make([]string, 0)
	m.FlushHistory = make([]string, 0)
}

/**
 * GetBindingCount returns the number of registered bindings
 */
func (m *MockContainer) GetBindingCount() int {
	return len(m.Bindings)
}

/**
 * GetSingletonCount returns the number of instantiated singletons
 */
func (m *MockContainer) GetSingletonCount() int {
	return len(m.Singletons)
}

/**
 * IsSingleton checks if a service is registered as a singleton
 */
func (m *MockContainer) IsSingleton(abstract string) bool {
	return m.SingletonBindings[abstract]
}

/**
 * HasSingletonInstance checks if a singleton instance exists
 */
func (m *MockContainer) HasSingletonInstance(abstract string) bool {
	_, exists := m.Singletons[abstract]
	return exists
}

// Mock Error Type
type MockContainerError struct {
	Message  string
	Abstract string
}

func (e *MockContainerError) Error() string {
	if e.Abstract != "" {
		return "mock container error (" + e.Abstract + "): " + e.Message
	}
	return "mock container error: " + e.Message
}

// Compile-time interface compliance check
var _ containerInterfaces.ContainerInterface = (*MockContainer)(nil)

/**
 * MockContainable provides a mock implementation of ContainableInterface for testing.
 */
type MockContainable struct {
	*MockContainer

	ContainerInstance containerInterfaces.ContainerInterface
	HasContainerValue bool
}

/**
 * NewMockContainable creates a new mock containable with default values
 */
func NewMockContainable() *MockContainable {
	mockContainer := NewMockContainer()
	return &MockContainable{
		MockContainer:     mockContainer,
		ContainerInstance: mockContainer,
		HasContainerValue: true,
	}
}

// ContainableInterface Implementation

func (m *MockContainable) GetContainer() containerInterfaces.ContainerInterface {
	return m.ContainerInstance
}

func (m *MockContainable) SetContainer(container interface{}) {
	if ctr, ok := container.(containerInterfaces.ContainerInterface); ok {
		m.ContainerInstance = ctr
		m.HasContainerValue = true
	} else if ctr, ok := container.(*MockContainer); ok {
		m.ContainerInstance = ctr
		m.HasContainerValue = true
	}
}

func (m *MockContainable) HasContainer() bool {
	return m.HasContainerValue
}

func (m *MockContainable) GetContainerInfo() map[string]interface{} {
	info := map[string]interface{}{
		"has_container":  m.HasContainerValue,
		"container_type": "mock",
	}

	if m.ContainerInstance != nil {
		if mockContainer, ok := m.ContainerInstance.(*MockContainer); ok {
			info["bindings_count"] = mockContainer.GetBindingCount()
			info["singletons_count"] = mockContainer.GetSingletonCount()
			info["bind_operations"] = len(mockContainer.GetBindHistory())
			info["make_operations"] = len(mockContainer.GetMakeHistory())
			info["singleton_operations"] = len(mockContainer.GetSingletonHistory())
			info["forget_operations"] = len(mockContainer.GetForgetHistory())
			info["flush_operations"] = len(mockContainer.GetFlushHistory())
		} else {
			// For non-mock containers, provide basic info
			info["bindings_count"] = "unknown"
			info["singletons_count"] = "unknown"
		}
	}

	return info
}

// Mock-specific helper methods for Containable

/**
 * SetHasContainer controls whether the containable reports having a container
 */
func (m *MockContainable) SetHasContainer(hasContainer bool) {
	m.HasContainerValue = hasContainer
}

/**
 * GetMockContainer returns the underlying MockContainer if available
 */
func (m *MockContainable) GetMockContainer() *MockContainer {
	if mockContainer, ok := m.ContainerInstance.(*MockContainer); ok {
		return mockContainer
	}
	return nil
}

// Compile-time interface compliance check
var _ containerInterfaces.ContainableInterface = (*MockContainable)(nil)
