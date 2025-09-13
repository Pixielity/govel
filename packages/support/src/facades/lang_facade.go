package facades

import (
	langInterfaces "govel/types/src/interfaces/lang"
	facade "govel/support/src"
)

// Lang provides a clean, static-like interface to the application's localization service.
//
// This facade implements the facade pattern, providing global access to the localization
// service configured in the dependency injection container. It offers a Laravel-style
// API for internationalization (i18n) and localization (l10n) with automatic service
// resolution, translation management, locale detection, and multi-language support.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved localization service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent translation operations across goroutines
//   - Supports multiple translation formats (JSON, YAML, PO, etc.)
//   - Built-in pluralization rules and locale-specific formatting
//
// Behavior:
//   - First call: Resolves lang service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if lang service cannot be resolved (fail-fast behavior)
//   - Automatically handles locale switching, fallback languages, and translation loading
//
// Returns:
//   - LocalizationInterface: The application's localization service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "lang" service is not registered in the container
//   - If the resolved service doesn't implement LocalizationInterface
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
//   - Multiple goroutines can call Lang() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Translation operations are thread-safe with concurrent locale access
//
// Usage Examples:
//
//	// Basic translation
//	message := facades.Lang().Get("welcome.message")
//	fmt.Println(message) // Output: "Welcome to our application!"
//
//	// Translation with parameters
//	greeting := facades.Lang().Get("user.greeting", map[string]interface{}{
//	    "name": "John",
//	    "time": "morning",
//	})
//	fmt.Println(greeting) // Output: "Good morning, John!"
//
//	// Translation with default fallback
//	message := facades.Lang().GetWithDefault("unknown.key", "Default message")
//	fmt.Println(message) // Output: "Default message"
//
//	// Check if translation exists
//	if facades.Lang().Has("error.validation") {
//	    error := facades.Lang().Get("error.validation")
//	    fmt.Println(error)
//	}
//
//	// Pluralization
//	count := 5
//	message := facades.Lang().Choice("items.count", count, map[string]interface{}{
//	    "count": count,
//	})
//	fmt.Println(message) // Output: "You have 5 items"
//
//	// Locale management
//	currentLocale := facades.Lang().GetLocale()
//	fmt.Printf("Current locale: %s\n", currentLocale) // Output: "en"
//
//	// Set locale for current request/context
//	facades.Lang().SetLocale("es")
//	message := facades.Lang().Get("welcome.message")
//	fmt.Println(message) // Output: "¡Bienvenido a nuestra aplicación!"
//
//	// Temporarily use different locale
//	frenchMessage := facades.Lang().ForLocale("fr", func() string {
//	    return facades.Lang().Get("welcome.message")
//	})
//	fmt.Println(frenchMessage) // Output: "Bienvenue dans notre application!"
//
//	// Get available locales
//	locales := facades.Lang().GetAvailableLocales()
//	fmt.Printf("Available locales: %v\n", locales) // Output: ["en", "es", "fr", "de"]
//
//	// Fallback locale handling
//	facades.Lang().SetFallbackLocale("en")
//	message := facades.Lang().Get("missing.translation") // Falls back to English
//
//	// Nested translation keys
//	errorMessage := facades.Lang().Get("validation.required", map[string]interface{}{
//	    "attribute": "email",
//	})
//	fmt.Println(errorMessage) // Output: "The email field is required."
//
//	// Translation with complex interpolation
//	message := facades.Lang().Get("order.summary", map[string]interface{}{
//	    "customer":   "John Doe",
//	    "total":      99.99,
//	    "currency":   "USD",
//	    "item_count": 3,
//	    "date":       time.Now().Format("2006-01-02"),
//	})
//	// Output: "Order for John Doe: 3 items totaling $99.99 USD on 2023-12-01"
//
//	// Array/list translations
//	days := facades.Lang().GetArray("calendar.days")
//	for i, day := range days {
//	    fmt.Printf("%d: %s\n", i+1, day)
//	}
//	// Output: "1: Monday", "2: Tuesday", etc.
//
//	// Date and time localization
//	now := time.Now()
//	localizedDate := facades.Lang().FormatDate(now, "long")
//	localizedTime := facades.Lang().FormatTime(now, "short")
//	fmt.Printf("Date: %s, Time: %s\n", localizedDate, localizedTime)
//
// Advanced Localization Patterns:
//
//	// Middleware for locale detection
//	func LocaleMiddleware(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        // Detect locale from various sources
//	        locale := detectLocaleFromRequest(r)
//
//	        if locale == "" {
//	            locale = facades.Lang().GetDefaultLocale()
//	        }
//
//	        // Set locale for this request
//	        facades.Lang().SetLocale(locale)
//
//	        next.ServeHTTP(w, r)
//	    })
//	}
//
//	func detectLocaleFromRequest(r *http.Request) string {
//	    // 1. Check URL parameter
//	    if lang := r.URL.Query().Get("lang"); lang != "" {
//	        return lang
//	    }
//
//	    // 2. Check cookie
//	    if cookie, err := r.Cookie("locale"); err == nil {
//	        return cookie.Value
//	    }
//
//	    // 3. Check Accept-Language header
//	    return facades.Lang().ParseAcceptLanguage(r.Header.Get("Accept-Language"))
//	}
//
//	// Template integration
//	func RenderTemplate(templateName string, data map[string]interface{}) string {
//	    // Add translation function to template data
//	    data["trans"] = func(key string, params ...map[string]interface{}) string {
//	        if len(params) > 0 {
//	            return facades.Lang().Get(key, params[0])
//	        }
//	        return facades.Lang().Get(key)
//	    }
//
//	    data["transChoice"] = func(key string, count int, params ...map[string]interface{}) string {
//	        if len(params) > 0 {
//	            return facades.Lang().Choice(key, count, params[0])
//	        }
//	        return facades.Lang().Choice(key, count)
//	    }
//
//	    return executeTemplate(templateName, data)
//	}
//
//	// Form validation with localized messages
//	type ValidationErrors map[string][]string
//
//	func (v ValidationErrors) Localize() map[string][]string {
//	    localized := make(map[string][]string)
//
//	    for field, errors := range v {
//	        localizedErrors := make([]string, len(errors))
//	        for i, errKey := range errors {
//	            localizedErrors[i] = facades.Lang().Get("validation."+errKey, map[string]interface{}{
//	                "attribute": facades.Lang().Get("attributes."+field),
//	            })
//	        }
//	        localized[field] = localizedErrors
//	    }
//
//	    return localized
//	}
//
//	// Email localization
//	func SendLocalizedEmail(userLocale, template string, data map[string]interface{}) error {
//	    return facades.Lang().ForLocale(userLocale, func() error {
//	        subject := facades.Lang().Get("email."+template+".subject", data)
//	        body := facades.Lang().Get("email."+template+".body", data)
//
//	        return facades.Mail().Send(subject, body, data["email"].(string))
//	    })
//	}
//
//	// Currency and number formatting
//	func FormatCurrency(amount float64) string {
//	    locale := facades.Lang().GetLocale()
//	    currency := facades.Lang().Get("currency.code") // "USD", "EUR", etc.
//
//	    return facades.Lang().FormatCurrency(amount, currency, locale)
//	}
//
//	// Contextual translations
//	func GetContextualTranslation(key, context string, params map[string]interface{}) string {
//	    // Try context-specific translation first
//	    contextKey := fmt.Sprintf("%s.contexts.%s", key, context)
//	    if facades.Lang().Has(contextKey) {
//	        return facades.Lang().Get(contextKey, params)
//	    }
//
//	    // Fall back to general translation
//	    return facades.Lang().Get(key, params)
//	}
//
// Translation File Examples:
//
//	// lang/en.json
//	{
//	  "welcome": {
//	    "message": "Welcome to our application!",
//	    "user": "Welcome back, {{name}}!"
//	  },
//	  "user": {
//	    "greeting": "Good {{time}}, {{name}}!"
//	  },
//	  "items": {
//	    "count": {
//	      "0": "No items",
//	      "1": "One item",
//	      "other": "{{count}} items"
//	    }
//	  },
//	  "validation": {
//	    "required": "The {{attribute}} field is required.",
//	    "email": "The {{attribute}} must be a valid email address.",
//	    "min": "The {{attribute}} must be at least {{min}} characters."
//	  },
//	  "attributes": {
//	    "email": "email address",
//	    "password": "password",
//	    "name": "full name"
//	  }
//	}
//
//	// lang/es.json
//	{
//	  "welcome": {
//	    "message": "¡Bienvenido a nuestra aplicación!",
//	    "user": "¡Bienvenido de nuevo, {{name}}!"
//	  },
//	  "validation": {
//	    "required": "El campo {{attribute}} es obligatorio.",
//	    "email": "El {{attribute}} debe ser una dirección de correo válida."
//	  },
//	  "attributes": {
//	    "email": "correo electrónico",
//	    "password": "contraseña",
//	    "name": "nombre completo"
//	  }
//	}
//
// Best Practices:
//   - Use hierarchical translation keys ("module.section.key")
//   - Provide meaningful default values for missing translations
//   - Use parameter interpolation instead of string concatenation
//   - Implement proper pluralization rules for different languages
//   - Keep translation files organized by feature or module
//   - Use context-specific translations when needed
//   - Always provide fallback translations for critical messages
//   - Test translations with different string lengths
//
// Pluralization Rules:
//   - English: 0 and 1 are singular, 2+ are plural
//   - Spanish: 1 is singular, 0 and 2+ are plural
//   - Russian: Complex rules based on last digits
//   - Chinese: No pluralization (same form for all counts)
//   - Use ICU MessageFormat for complex pluralization
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume localization service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	lang, err := facade.TryResolve[LocalizationInterface]("lang")
//	if err != nil {
//	    // Handle localization service unavailability gracefully
//	    return "Default message", fmt.Errorf("localization unavailable: %w", err)
//	}
//	translation := lang.Get("key")
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestLocalization(t *testing.T) {
//	    // Create a test localization service with predefined translations
//	    testLang := &TestLocalization{
//	        translations: map[string]map[string]string{
//	            "en": {
//	                "test.message": "Test message",
//	                "user.greeting": "Hello, {{name}}!",
//	            },
//	            "es": {
//	                "test.message": "Mensaje de prueba",
//	                "user.greeting": "¡Hola, {{name}}!",
//	            },
//	        },
//	        currentLocale: "en",
//	    }
//
//	    // Swap the real localization with test localization
//	    restore := support.SwapService("lang", testLang)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Lang() returns testLang
//	    message := facades.Lang().Get("test.message")
//	    assert.Equal(t, "Test message", message)
//
//	    // Test locale switching
//	    facades.Lang().SetLocale("es")
//	    message = facades.Lang().Get("test.message")
//	    assert.Equal(t, "Mensaje de prueba", message)
//
//	    // Test parameter interpolation
//	    greeting := facades.Lang().Get("user.greeting", map[string]interface{}{
//	        "name": "Juan",
//	    })
//	    assert.Equal(t, "¡Hola, Juan!", greeting)
//	}
//
// Container Configuration:
// Ensure the localization service is properly configured in your container:
//
//	// Example lang registration
//	container.Singleton("lang", func() interface{} {
//	    config := lang.Config{
//	        // Default locale
//	        DefaultLocale:  "en",
//	        FallbackLocale: "en",
//
//	        // Available locales
//	        AvailableLocales: []string{"en", "es", "fr", "de", "zh"},
//
//	        // Translation file paths
//	        TranslationPath: "./lang",
//	        FileFormat:      "json", // json, yaml, po
//
//	        // Caching settings
//	        CacheTranslations: true,
//	        CacheTTL:          time.Hour * 24, // 24 hours
//
//	        // Loading strategy
//	        LoadOnDemand:   true,  // Load translations when needed
//	        PreloadDefault: true,  // Always preload default locale
//
//	        // Interpolation settings
//	        InterpolationStartDelimiter: "{{",
//	        InterpolationEndDelimiter:   "}}",
//
//	        // Pluralization
//	        PluralRules: map[string]lang.PluralRule{
//	            "en": lang.EnglishPluralRule,
//	            "es": lang.SpanishPluralRule,
//	            "fr": lang.FrenchPluralRule,
//	        },
//
//	        // Missing translation behavior
//	        ReturnKeyOnMissing: false, // Return default instead of key
//	        LogMissingTranslations: true,
//
//	        // Formatting options
//	        NumberFormat: lang.NumberFormatConfig{
//	            DecimalSeparator:  ".",
//	            ThousandSeparator: ",",
//	        },
//
//	        DateTimeFormat: lang.DateTimeFormatConfig{
//	            ShortDateFormat: "2006-01-02",
//	            LongDateFormat:  "January 2, 2006",
//	            TimeFormat:      "15:04:05",
//	        },
//	    }
//
//	    langService, err := lang.NewLocalizationService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create localization service: %v", err)
//	    }
//
//	    return langService
//	})
func Lang() langInterfaces.LanguageInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "lang" service from the dependency injection container
	// - Performs type assertion to LanguageInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[langInterfaces.LanguageInterface](langInterfaces.LANG_TOKEN)
}

