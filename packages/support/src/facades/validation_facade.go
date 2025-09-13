package facades

import (
	validationInterfaces "govel/types/src/interfaces/validation"
	facade "govel/support/src"
)

// Validation provides a clean, static-like interface to the application's input validation service.
//
// This facade implements the facade pattern, providing global access to the validation
// service configured in the dependency injection container. It offers a Laravel-style
// API for input validation with automatic service resolution, custom rules,
// comprehensive error handling, and flexible validation patterns.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved validation service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent validation operations across goroutines
//   - Supports custom rules, conditional validation, and nested data structures
//   - Built-in localization support and custom error message formatting
//
// Behavior:
//   - First call: Resolves validation service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if validation service cannot be resolved (fail-fast behavior)
//   - Automatically handles rule parsing, data validation, and error collection
//
// Returns:
//   - ValidationInterface: The application's validation service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "validation" service is not registered in the container
//   - If the resolved service doesn't implement ValidationInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Validation() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Validation operations are thread-safe with proper data isolation
//
// Usage Examples:
//
//	// Basic validation
//	data := map[string]interface{}{
//	    "name":     "John Doe",
//	    "email":    "john@example.com",
//	    "age":      25,
//	    "password": "password123",
//	}
//
//	rules := map[string]string{
//	    "name":     "required|string|min:2|max:50",
//	    "email":    "required|email|unique:users,email",
//	    "age":      "required|integer|min:18|max:120",
//	    "password": "required|string|min:8|confirmed",
//	}
//
//	// Validate data
//	validator := facades.Validation().Make(data, rules)
//	if validator.Fails() {
//	    errors := validator.Errors()
//	    for field, messages := range errors {
//	        for _, message := range messages {
//	            fmt.Printf("%s: %s\n", field, message)
//	        }
//	    }
//	} else {
//	    fmt.Println("Validation passed!")
//	}
//
//	// Struct validation with tags
//	type User struct {
//	    Name     string `validate:"required,min=2,max=50"`
//	    Email    string `validate:"required,email"`
//	    Age      int    `validate:"required,min=18,max=120"`
//	    Password string `validate:"required,min=8"`
//	}
//
//	user := User{
//	    Name:     "John Doe",
//	    Email:    "invalid-email",
//	    Age:      17,
//	    Password: "short",
//	}
//
//	// Validate struct
//	err := facades.Validation().Struct(&user)
//	if err != nil {
//	    if validationErrors, ok := err.(ValidationErrors); ok {
//	        for _, fieldError := range validationErrors {
//	            fmt.Printf("%s: %s\n", fieldError.Field(), fieldError.Message())
//	        }
//	    }
//	}
//
//	// Custom validation rules
//	// Register custom rule
//	facades.Validation().RegisterRule("phone", func(value interface{}) bool {
//	    str, ok := value.(string)
//	    if !ok {
//	        return false
//	    }
//
//	    // Basic phone number validation (adjust regex as needed)
//	    phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
//	    return phoneRegex.MatchString(str)
//	}, "The :attribute must be a valid phone number.")
//
//	// Use custom rule
//	contactData := map[string]interface{}{
//	    "phone": "+1234567890",
//	}
//
//	contactRules := map[string]string{
//	    "phone": "required|phone",
//	}
//
//	validator = facades.Validation().Make(contactData, contactRules)
//	if validator.Passes() {
//	    fmt.Println("Phone number is valid")
//	}
//
//	// Conditional validation
//	formData := map[string]interface{}{
//	    "account_type": "premium",
//	    "credit_limit": 5000,
//	    "referral_code": "FRIEND123",
//	}
//
//	conditionalRules := map[string]string{
//	    "account_type": "required|in:basic,premium,enterprise",
//	    "credit_limit": "required_if:account_type,premium|integer|min:1000",
//	    "referral_code": "sometimes|string|exists:referral_codes,code",
//	}
//
//	validator = facades.Validation().Make(formData, conditionalRules)
//	if validator.Fails() {
//	    fmt.Println("Conditional validation failed")
//	}
//
// Advanced Validation Patterns:
//
//	// Nested data validation
//	nestedData := map[string]interface{}{
//	    "user": map[string]interface{}{
//	        "name":    "John Doe",
//	        "email":   "john@example.com",
//	        "profile": map[string]interface{}{
//	            "bio":      "Software Developer",
//	            "website":  "https://johndoe.com",
//	            "location": "New York",
//	        },
//	    },
//	    "preferences": map[string]interface{}{
//	        "notifications": true,
//	        "theme":        "dark",
//	        "language":     "en",
//	    },
//	}
//
//	nestedRules := map[string]string{
//	    "user.name":              "required|string|max:100",
//	    "user.email":             "required|email",
//	    "user.profile.bio":       "sometimes|string|max:500",
//	    "user.profile.website":   "sometimes|url",
//	    "user.profile.location":  "sometimes|string|max:100",
//	    "preferences.theme":      "required|in:light,dark",
//	    "preferences.language":   "required|string|size:2",
//	}
//
//	validator = facades.Validation().Make(nestedData, nestedRules)
//
//	// Array validation
//	arrayData := map[string]interface{}{
//	    "tags": []string{"golang", "web", "api", "backend"},
//	    "items": []map[string]interface{}{
//	        {"name": "Item 1", "price": 19.99, "category": "electronics"},
//	        {"name": "Item 2", "price": 29.99, "category": "books"},
//	        {"name": "Item 3", "price": 39.99, "category": "clothing"},
//	    },
//	}
//
//	arrayRules := map[string]string{
//	    "tags":           "required|array|min:1|max:10",
//	    "tags.*":         "string|max:50",
//	    "items":          "required|array|min:1",
//	    "items.*.name":   "required|string|max:100",
//	    "items.*.price":  "required|numeric|min:0",
//	    "items.*.category": "required|in:electronics,books,clothing,home",
//	}
//
//	validator = facades.Validation().Make(arrayData, arrayRules)
//
//	// File validation
//	func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
//	    err := r.ParseMultipartForm(10 << 20) // 10MB limit
//	    if err != nil {
//	        http.Error(w, "Failed to parse form", http.StatusBadRequest)
//	        return
//	    }
//
//	    file, header, err := r.FormFile("upload")
//	    if err != nil {
//	        http.Error(w, "No file provided", http.StatusBadRequest)
//	        return
//	    }
//	    defer file.Close()
//
//	    // Validate file
//	    fileData := map[string]interface{}{
//	        "file": header,
//	        "name": r.FormValue("name"),
//	    }
//
//	    fileRules := map[string]string{
//	        "file": "required|file|mimes:jpeg,png,gif,pdf|max:5120", // 5MB
//	        "name": "required|string|max:100",
//	    }
//
//	    validator := facades.Validation().Make(fileData, fileRules)
//	    if validator.Fails() {
//	        errors := validator.Errors()
//	        w.Header().Set("Content-Type", "application/json")
//	        w.WriteHeader(http.StatusUnprocessableEntity)
//	        json.NewEncoder(w).Encode(map[string]interface{}{
//	            "errors": errors,
//	        })
//	        return
//	    }
//
//	    // File is valid, process upload
//	    // ... upload logic
//	}
//
// Custom Validation Rules:
//
//	// Complex custom rule with parameters
//	facades.Validation().RegisterRuleWithParams("between_dates",
//	    func(value interface{}, params []string) bool {
//	        if len(params) != 2 {
//	            return false
//	        }
//
//	        dateStr, ok := value.(string)
//	        if !ok {
//	            return false
//	        }
//
//	        date, err := time.Parse("2006-01-02", dateStr)
//	        if err != nil {
//	            return false
//	        }
//
//	        startDate, err := time.Parse("2006-01-02", params[0])
//	        if err != nil {
//	            return false
//	        }
//
//	        endDate, err := time.Parse("2006-01-02", params[1])
//	        if err != nil {
//	            return false
//	        }
//
//	        return date.After(startDate) && date.Before(endDate)
//	    },
//	    "The :attribute must be between :param1 and :param2."
//	)
//
//	// Use parameterized custom rule
//	eventData := map[string]interface{}{
//	    "event_date": "2023-07-15",
//	}
//
//	eventRules := map[string]string{
//	    "event_date": "required|date|between_dates:2023-06-01,2023-08-31",
//	}
//
//	// Database validation rule
//	facades.Validation().RegisterRule("unique_username",
//	    func(value interface{}) bool {
//	        username, ok := value.(string)
//	        if !ok {
//	            return false
//	        }
//
//	        var count int64
//	        err := facades.ORM().Model(&User{}).Where("username = ?", username).Count(&count).Error
//	        if err != nil {
//	            return false
//	        }
//
//	        return count == 0
//	    },
//	    "The :attribute has already been taken."
//	)
//
// Form Validation Patterns:
//
//	// Registration form validation
//	type RegistrationRequest struct {
//	    FirstName            string `json:"first_name" validate:"required,min=2,max=50"`
//	    LastName             string `json:"last_name" validate:"required,min=2,max=50"`
//	    Email                string `json:"email" validate:"required,email"`
//	    Password             string `json:"password" validate:"required,min=8,containsany=!@#$%^&*"`
//	    PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
//	    DateOfBirth          string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
//	    PhoneNumber          string `json:"phone_number" validate:"omitempty,phone"`
//	    TermsAccepted        bool   `json:"terms_accepted" validate:"required,eq=true"`
//	}
//
//	func ValidateRegistration(w http.ResponseWriter, r *http.Request) {
//	    var req RegistrationRequest
//
//	    // Parse JSON request
//	    err := json.NewDecoder(r.Body).Decode(&req)
//	    if err != nil {
//	        http.Error(w, "Invalid JSON", http.StatusBadRequest)
//	        return
//	    }
//
//	    // Validate struct
//	    err = facades.Validation().Struct(&req)
//	    if err != nil {
//	        validationErrors := facades.Validation().FormatErrors(err)
//
//	        w.Header().Set("Content-Type", "application/json")
//	        w.WriteHeader(http.StatusUnprocessableEntity)
//	        json.NewEncoder(w).Encode(map[string]interface{}{
//	            "message": "Validation failed",
//	            "errors":  validationErrors,
//	        })
//	        return
//	    }
//
//	    // Validation passed, process registration
//	    // ... registration logic
//	}
//
//	// API endpoint validation
//	type ProductCreateRequest struct {
//	    Name        string   `json:"name" validate:"required,min=1,max=200"`
//	    Description string   `json:"description" validate:"max=1000"`
//	    Price       float64  `json:"price" validate:"required,min=0"`
//	    CategoryID  int      `json:"category_id" validate:"required,min=1"`
//	    Tags        []string `json:"tags" validate:"max=10,dive,min=1,max=50"`
//	    InStock     bool     `json:"in_stock"`
//	    Weight      float64  `json:"weight" validate:"omitempty,min=0"`
//	}
//
//	func CreateProduct(w http.ResponseWriter, r *http.Request) {
//	    var req ProductCreateRequest
//
//	    // Parse and validate request
//	    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//	        http.Error(w, "Invalid JSON", http.StatusBadRequest)
//	        return
//	    }
//
//	    // Additional business logic validation
//	    businessRules := map[string]interface{}{
//	        "category_id": req.CategoryID,
//	        "name":        req.Name,
//	    }
//
//	    businessValidation := map[string]string{
//	        "category_id": "exists:categories,id",
//	        "name":        "unique:products,name",
//	    }
//
//	    // Validate struct first
//	    if err := facades.Validation().Struct(&req); err != nil {
//	        respondWithValidationErrors(w, err)
//	        return
//	    }
//
//	    // Then validate business rules
//	    validator := facades.Validation().Make(businessRules, businessValidation)
//	    if validator.Fails() {
//	        w.Header().Set("Content-Type", "application/json")
//	        w.WriteHeader(http.StatusUnprocessableEntity)
//	        json.NewEncoder(w).Encode(map[string]interface{}{
//	            "message": "Business validation failed",
//	            "errors":  validator.Errors(),
//	        })
//	        return
//	    }
//
//	    // Create product
//	    product := Product{
//	        Name:        req.Name,
//	        Description: req.Description,
//	        Price:       req.Price,
//	        CategoryID:  req.CategoryID,
//	        InStock:     req.InStock,
//	        Weight:      req.Weight,
//	    }
//
//	    err := facades.ORM().Create(&product).Error
//	    if err != nil {
//	        http.Error(w, "Failed to create product", http.StatusInternalServerError)
//	        return
//	    }
//
//	    w.Header().Set("Content-Type", "application/json")
//	    w.WriteHeader(http.StatusCreated)
//	    json.NewEncoder(w).Encode(product)
//	}
//
// Localized Validation Messages:
//
//	// Custom messages for specific fields
//	customMessages := map[string]string{
//	    "email.required":    "Please provide your email address.",
//	    "email.email":       "Please provide a valid email address.",
//	    "password.required": "Password is required.",
//	    "password.min":      "Password must be at least 8 characters long.",
//	    "age.min":           "You must be at least 18 years old.",
//	}
//
//	validator = facades.Validation().Make(data, rules)
//	validator.SetCustomMessages(customMessages)
//
//	if validator.Fails() {
//	    errors := validator.Errors()
//	    // Errors will use custom messages where defined
//	}
//
//	// Attribute name customization
//	attributeNames := map[string]string{
//	    "first_name": "First Name",
//	    "last_name":  "Last Name",
//	    "dob":        "Date of Birth",
//	}
//
//	validator = facades.Validation().Make(data, rules)
//	validator.SetAttributeNames(attributeNames)
//
//	// Language-specific validation
//	facades.Validation().SetLocale("es") // Spanish validation messages
//	validator = facades.Validation().Make(data, rules)
//
//	// Load custom language file
//	facades.Validation().LoadLanguageFile("validation", "./lang/es/validation.json")
//
// Conditional and Contextual Validation:
//
//	// Validation with context
//	type UpdateProfileContext struct {
//	    UserID      int
//	    CurrentUser *User
//	    IsAdmin     bool
//	}
//
//	func ValidateProfileUpdate(data map[string]interface{}, ctx UpdateProfileContext) error {
//	    rules := map[string]string{
//	        "name":  "required|string|max:100",
//	        "email": fmt.Sprintf("required|email|unique:users,email,%d", ctx.UserID),
//	        "bio":   "sometimes|string|max:500",
//	    }
//
//	    // Admin can update additional fields
//	    if ctx.IsAdmin {
//	        rules["role"] = "sometimes|in:user,admin,moderator"
//	        rules["verified"] = "sometimes|boolean"
//	    }
//
//	    validator := facades.Validation().Make(data, rules)
//
//	    // Add context-specific custom rules
//	    validator.After(func(v *Validator) {
//	        if email, ok := data["email"].(string); ok {
//	            if email != ctx.CurrentUser.Email {
//	                // Check if email change is allowed
//	                if !ctx.CurrentUser.CanChangeEmail() {
//	                    v.AddError("email", "You cannot change your email address.")
//	                }
//	            }
//	        }
//	    })
//
//	    if validator.Fails() {
//	        return &ValidationError{
//	            Errors: validator.Errors(),
//	        }
//	    }
//
//	    return nil
//	}
//
// Batch Validation:
//
//	// Validate multiple records
//	type BatchUserImport struct {
//	    Users []UserImportRecord `json:"users" validate:"required,min=1,max=1000,dive"`
//	}
//
//	type UserImportRecord struct {
//	    Name     string `json:"name" validate:"required,min=2,max=100"`
//	    Email    string `json:"email" validate:"required,email"`
//	    Role     string `json:"role" validate:"required,oneof=user admin moderator"`
//	    Active   bool   `json:"active"`
//	    Line     int    `json:"line,omitempty"` // For error reporting
//	}
//
//	func ValidateBatchImport(w http.ResponseWriter, r *http.Request) {
//	    var batch BatchUserImport
//
//	    if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
//	        http.Error(w, "Invalid JSON", http.StatusBadRequest)
//	        return
//	    }
//
//	    // Add line numbers for error reporting
//	    for i := range batch.Users {
//	        batch.Users[i].Line = i + 1
//	    }
//
//	    // Validate entire batch
//	    err := facades.Validation().Struct(&batch)
//	    if err != nil {
//	        validationErrors := facades.Validation().FormatBatchErrors(err, "users")
//
//	        w.Header().Set("Content-Type", "application/json")
//	        w.WriteHeader(http.StatusUnprocessableEntity)
//	        json.NewEncoder(w).Encode(map[string]interface{}{
//	            "message":       "Batch validation failed",
//	            "errors":        validationErrors,
//	            "valid_count":   len(batch.Users) - len(validationErrors),
//	            "invalid_count": len(validationErrors),
//	        })
//	        return
//	    }
//
//	    // Additional business validation for duplicates
//	    emailMap := make(map[string][]int)
//	    for i, user := range batch.Users {
//	        emailMap[user.Email] = append(emailMap[user.Email], i+1)
//	    }
//
//	    duplicateErrors := []map[string]interface{}{}
//	    for email, lines := range emailMap {
//	        if len(lines) > 1 {
//	            for _, line := range lines {
//	                duplicateErrors = append(duplicateErrors, map[string]interface{}{
//	                    "line":    line,
//	                    "field":   "email",
//	                    "value":   email,
//	                    "message": "Duplicate email address in batch",
//	                })
//	            }
//	        }
//	    }
//
//	    if len(duplicateErrors) > 0 {
//	        w.Header().Set("Content-Type", "application/json")
//	        w.WriteHeader(http.StatusUnprocessableEntity)
//	        json.NewEncoder(w).Encode(map[string]interface{}{
//	            "message": "Duplicate entries found",
//	            "errors":  duplicateErrors,
//	        })
//	        return
//	    }
//
//	    // Process valid batch
//	    // ... import logic
//	}
//
// Performance Optimizations:
//
//	// Cached validation rules
//	type CachedValidator struct {
//	    rules map[string]map[string]string
//	    mutex sync.RWMutex
//	}
//
//	func (cv *CachedValidator) GetRules(ruleSet string) map[string]string {
//	    cv.mutex.RLock()
//	    defer cv.mutex.RUnlock()
//
//	    if rules, exists := cv.rules[ruleSet]; exists {
//	        return rules
//	    }
//
//	    return nil
//	}
//
//	func (cv *CachedValidator) SetRules(ruleSet string, rules map[string]string) {
//	    cv.mutex.Lock()
//	    defer cv.mutex.Unlock()
//
//	    if cv.rules == nil {
//	        cv.rules = make(map[string]map[string]string)
//	    }
//
//	    cv.rules[ruleSet] = rules
//	}
//
//	// Use cached validator
//	cachedValidator := &CachedValidator{}
//
//	// Cache common validation rules
//	userRules := map[string]string{
//	    "name":  "required|string|max:100",
//	    "email": "required|email|unique:users,email",
//	    "age":   "required|integer|min:18",
//	}
//	cachedValidator.SetRules("user", userRules)
//
//	// Use cached rules
//	rules := cachedValidator.GetRules("user")
//	validator := facades.Validation().Make(userData, rules)
//
// Testing Validation:
//
//	// Test validation rules
//	func TestUserValidation(t *testing.T) {
//	    testCases := []struct {
//	        name        string
//	        data        map[string]interface{}
//	        shouldPass  bool
//	        expectedErrors []string
//	    }{
//	        {
//	            name: "valid user data",
//	            data: map[string]interface{}{
//	                "name":  "John Doe",
//	                "email": "john@example.com",
//	                "age":   25,
//	            },
//	            shouldPass: true,
//	        },
//	        {
//	            name: "missing required fields",
//	            data: map[string]interface{}{
//	                "name": "John Doe",
//	            },
//	            shouldPass: false,
//	            expectedErrors: []string{"email", "age"},
//	        },
//	        {
//	            name: "invalid email",
//	            data: map[string]interface{}{
//	                "name":  "John Doe",
//	                "email": "invalid-email",
//	                "age":   25,
//	            },
//	            shouldPass: false,
//	            expectedErrors: []string{"email"},
//	        },
//	    }
//
//	    rules := map[string]string{
//	        "name":  "required|string|max:100",
//	        "email": "required|email",
//	        "age":   "required|integer|min:18",
//	    }
//
//	    for _, tc := range testCases {
//	        t.Run(tc.name, func(t *testing.T) {
//	            validator := facades.Validation().Make(tc.data, rules)
//
//	            if tc.shouldPass {
//	                assert.True(t, validator.Passes(), "Validation should pass")
//	            } else {
//	                assert.True(t, validator.Fails(), "Validation should fail")
//
//	                errors := validator.Errors()
//	                for _, expectedError := range tc.expectedErrors {
//	                    assert.Contains(t, errors, expectedError, "Should contain error for %s", expectedError)
//	                }
//	            }
//	        })
//	    }
//	}
//
// Best Practices:
//   - Always validate user input at API boundaries
//   - Use struct tags for consistent validation rules
//   - Implement custom rules for business-specific validation
//   - Provide clear, user-friendly error messages
//   - Cache validation rules for better performance
//   - Use conditional validation for complex business logic
//   - Validate file uploads for security
//   - Implement proper error handling and response formatting
//   - Test validation rules thoroughly
//   - Use localized messages for international applications
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume validation service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	validation, err := facade.TryResolve[ValidationInterface]("validation")
//	if err != nil {
//	    // Handle validation service unavailability gracefully
//	    log.Printf("Validation service unavailable: %v", err)
//	    return // Skip validation operations
//	}
//	validator := validation.Make(data, rules)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestValidationBehavior(t *testing.T) {
//	    // Create a test validation service
//	    testValidation := &TestValidation{
//	        rules: make(map[string]string),
//	    }
//
//	    // Swap the real validation with test validation
//	    restore := support.SwapService("validation", testValidation)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Validation() returns testValidation
//	    validator := facades.Validation().Make(data, rules)
//
//	    // Verify validation behavior
//	    assert.NotNil(t, validator)
//	}
//
// Container Configuration:
// Ensure the validation service is properly configured in your container:
//
//	// Example validation registration
//	container.Singleton("validation", func() interface{} {
//	    config := validation.Config{
//	        // Default locale for validation messages
//	        DefaultLocale: "en",
//
//	        // Custom validation rules
//	        CustomRules: map[string]validation.RuleFunc{
//	            "phone": validation.PhoneRule,
//	            "slug":  validation.SlugRule,
//	        },
//
//	        // Database connection for database validation rules
//	        Database: facades.ORM(),
//
//	        // Language files path
//	        LanguagePath: "./lang",
//
//	        // Supported locales
//	        SupportedLocales: []string{"en", "es", "fr", "de"},
//
//	        // Custom message templates
//	        MessageTemplates: map[string]string{
//	            "required": "The :attribute field is required.",
//	            "email":    "The :attribute must be a valid email address.",
//	            "min":      "The :attribute must be at least :min characters.",
//	            "max":      "The :attribute may not be greater than :max characters.",
//	        },
//
//	        // Performance settings
//	        CacheRules:     true,
//	        CacheMessages:  true,
//	        CacheTimeout:   time.Hour,
//
//	        // File validation settings
//	        FileValidation: validation.FileConfig{
//	            MaxSize:      10 * 1024 * 1024, // 10MB
//	            AllowedMimes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
//	            ScanViruses:   true,
//	        },
//	    }
//
//	    validationService, err := validation.NewValidationService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create validation service: %v", err)
//	    }
//
//	    return validationService
//	})
func Validation() validationInterfaces.ValidationInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "validation" service from the dependency injection container
	// - Performs type assertion to ValidationInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[validationInterfaces.ValidationInterface](validationInterfaces.VALIDATION_TOKEN)
}

// ValidationWithError provides error-safe access to the validation service.
//
// This function offers the same functionality as Validation() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle validation service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Validation() but with error handling.
//
// Returns:
//   - ValidationInterface: The resolved validation instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement ValidationInterface
//
// Usage Examples:
//
//	// Basic error-safe validation
//	validation, err := facades.ValidationWithError()
//	if err != nil {
//	    log.Printf("Validation service unavailable: %v", err)
//	    return // Skip validation operations
//	}
//	validator := validation.Make(data, rules)
//
//	// Conditional validation
//	if validation, err := facades.ValidationWithError(); err == nil {
//	    // Perform optional input validation
//	    validator := validation.Make(optionalData, rules)
//	    if validator.Fails() {
//	        log.Printf("Optional validation failed: %v", validator.Errors())
//	    }
//	}
func ValidationWithError() (validationInterfaces.ValidationInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "validation" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[validationInterfaces.ValidationInterface](validationInterfaces.VALIDATION_TOKEN)
}
