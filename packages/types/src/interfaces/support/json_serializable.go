package interfaces

/**
 * JsonSerializableInterface defines the contract for custom JSON serialization.
 *
 * This interface provides a PHP/Laravel-inspired JsonSerializable pattern that allows
 * objects to control their JSON serialization behavior. Instead of relying on Go's
 * default JSON marshaling of struct fields, implementing objects can provide custom
 * serialization logic through the JsonSerialize method.
 *
 * Key Features:
 *   - Custom serialization logic for complex data structures
 *   - Control over which data gets included in JSON output
 *   - Ability to transform data before serialization
 *   - Support for nested JsonSerializable objects
 *   - Compatibility with Go's json package through custom MarshalJSON
 *
 * Example Usage:
 *   type User struct {
 *       ID       int    `json:"-"`        // Hidden from default JSON
 *       Name     string `json:"name"`
 *       Email    string `json:"email"`
 *       Password string `json:"-"`        // Hidden from default JSON
 *       IsActive bool   `json:"active"`
 *   }
 *
 *   func (u *User) JsonSerialize() interface{} {
 *       return map[string]interface{}{
 *           "id":     u.ID,
 *           "name":   u.Name,
 *           "email":  u.Email,
 *           "active": u.IsActive,
 *           // Password intentionally omitted for security
 *           "display_name": strings.Title(u.Name),  // Computed field
 *       }
 *   }
 *
 *   func (u *User) MarshalJSON() ([]byte, error) {
 *       return json.Marshal(u.JsonSerialize())
 *   }
 *
 *   // Usage:
 *   user := &User{ID: 1, Name: "john", Email: "john@example.com", Password: "secret"}
 *   jsonData, _ := json.Marshal(user)
 *   // Result: {"id":1,"name":"john","email":"john@example.com","active":false,"display_name":"John"}
 *
 * Design Patterns:
 *   - Facade Pattern: Provides a simplified interface for complex serialization
 *   - Template Method: Defines the serialization algorithm structure
 *   - Strategy Pattern: Allows different serialization strategies per type
 *
 * Security Considerations:
 *   - Always exclude sensitive data (passwords, tokens, etc.)
 *   - Validate and sanitize data before serialization
 *   - Be careful with recursive structures to avoid infinite loops
 *   - Consider rate limiting for expensive serialization operations
 *
 * Performance Considerations:
 *   - Cache serialized data when appropriate
 *   - Use lazy evaluation for expensive computed fields
 *   - Minimize memory allocations in hot paths
 *   - Consider streaming for large data sets
 *
 * Thread Safety:
 *   - Implementations should be safe for concurrent access
 *   - Consider read-only access patterns where possible
 *   - Use appropriate synchronization for mutable state
 */
type JsonSerializableInterface interface {
	/**
	 * JsonSerialize returns data which can be serialized by json.Marshal().
	 *
	 * This method should return a representation of the object that can be
	 * safely converted to JSON. The returned value will typically be passed
	 * to json.Marshal() to produce the final JSON output.
	 *
	 * Returns:
	 *   interface{} - A JSON-serializable representation of the object.
	 *                 Common return types include:
	 *                 - map[string]interface{} for object-like structures
	 *                 - []interface{} for array-like structures
	 *                 - Primitive types (string, int, bool, float64)
	 *                 - Nested structures implementing JsonSerializable
	 *
	 * Implementation Guidelines:
	 *
	 * 1. **Return JSON-Compatible Types**:
	 *    Only return types that json.Marshal can handle:
	 *    - Basic types: bool, int, float64, string
	 *    - Arrays and slices of JSON-compatible types
	 *    - Maps with string keys and JSON-compatible values
	 *    - Structs with proper JSON tags
	 *    - nil for null values
	 *
	 * 2. **Handle Nested Serialization**:
	 *    For nested objects that implement JsonSerializable:
	 *    if nested, ok := obj.(JsonSerializableInterface); ok {
	 *        return nested.JsonSerialize()
	 *    }
	 *    return obj  // Let json.Marshal handle it
	 *
	 * 3. **Security Best Practices**:
	 *    - Never include sensitive data (passwords, API keys, etc.)
	 *    - Validate data before including it
	 *    - Consider data visibility rules and user permissions
	 *    - Sanitize user-generated content
	 *
	 * 4. **Performance Optimization**:
	 *    - Cache expensive computations when appropriate
	 *    - Use lazy loading for optional fields
	 *    - Minimize memory allocations
	 *    - Consider pagination for large collections
	 *
	 * 5. **Error Handling**:
	 *    - Handle conversion errors gracefully
	 *    - Provide fallback values for problematic data
	 *    - Log errors for debugging purposes
	 *    - Don't panic; return safe defaults instead
	 *
	 * Example Implementations:
	 *
	 *   // Simple object serialization
	 *   func (p *Person) JsonSerialize() interface{} {
	 *       return map[string]interface{}{
	 *           "name":     p.Name,
	 *           "age":      p.Age,
	 *           "email":    p.Email,
	 *           "created":  p.CreatedAt.Format(time.RFC3339),
	 *       }
	 *   }
	 *
	 *   // Collection serialization
	 *   func (users *UserCollection) JsonSerialize() interface{} {
	 *       result := make([]interface{}, len(users.items))
	 *       for i, user := range users.items {
	 *           if serializable, ok := user.(JsonSerializableInterface); ok {
	 *               result[i] = serializable.JsonSerialize()
	 *           } else {
	 *               result[i] = user
	 *           }
	 *       }
	 *       return result
	 *   }
	 *
	 *   // Conditional serialization
	 *   func (u *User) JsonSerialize() interface{} {
	 *       data := map[string]interface{}{
	 *           "id":   u.ID,
	 *           "name": u.Name,
	 *       }
	 *       
	 *       // Only include email for verified users
	 *       if u.EmailVerified {
	 *           data["email"] = u.Email
	 *       }
	 *       
	 *       // Include admin fields for admin users
	 *       if u.IsAdmin {
	 *           data["permissions"] = u.Permissions
	 *           data["last_login"] = u.LastLogin
	 *       }
	 *       
	 *       return data
	 *   }
	 *
	 *   // Recursive structure handling
	 *   func (n *TreeNode) JsonSerialize() interface{} {
	 *       data := map[string]interface{}{
	 *           "id":    n.ID,
	 *           "value": n.Value,
	 *       }
	 *       
	 *       if len(n.Children) > 0 {
	 *           children := make([]interface{}, len(n.Children))
	 *           for i, child := range n.Children {
	 *               children[i] = child.JsonSerialize()
	 *           }
	 *           data["children"] = children
	 *       }
	 *       
	 *       return data
	 *   }
	 *
	 * Integration with Go's json package:
	 *   To integrate with Go's standard JSON marshaling, implement MarshalJSON:
	 *
	 *   func (obj *MyStruct) MarshalJSON() ([]byte, error) {
	 *       return json.Marshal(obj.JsonSerialize())
	 *   }
	 *
	 * Thread Safety:
	 *   - Should be safe for concurrent read operations
	 *   - Should not modify the object state during serialization
	 *   - Use appropriate synchronization for shared mutable data
	 *
	 * Testing Considerations:
	 *   - Test with various data combinations
	 *   - Verify security rules (no sensitive data leakage)
	 *   - Test recursive structures don't cause infinite loops
	 *   - Validate performance with large data sets
	 *   - Test error conditions and edge cases
	 */
	JsonSerialize() interface{}
}