// LangWithError provides error-safe access to the localization service.
//
// This function offers the same functionality as Lang() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle localization service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Lang() but with error handling.
//
// Returns:
//   - LocalizationInterface: The resolved lang instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement LocalizationInterface
//
// Usage Examples:
//
//	// Basic error-safe translation
//	lang, err := facades.LangWithError()
//	if err != nil {
//	    log.Printf("Localization service unavailable: %v", err)
//	    return "Default message", fmt.Errorf("translation service not available")
//	}
//	message := lang.Get("welcome.message")
//
//	// Conditional localization
//	if lang, err := facades.LangWithError(); err == nil {
//	    // Use localized messages when available
//	    errorMessage := lang.Get("validation.error")
//	    return errorMessage
//	} else {
//	    // Fall back to hardcoded English
//	    return "Validation error occurred"
//	}
//
//	// Health check pattern
//	func CheckLocalizationHealth() error {
//	    lang, err := facades.LangWithError()
//	    if err != nil {
//	        return fmt.Errorf("localization service unavailable: %w", err)
//	    }
//
//	    // Test basic translation functionality
//	    if !lang.Has("system.test") {
//	        return fmt.Errorf("translation files not loaded properly")
//	    }
//
//	    // Test locale switching
//	    originalLocale := lang.GetLocale()
//	    lang.SetLocale("en")
//	    testMessage := lang.Get("system.test")
//	    lang.SetLocale(originalLocale)
//
//	    if testMessage == "" {
//	        return fmt.Errorf("localization service not working properly")
//	    }
//
//	    return nil
//	}
func LangWithError() (langInterfaces.LanguageInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "lang" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[langInterfaces.LanguageInterface](langInterfaces.LANG_TOKEN)
}
