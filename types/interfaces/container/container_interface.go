package container

// ContainerInterface defines the interface for dependency injection container.
type ContainerInterface interface {
	// Bind registers a binding in the container.
	Bind(abstract string, concrete interface{}, shared ...bool) error

	// BindIf registers a binding in the container if it doesn't already exist.
	BindIf(abstract string, concrete interface{}, shared ...bool) error

	// Singleton registers a shared binding in the container.
	Singleton(abstract string, concrete interface{}) error

	// Instance registers an existing instance as shared in the container.
	Instance(abstract string, instance interface{}) error

	// Make resolves the given type from the container.
	Make(abstract string) (interface{}, error)

	// MakeWith resolves the given type from the container with parameters.
	MakeWith(abstract string, parameters map[string]interface{}) (interface{}, error)

	// Bound checks if the given abstract type has been bound.
	Bound(abstract string) bool

	// Resolved checks if the given abstract type has been resolved.
	Resolved(abstract string) bool

	// Forget removes a binding from the container.
	Forget(abstract string) error

	// Flush flushes all bindings and resolved instances from the container.
	Flush() error

	// Tag adds a tag to a binding.
	Tag(abstracts []string, tag string) error

	// Tagged resolves all bindings for a given tag.
	Tagged(tag string) ([]interface{}, error)

	// When registers a resolver for when a given type is requested.
	When(concrete string) ConditionalBinding

	// Extend extends an abstract type with a closure.
	Extend(abstract string, closure func(interface{}, ContainerInterface) interface{}) error
}

// ConditionalBinding interface for conditional bindings.
type ConditionalBinding interface {
	// Needs specifies what the binding needs.
	Needs(abstract string) ContextualBinding
}

// ContextualBinding interface for contextual bindings.
type ContextualBinding interface {
	// Give specifies what to give when the binding is requested.
	Give(concrete interface{}) error

	// GiveTagged specifies to give all tagged instances when the binding is requested.
	GiveTagged(tag string) error
}