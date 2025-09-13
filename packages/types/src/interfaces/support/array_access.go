package interfaces

/**
 * ArrayAccessInterface defines the contract for array-like access to data structures.
 *
 * This interface provides a PHP/Laravel-inspired ArrayAccess pattern that allows
 * objects to be accessed like arrays using bracket notation. It enables treating
 * objects as associative arrays or indexed collections while maintaining type safety
 * and custom access logic.
 *
 * The interface supports:
 *   - Checking if an offset/key exists (OffsetExists)
 *   - Getting values by offset/key (OffsetGet)
 *   - Setting values by offset/key (OffsetSet)
 *   - Unsetting/removing values by offset/key (OffsetUnset)
 *
 * Example Usage:
 *   type MyCollection struct {
 *       data map[string]interface{}
 *   }
 *
 *   func (c *MyCollection) OffsetExists(offset interface{}) bool {
 *       key, ok := offset.(string)
 *       if !ok { return false }
 *       _, exists := c.data[key]
 *       return exists
 *   }
 *
 *   // Usage:
 *   collection := &MyCollection{data: make(map[string]interface{})}
 *   exists := collection.OffsetExists("name")  // Check existence
 *   value := collection.OffsetGet("name")     // Get value
 *   collection.OffsetSet("name", "John")     // Set value
 *   collection.OffsetUnset("name")           // Remove value
 *
 * Design Considerations:
 *   - Uses interface{} for maximum flexibility in key/value types
 *   - Implementations should handle type conversion and validation
 *   - Should gracefully handle invalid or missing keys
 *   - Maintains consistency with PHP ArrayAccess behavior
 *
 * Thread Safety:
 *   Implementations must handle concurrent access appropriately.
 *   Consider using mutexes for thread-safe implementations.
 *
 * Performance:
 *   Methods should be optimized for frequent access patterns.
 *   Consider caching and lazy evaluation where appropriate.
 */
