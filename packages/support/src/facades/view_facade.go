package facades

import (
	viewInterfaces "govel/packages/types/src/interfaces/view"
	facade "govel/packages/support/src"
)

// View provides a clean, static-like interface to the application's template rendering service.
//
// This facade implements the facade pattern, providing global access to the view
// service configured in the dependency injection container. It offers a Laravel-style
// API for template rendering with automatic service resolution, multiple template engines,
// view composition, and comprehensive data binding capabilities.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved view service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent template rendering across goroutines
//   - Supports multiple template engines (HTML, Pug, Mustache, Go templates)
//   - Built-in view composition, layouts, partials, and component systems
//
// Behavior:
//   - First call: Resolves view service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if view service cannot be resolved (fail-fast behavior)
//   - Automatically handles template compilation, caching, and rendering
//
// Returns:
//   - ViewInterface: The application's view service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "view" service is not registered in the container
//   - If the resolved service doesn't implement ViewInterface
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
//   - Multiple goroutines can call View() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Template rendering is thread-safe with proper data isolation
//
// Usage Examples:
//
//	// Basic template rendering
//	func HandleHomePage(w http.ResponseWriter, r *http.Request) {
//	    data := map[string]interface{}{
//	        "title":   "Welcome to My Website",
//	        "user":    getCurrentUser(r),
//	        "posts":   getRecentPosts(),
//	        "version": "1.0.0",
//	    }
//
//	    html, err := facades.View().Render("home", data)
//	    if err != nil {
//	        http.Error(w, "Template rendering failed", http.StatusInternalServerError)
//	        return
//	    }
//
//	    w.Header().Set("Content-Type", "text/html")
//	    w.Write([]byte(html))
//	}
//
//	// Render template with layout
//	func HandleUserProfile(w http.ResponseWriter, r *http.Request) {
//	    userID := getUserID(r)
//	    user := getUserByID(userID)
//
//	    data := map[string]interface{}{
//	        "user":       user,
//	        "page_title": user.Name + "'s Profile",
//	        "breadcrumb": []string{"Home", "Users", user.Name},
//	    }
//
//	    // Render with layout
//	    html, err := facades.View().RenderWithLayout("users/profile", "layouts/app", data)
//	    if err != nil {
//	        http.Error(w, "Profile rendering failed", http.StatusInternalServerError)
//	        return
//	    }
//
//	    w.Header().Set("Content-Type", "text/html")
//	    w.Write([]byte(html))
//	}
//
//	// Stream large templates directly to response
//	func HandleLargeReport(w http.ResponseWriter, r *http.Request) {
//	    reportData := generateLargeReportData()
//
//	    data := map[string]interface{}{
//	        "report":     reportData,
//	        "generated": time.Now().Format("2006-01-02 15:04:05"),
//	        "user":      getCurrentUser(r),
//	    }
//
//	    w.Header().Set("Content-Type", "text/html")
//	    w.Header().Set("Content-Disposition", "attachment; filename=report.html")
//
//	    // Stream template directly to response for memory efficiency
//	    err := facades.View().RenderToWriter(w, "reports/large-report", data)
//	    if err != nil {
//	        http.Error(w, "Report generation failed", http.StatusInternalServerError)
//	    }
//	}
//
//	// JSON API with template rendering
//	func HandleAPIWithTemplate(w http.ResponseWriter, r *http.Request) {
//	    // Accept header determines response format
//	    acceptHeader := r.Header.Get("Accept")
//
//	    data := map[string]interface{}{
//	        "products": getProducts(),
//	        "total":    getTotalProducts(),
//	    }
//
//	    if strings.Contains(acceptHeader, "text/html") {
//	        // Render HTML template
//	        html, err := facades.View().Render("products/index", data)
//	        if err != nil {
//	            http.Error(w, "Template error", http.StatusInternalServerError)
//	            return
//	        }
//
//	        w.Header().Set("Content-Type", "text/html")
//	        w.Write([]byte(html))
//	    } else {
//	        // Return JSON
//	        w.Header().Set("Content-Type", "application/json")
//	        json.NewEncoder(w).Encode(data)
//	    }
//	}
//
// Advanced View Patterns:
//
//	// View composition with shared data
//	type ViewComposer struct {
//	    sharedData map[string]interface{}
//	}
//
//	func NewViewComposer() *ViewComposer {
//	    return &ViewComposer{
//	        sharedData: map[string]interface{}{
//	            "app_name":    "My Application",
//	            "app_version": "1.0.0",
//	            "year":        time.Now().Year(),
//	        },
//	    }
//	}
//
//	func (vc *ViewComposer) AddGlobalData(key string, value interface{}) {
//	    vc.sharedData[key] = value
//	}
//
//	func (vc *ViewComposer) Render(template string, data map[string]interface{}) (string, error) {
//	    // Merge shared data with page-specific data
//	    mergedData := make(map[string]interface{})
//
//	    // Add shared data first
//	    for k, v := range vc.sharedData {
//	        mergedData[k] = v
//	    }
//
//	    // Override with page-specific data
//	    for k, v := range data {
//	        mergedData[k] = v
//	    }
//
//	    return facades.View().Render(template, mergedData)
//	}
//
//	// Use composer
//	composer := NewViewComposer()
//	composer.AddGlobalData("current_user", getCurrentUser(r))
//	composer.AddGlobalData("nav_items", getNavigationItems())
//
//	html, err := composer.Render("dashboard", map[string]interface{}{
//	    "stats": getDashboardStats(),
//	    "alerts": getAlerts(),
//	})
//
//	// Partial templates and includes
//	func RenderProductCard(product Product) (string, error) {
//	    data := map[string]interface{}{
//	        "product": product,
//	        "currency": "$",
//	    }
//
//	    return facades.View().RenderPartial("partials/product-card", data)
//	}
//
//	func HandleProductListing(w http.ResponseWriter, r *http.Request) {
//	    products := getProducts()
//
//	    // Render individual product cards
//	    productCards := make([]string, 0, len(products))
//	    for _, product := range products {
//	        card, err := RenderProductCard(product)
//	        if err != nil {
//	            http.Error(w, "Card rendering failed", http.StatusInternalServerError)
//	            return
//	        }
//	        productCards = append(productCards, card)
//	    }
//
//	    // Render main template with product cards
//	    data := map[string]interface{}{
//	        "title": "Product Listing",
//	        "product_cards": productCards,
//	        "total_products": len(products),
//	    }
//
//	    html, err := facades.View().Render("products/listing", data)
//	    if err != nil {
//	        http.Error(w, "Template error", http.StatusInternalServerError)
//	        return
//	    }
//
//	    w.Header().Set("Content-Type", "text/html")
//	    w.Write([]byte(html))
//	}
//
//	// Component-based rendering
//	type Component interface {
//	    Render(data map[string]interface{}) (string, error)
//	    GetName() string
//	}
//
//	type ButtonComponent struct {
//	    template string
//	}
//
//	func (bc *ButtonComponent) GetName() string {
//	    return "button"
//	}
//
//	func (bc *ButtonComponent) Render(data map[string]interface{}) (string, error) {
//	    // Set default values
//	    if data["type"] == nil {
//	        data["type"] = "button"
//	    }
//	    if data["class"] == nil {
//	        data["class"] = "btn btn-primary"
//	    }
//
//	    return facades.View().RenderPartial("components/button", data)
//	}
//
//	// Register and use components
//	facades.View().RegisterComponent(&ButtonComponent{})
//
//	// Use component in templates
//	componentHTML, err := facades.View().RenderComponent("button", map[string]interface{}{
//	    "text":  "Click Me",
//	    "type":  "submit",
//	    "class": "btn btn-success btn-lg",
//	})
//
// Template Engine Integration:
//
//	// Multiple template engines
//	// HTML/Go templates (default)
//	html, err := facades.View().Render("home.html", data)
//
//	// Pug templates
//	html, err = facades.View().WithEngine("pug").Render("home.pug", data)
//
//	// Mustache templates
//	html, err = facades.View().WithEngine("mustache").Render("home.mustache", data)
//
//	// Handlebars templates
//	html, err = facades.View().WithEngine("handlebars").Render("home.hbs", data)
//
//	// Engine-specific helpers
//	facades.View().AddHelper("formatDate", func(date time.Time, format string) string {
//	    return date.Format(format)
//	})
//
//	facades.View().AddHelper("currency", func(amount float64, currency string) string {
//	    return fmt.Sprintf("%s%.2f", currency, amount)
//	})
//
//	facades.View().AddHelper("truncate", func(text string, length int) string {
//	    if len(text) <= length {
//	        return text
//	    }
//	    return text[:length] + "..."
//	})
//
//	// Custom filters and functions
//	facades.View().AddFilter("upper", strings.ToUpper)
//	facades.View().AddFilter("lower", strings.ToLower)
//	facades.View().AddFilter("title", strings.Title)
//
//	// Use in templates:
//	// {{name | upper}} or {{upper name}}
//	// {{description | truncate 100}}
//	// {{created_at | formatDate "2006-01-02"}}
//
// Caching and Performance:
//
//	// Template caching
//	facades.View().EnableCache(true)
//	facades.View().SetCacheTimeout(time.Hour)
//
//	// Pre-compile templates
//	err := facades.View().CompileAll("./templates")
//	if err != nil {
//	    log.Fatalf("Template compilation failed: %v", err)
//	}
//
//	// Template minification
//	facades.View().EnableMinification(true)
//
//	// Gzip compression
//	facades.View().EnableCompression(true)
//
//	// Fragment caching for expensive partials
//	type ExpensiveFragment struct {
//	    cacheKey string
//	    ttl      time.Duration
//	}
//
//	func (ef *ExpensiveFragment) Render(data map[string]interface{}) (string, error) {
//	    // Check cache first
//	    if cached := facades.Cache().Get(ef.cacheKey); cached != nil {
//	        return cached.(string), nil
//	    }
//
//	    // Render expensive template
//	    html, err := facades.View().Render("expensive-partial", data)
//	    if err != nil {
//	        return "", err
//	    }
//
//	    // Cache result
//	    facades.Cache().Put(ef.cacheKey, html, ef.ttl)
//
//	    return html, nil
//	}
//
// Email Template Rendering:
//
//	// Email templates with layout
//	type EmailRenderer struct {
//	    baseLayout string
//	}
//
//	func NewEmailRenderer() *EmailRenderer {
//	    return &EmailRenderer{
//	        baseLayout: "emails/layout",
//	    }
//	}
//
//	func (er *EmailRenderer) RenderWelcomeEmail(user User) (string, string, error) {
//	    data := map[string]interface{}{
//	        "user":        user,
//	        "app_name":    "My App",
//	        "support_url": "https://myapp.com/support",
//	        "year":        time.Now().Year(),
//	    }
//
//	    // Render HTML version
//	    htmlBody, err := facades.View().RenderWithLayout(
//	        "emails/welcome-html",
//	        er.baseLayout,
//	        data,
//	    )
//	    if err != nil {
//	        return "", "", err
//	    }
//
//	    // Render plain text version
//	    textBody, err := facades.View().Render("emails/welcome-text", data)
//	    if err != nil {
//	        return "", "", err
//	    }
//
//	    return htmlBody, textBody, nil
//	}
//
//	// Use email renderer
//	renderer := NewEmailRenderer()
//	htmlBody, textBody, err := renderer.RenderWelcomeEmail(user)
//	if err != nil {
//	    log.Printf("Email rendering failed: %v", err)
//	    return
//	}
//
//	// Send email
//	facades.Mail().Send(mail.Message{
//	    To:       []string{user.Email},
//	    Subject:  "Welcome to My App!",
//	    HTMLBody: htmlBody,
//	    TextBody: textBody,
//	})
//
// PDF Generation from Templates:
//
//	// Render template to PDF
//	func GenerateInvoicePDF(invoice Invoice) ([]byte, error) {
//	    data := map[string]interface{}{
//	        "invoice":     invoice,
//	        "company":     getCompanyInfo(),
//	        "generated":   time.Now(),
//	        "total_words": numberToWords(invoice.Total),
//	    }
//
//	    // Render HTML first
//	    html, err := facades.View().Render("invoices/pdf-template", data)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Convert HTML to PDF
//	    pdf, err := facades.View().RenderToPDF(html, map[string]interface{}{
//	        "format":      "A4",
//	        "orientation": "portrait",
//	        "margins": map[string]string{
//	            "top":    "1cm",
//	            "bottom": "1cm",
//	            "left":   "1cm",
//	            "right":  "1cm",
//	        },
//	    })
//
//	    return pdf, err
//	}
//
//	// Stream PDF to response
//	func HandleInvoicePDF(w http.ResponseWriter, r *http.Request) {
//	    invoiceID := getInvoiceID(r)
//	    invoice := getInvoice(invoiceID)
//
//	    pdf, err := GenerateInvoicePDF(invoice)
//	    if err != nil {
//	        http.Error(w, "PDF generation failed", http.StatusInternalServerError)
//	        return
//	    }
//
//	    w.Header().Set("Content-Type", "application/pdf")
//	    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%d.pdf", invoice.ID))
//	    w.Header().Set("Content-Length", strconv.Itoa(len(pdf)))
//
//	    w.Write(pdf)
//	}
//
// Internationalization:
//
//	// Multi-language template rendering
//	func HandleMultiLanguagePage(w http.ResponseWriter, r *http.Request) {
//	    locale := getLocale(r) // From URL, cookie, or header
//
//	    // Set locale for this request
//	    facades.View().SetLocale(locale)
//
//	    data := map[string]interface{}{
//	        "user":     getCurrentUser(r),
//	        "products": getProducts(),
//	        "locale":   locale,
//	    }
//
//	    // Render localized template
//	    template := fmt.Sprintf("pages/home.%s", locale)
//	    html, err := facades.View().Render(template, data)
//	    if err != nil {
//	        // Fallback to default locale
//	        html, err = facades.View().Render("pages/home.en", data)
//	        if err != nil {
//	            http.Error(w, "Template error", http.StatusInternalServerError)
//	            return
//	        }
//	    }
//
//	    w.Header().Set("Content-Type", "text/html")
//	    w.Header().Set("Content-Language", locale)
//	    w.Write([]byte(html))
//	}
//
//	// Translation helpers in templates
//	facades.View().AddHelper("t", func(key string, params ...interface{}) string {
//	    return facades.Lang().Get(key, params...)
//	})
//
//	// Use in templates:
//	// {{t "welcome.message" .user.name}}
//	// {{t "product.count" .count}}
//
// Error Handling and Debugging:
//
//	// Custom error templates
//	func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
//	    data := map[string]interface{}{
//	        "status_code": statusCode,
//	        "error":       err.Error(),
//	        "request_id":  getRequestID(r),
//	        "timestamp":   time.Now(),
//	    }
//
//	    // Try to render custom error page
//	    template := fmt.Sprintf("errors/%d", statusCode)
//	    html, renderErr := facades.View().Render(template, data)
//	    if renderErr != nil {
//	        // Fallback to generic error page
//	        html, renderErr = facades.View().Render("errors/generic", data)
//	        if renderErr != nil {
//	            // Final fallback to plain text
//	            http.Error(w, err.Error(), statusCode)
//	            return
//	        }
//	    }
//
//	    w.Header().Set("Content-Type", "text/html")
//	    w.WriteHeader(statusCode)
//	    w.Write([]byte(html))
//	}
//
//	// Template debugging
//	facades.View().EnableDebug(true)
//	facades.View().SetDebugLevel("verbose")
//
//	// Debug information in templates
//	facades.View().AddHelper("debug", func(value interface{}) string {
//	    if !facades.App().IsDebug() {
//	        return ""
//	    }
//
//	    jsonData, _ := json.MarshalIndent(value, "", "  ")
//	    return fmt.Sprintf("<pre>%s</pre>", string(jsonData))
//	})
//
// Testing Templates:
//
//	// Test template rendering
//	func TestTemplateRendering(t *testing.T) {
//	    testCases := []struct {
//	        name     string
//	        template string
//	        data     map[string]interface{}
//	        expected string
//	        hasError bool
//	    }{
//	        {
//	            name:     "basic template",
//	            template: "basic",
//	            data: map[string]interface{}{
//	                "title": "Test Title",
//	                "name":  "John",
//	            },
//	            expected: "<h1>Test Title</h1><p>Hello, John!</p>",
//	            hasError: false,
//	        },
//	        {
//	            name:     "missing template",
//	            template: "nonexistent",
//	            data:     map[string]interface{}{},
//	            expected: "",
//	            hasError: true,
//	        },
//	    }
//
//	    for _, tc := range testCases {
//	        t.Run(tc.name, func(t *testing.T) {
//	            html, err := facades.View().Render(tc.template, tc.data)
//
//	            if tc.hasError {
//	                assert.Error(t, err)
//	            } else {
//	                assert.NoError(t, err)
//	                assert.Equal(t, tc.expected, html)
//	            }
//	        })
//	    }
//	}
//
//	// Mock view service for testing
//	func TestWithMockView(t *testing.T) {
//	    mockView := &MockViewService{
//	        templates: map[string]string{
//	            "test": "<p>{{.message}}</p>",
//	        },
//	    }
//
//	    // Swap view service
//	    restore := support.SwapService("view", mockView)
//	    defer restore()
//
//	    // Test with mock
//	    html, err := facades.View().Render("test", map[string]interface{}{
//	        "message": "Hello, World!",
//	    })
//
//	    assert.NoError(t, err)
//	    assert.Equal(t, "<p>Hello, World!</p>", html)
//	}
//
// Best Practices:
//   - Separate templates by feature/module for better organization
//   - Use layouts and partials to avoid code duplication
//   - Cache compiled templates in production for performance
//   - Validate template data before rendering to avoid runtime errors
//   - Use helpers and filters for common formatting tasks
//   - Implement proper error handling with fallback templates
//   - Use components for reusable UI elements
//   - Test templates thoroughly with various data scenarios
//   - Optimize template performance with fragment caching
//   - Use proper escaping to prevent XSS vulnerabilities
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume view service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	view, err := facade.TryResolve[ViewInterface]("view")
//	if err != nil {
//	    // Handle view service unavailability gracefully
//	    log.Printf("View service unavailable: %v", err)
//	    return // Skip template rendering
//	}
//	html, err := view.Render("template", data)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestViewBehavior(t *testing.T) {
//	    // Create a test view service
//	    testView := &TestView{
//	        templates: make(map[string]string),
//	    }
//
//	    // Swap the real view with test view
//	    restore := support.SwapService("view", testView)
//	    defer restore() // Always restore after test
//
//	    // Now facades.View() returns testView
//	    html, err := facades.View().Render("test", data)
//
//	    // Verify view behavior
//	    assert.NoError(t, err)
//	    assert.NotEmpty(t, html)
//	}
//
// Container Configuration:
// Ensure the view service is properly configured in your container:
//
//	// Example view registration
//	container.Singleton("view", func() interface{} {
//	    config := view.Config{
//	        // Template directories
//	        TemplateDirs: []string{
//	            "./templates",
//	            "./resources/views",
//	        },
//
//	        // Default template engine
//	        DefaultEngine: "html",
//
//	        // Supported engines
//	        Engines: map[string]view.EngineConfig{
//	            "html": {
//	                Engine:     "go-template",
//	                Extensions: []string{".html", ".tmpl"},
//	                Delimiters: []string{"{{", "}}"},
//	            },
//	            "pug": {
//	                Engine:     "pug",
//	                Extensions: []string{".pug"},
//	            },
//	            "mustache": {
//	                Engine:     "mustache",
//	                Extensions: []string{".mustache"},
//	            },
//	        },
//
//	        // Caching configuration
//	        Cache: view.CacheConfig{
//	            Enabled:    !facades.App().IsDebug(),
//	            TTL:        time.Hour,
//	            Driver:     "memory", // or "redis"
//	            KeyPrefix:  "view_cache:",
//	        },
//
//	        // Performance settings
//	        Performance: view.PerformanceConfig{
//	            PrecompileTemplates: true,
//	            EnableMinification:  true,
//	            EnableCompression:   true,
//	            ConcurrentRendering: true,
//	        },
//
//	        // Security settings
//	        Security: view.SecurityConfig{
//	            AutoEscape:     true,
//	            CSPNonce:       true,
//	            AllowedHosts:   []string{"localhost", "myapp.com"},
//	        },
//
//	        // Helper functions
//	        Helpers: map[string]interface{}{
//	            "formatDate": func(date time.Time) string {
//	                return date.Format("January 2, 2006")
//	            },
//	            "currency": func(amount float64) string {
//	                return fmt.Sprintf("$%.2f", amount)
//	            },
//	            "asset": func(path string) string {
//	                return fmt.Sprintf("/assets/%s", path)
//	            },
//	        },
//
//	        // Layout configuration
//	        Layouts: view.LayoutConfig{
//	            DefaultLayout: "layouts/app",
//	            LayoutDir:     "layouts",
//	            PartialDir:    "partials",
//	            ComponentDir:  "components",
//	        },
//	    }
//
//	    viewService, err := view.NewViewService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create view service: %v", err)
//	    }
//
//	    return viewService
//	})
func View() viewInterfaces.ViewInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "view" service from the dependency injection container
	// - Performs type assertion to ViewInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[viewInterfaces.ViewInterface](viewInterfaces.VIEW_TOKEN)
}

// ViewWithError provides error-safe access to the view service.
//
// This function offers the same functionality as View() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle view service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as View() but with error handling.
//
// Returns:
//   - ViewInterface: The resolved view instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement ViewInterface
//
// Usage Examples:
//
//	// Basic error-safe template rendering
//	view, err := facades.ViewWithError()
//	if err != nil {
//	    log.Printf("View service unavailable: %v", err)
//	    return // Skip template rendering
//	}
//	html, err := view.Render("template", data)
//
//	// Conditional template rendering
//	if view, err := facades.ViewWithError(); err == nil {
//	    // Render optional template
//	    html, err := view.Render("optional-template", data)
//	    if err == nil {
//	        // Use rendered HTML
//	        w.Write([]byte(html))
//	    }
//	}
func ViewWithError() (viewInterfaces.ViewInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "view" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[viewInterfaces.ViewInterface](viewInterfaces.VIEW_TOKEN)
}
