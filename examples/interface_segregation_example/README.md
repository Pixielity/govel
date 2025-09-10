# Interface Segregation Example

This example demonstrates how the `When()` method has been made optional in service providers through proper interface segregation, following the Interface Segregation Principle (ISP).

## Problem

Previously, the `DeferrableProvider` interface required all implementing providers to define a `When()` method, even if they didn't need event-based triggering. This violated the ISP by forcing providers to implement functionality they didn't need.

```go
// OLD - Violates ISP
type DeferrableProvider interface {
    Provides() []string  // Required for deferred loading
    When() []string      // Not always needed, but was mandatory
}
```

## Solution

The solution separates concerns into two distinct interfaces:

1. **`DeferrableProvider`** - Core deferred loading functionality
2. **`EventTriggeredProvider`** - Optional event-based triggering

```go
// NEW - Follows ISP
type DeferrableProvider interface {
    Provides() []string  // Only what's needed for deferred loading
}

type EventTriggeredProvider interface {
    When() []string      // Only for providers that need event triggering
}
```

## Usage Examples

### Simple Provider (Deferred Only)

The `SimpleProvider` in this example implements only `DeferrableProvider`:

```go
type SimpleProvider struct {
    // ... fields
}

// Implements DeferrableProvider
func (p *SimpleProvider) Provides() []string {
    return []string{"simple.service", "simple.helper"}
}

// No When() method needed!

// Interface compliance
var _ providerInterfaces.DeferrableProvider = (*SimpleProvider)(nil)
// var _ providerInterfaces.EventTriggeredProvider = (*SimpleProvider)(nil) // Would fail
```

### Full-Featured Provider (Deferred + Event-Triggered)

Providers that need event triggering can implement both interfaces:

```go
type AdvancedProvider struct {
    // ... fields
}

// Implements DeferrableProvider
func (p *AdvancedProvider) Provides() []string {
    return []string{"advanced.service"}
}

// Implements EventTriggeredProvider
func (p *AdvancedProvider) When() []string {
    return []string{"booting", "request.started"}
}

// Interface compliance
var _ providerInterfaces.DeferrableProvider = (*AdvancedProvider)(nil)
var _ providerInterfaces.EventTriggeredProvider = (*AdvancedProvider)(nil)
```

## Runtime Detection

The provider manifest system automatically detects which interfaces are implemented:

```go
// Check for deferred loading capability
if provider.IsDeferred() {
    services := provider.GetProvides()
    // Register deferred services...
    
    // Check for optional event triggering
    if eventTriggered, ok := provider.(EventTriggeredProvider); ok {
        events := eventTriggered.When()
        if len(events) > 0 {
            // Register event triggers...
        }
    }
}
```

## Benefits

1. **ISP Compliance**: Providers only implement what they need
2. **Backward Compatibility**: Existing providers work without changes
3. **Flexibility**: New providers can choose their level of functionality
4. **Clear Separation**: Event triggering is clearly separated from deferred loading
5. **No Compilation Errors**: Providers without event needs compile successfully

## Files Changed

- `event_triggered_provider_interface.go` - New dedicated interface
- `deferrable_provider_interface.go` - Simplified to core functionality  
- `provider_manifest.go` - Updated to use new interface detection
- `interfaces.go` - Added new interface export
- Example providers - Added `When()` method implementations

## Running the Example

```bash
cd /Users/akouta/Projects/govel/examples/interface_segregation_example
go run simple_provider.go
```

This demonstrates a working provider that implements deferred loading without event triggering, showing that the `When()` method is now truly optional.