type ArrayAccessInterface interface {
	/**
	 * OffsetExists checks whether an offset exists in the data structure.
	 *
	 * This method determines if a given key/offset exists within the implementing
	 * data structure. It should return true if the offset exists, regardless of
	 * whether the associated value is null, empty, or falsy.
	 *
	 * Parameters:
	 *   offset interface{} - The offset/key to check for existence.
	 *                       Can be string, int, or any comparable type.
	 *
	 * Returns:
	 *   bool - True if the offset exists, false otherwise.
	 *
	 * Implementation Notes:
	 *   - Should handle type conversion gracefully (e.g., int to string keys)
	 *   - Should return false for invalid/unsupported key types
	 *   - Should be consistent with the data structure's key semantics
	 *   - Should not modify the data structure during check
	 *
	 * Example Implementations:
	 *   // For map-based structures
	 *   key, ok := offset.(string)
	 *   if !ok { return false }
	 *   _, exists := data[key]
	 *   return exists
	 *
	 *   // For slice-based structures
	 *   index, ok := offset.(int)
	 *   if !ok { return false }
	 *   return index >= 0 && index < len(slice)
	 *
	 * Thread Safety:
	 *   Should be safe for concurrent reads if the underlying structure allows it.
	 */
	OffsetExists(offset interface{}) bool

	/**
	 * OffsetGet retrieves the value at the specified offset.
	 *
	 * This method returns the value associated with the given key/offset.
	 * If the offset doesn't exist, it should return nil or a suitable default
	 * value based on the implementation's semantics.
	 *
	 * Parameters:
	 *   offset interface{} - The offset/key to retrieve the value for.
	 *                       Should match the key type used in OffsetExists.
	 *
	 * Returns:
	 *   interface{} - The value at the specified offset, or nil if not found.
	 *
	 * Implementation Notes:
	 *   - Should handle type conversion consistently with OffsetExists
	 *   - Should return nil for non-existent or invalid keys
	 *   - May return typed zero values instead of nil for specific implementations
	 *   - Should not panic on invalid keys, return nil instead
	 *
	 * Example Implementations:
	 *   // For map-based structures
	 *   key, ok := offset.(string)
	 *   if !ok { return nil }
	 *   return data[key]  // Returns nil if key doesn't exist
	 *
	 *   // For slice-based structures with bounds checking
	 *   index, ok := offset.(int)
	 *   if !ok || index < 0 || index >= len(slice) { return nil }
	 *   return slice[index]
	 *
	 * Error Handling:
	 *   - Should not panic on invalid offsets
	 *   - Should return nil for out-of-bounds or invalid keys
	 *   - May log warnings for invalid access attempts
	 *
	 * Thread Safety:
	 *   Should be safe for concurrent reads if the underlying structure supports it.
	 */
	OffsetGet(offset interface{}) interface{}

	/**
	 * OffsetSet assigns a value to the specified offset.
	 *
	 * This method sets the value at the given key/offset. If the offset already
	 * exists, it should update the existing value. If the offset doesn't exist,
	 * it should create a new entry.
	 *
	 * Parameters:
	 *   offset interface{} - The offset/key where to assign the value.
	 *                       For append operations, may be nil.
	 *   value interface{}  - The value to assign at the specified offset.
	 *
	 * Behavior:
	 *   - Creates new entries for non-existent offsets
	 *   - Updates existing entries when offset already exists
	 *   - May append to collections when offset is nil (implementation-specific)
	 *   - Should handle type conversion and validation
	 *
	 * Implementation Notes:
	 *   - Should handle nil offset appropriately (e.g., append for slices)
	 *   - Should validate key and value types as needed
	 *   - May resize or reallocate underlying storage as needed
	 *   - Should maintain data structure invariants
	 *
	 * Example Implementations:
	 *   // For map-based structures
	 *   key, ok := offset.(string)
	 *   if ok {
	 *       data[key] = value
	 *   }
	 *
	 *   // For slice-based structures
	 *   if offset == nil {
	 *       slice = append(slice, value)  // Append mode
	 *   } else if index, ok := offset.(int); ok && index >= 0 && index < len(slice) {
	 *       slice[index] = value  // Update mode
	 *   }
	 *
	 * Error Handling:
	 *   - May panic for immutable structures (e.g., strings)
	 *   - Should validate input types and ranges
	 *   - May ignore invalid assignments silently or log warnings
	 *
	 * Thread Safety:
	 *   Implementations should handle concurrent writes appropriately.
	 *   Consider using mutexes for thread-safe modifications.
	 */
	OffsetSet(offset, value interface{})

	/**
	 * OffsetUnset removes the entry at the specified offset.
	 *
	 * This method removes the key/value pair at the given offset from the
	 * data structure. After this operation, OffsetExists should return false
	 * for the same offset.
	 *
	 * Parameters:
	 *   offset interface{} - The offset/key to remove from the data structure.
	 *
	 * Behavior:
	 *   - Removes the entry if it exists
	 *   - Should be idempotent (safe to call on non-existent keys)
	 *   - May compact or reorganize the underlying storage
	 *   - Should maintain data structure consistency
	 *
	 * Implementation Notes:
	 *   - Should handle type conversion consistently with other methods
	 *   - Should be safe to call on non-existent keys (no-op)
	 *   - May trigger garbage collection of removed values
	 *   - Should update internal size/count tracking
	 *
	 * Example Implementations:
	 *   // For map-based structures
	 *   key, ok := offset.(string)
	 *   if ok {
	 *       delete(data, key)
	 *   }
	 *
	 *   // For slice-based structures (removing by index)
	 *   index, ok := offset.(int)
	 *   if ok && index >= 0 && index < len(slice) {
	 *       slice = append(slice[:index], slice[index+1:]...)
	 *   }
	 *
	 * Error Handling:
	 *   - May panic for immutable structures
	 *   - Should handle out-of-bounds indices gracefully
	 *   - Should ignore invalid key types silently
	 *
	 * Thread Safety:
	 *   Implementations should handle concurrent modifications appropriately.
	 *   Consider using mutexes for thread-safe removals.
	 */
	OffsetUnset(offset interface{})
}
