package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"govel/packages/ignition/config"
	"govel/packages/ignition/models"
	"govel/packages/ignition/renderer"
)

func main() {
	// Create HTTP server
	mux := http.NewServeMux()

	// Home page
	mux.HandleFunc("/", homePage)

	// Static assets for GoVel modules
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("/Users/akouta/Projects/govel/ignition/views/assets/"))))

	// Error endpoints
	mux.HandleFunc("/panic", handlePanic)
	mux.HandleFunc("/error", handleError)
	mux.HandleFunc("/nil-pointer", handleNilPointer)
	mux.HandleFunc("/solution", handleSolution)

	// Ignition-compatible endpoints
	mux.HandleFunc("/_ignition/health-check", handleHealthCheck)
	mux.HandleFunc("/_ignition/execute-solution", handleExecuteSolution)

	// Action endpoints for runnable solutions
	mux.HandleFunc("/generate-app-key", handleGenerateAppKey)
	mux.HandleFunc("/create-config", handleCreateConfig)
	mux.HandleFunc("/generate-recovery", handleGenerateRecovery)

	port := 3000
	fmt.Printf("🔥 GoVel Demo Server Starting!\n")
	fmt.Printf("🌐 Open your browser and visit: http://localhost:%d\n", port)
	fmt.Printf("⚡ The server will catch and display beautiful error pages with GoVel customizations\n\n")

	log.Printf("Server running on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoVel Demo</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            background: #f8fafc;
            color: #2d3748;
        }
        .header {
            text-align: center;
            margin-bottom: 3rem;
        }
        .demo-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 1.5rem;
            margin-top: 2rem;
        }
        .demo-card {
            background: white;
            padding: 1.5rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            border: 1px solid #e2e8f0;
            text-decoration: none;
            color: inherit;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .demo-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        .demo-title {
            font-size: 1.2rem;
            font-weight: 600;
            margin-bottom: 0.5rem;
            color: #e53e3e;
        }
        .demo-desc {
            color: #718096;
            font-size: 0.9rem;
        }
    </style>
</head>
<body>
    <div class="header">
        <div style="font-size: 3rem; margin-bottom: 1rem;">🐹</div>
        <h1>GoVel Error Pages</h1>
        <p>Beautiful error pages for Go applications with AMD module customizations</p>
    </div>

    <div class="demo-grid">
        <a href="/panic" class="demo-card">
            <div class="demo-title">Panic Error</div>
            <div class="demo-desc">Demonstrates panic recovery with GoVel customizations</div>
        </a>
        
        <a href="/error" class="demo-card">
            <div class="demo-title">Slice Error</div>
            <div class="demo-desc">Index out of range error with stack trace</div>
        </a>
        
        <a href="/nil-pointer" class="demo-card">
            <div class="demo-title">Nil Pointer</div>
            <div class="demo-desc">Classic nil pointer dereference error</div>
        </a>
        
        <a href="/solution" class="demo-card">
            <div class="demo-title">💡 Solution Demo</div>
            <div class="demo-desc">Error with suggested solutions and links</div>
        </a>

        <a href="/assets/js/govel/test.html" class="demo-card">
            <div class="demo-title">🧪 Module Tests</div>
            <div class="demo-desc">Test the GoVel AMD modules directly</div>
        </a>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// Error handler functions
func handlePanic(w http.ResponseWriter, r *http.Request) {
	defer recoverWithGoVel(w, r)
	panic("This is a deliberate panic to demonstrate GoVel error pages! 🐹")
}

func handleError(w http.ResponseWriter, r *http.Request) {
	defer recoverWithGoVel(w, r)
	// Cause an index out of range panic
	var slice []string
	_ = slice[10]
}

func handleNilPointer(w http.ResponseWriter, r *http.Request) {
	defer recoverWithGoVel(w, r)
	var ptr *string
	_ = *ptr // Nil pointer dereference
}

func handleSolution(w http.ResponseWriter, r *http.Request) {
	defer recoverWithGoVelSolutions(w, r)
	panic("Example error with solutions! Check out the suggested fixes below 🔧")
}

// Recovery function that uses GoVel renderer
func recoverWithGoVel(w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		log.Printf("Panic recovered: %v", rec)

		// Create error report
		report := models.NewErrorReport()
		report.SetType("runtime.panic")
		report.SetMessage(fmt.Sprintf("%v", rec))

		// Enhanced stack trace capture
		captureEnhancedStackTrace(report)

		// Create renderer and render error page
		htmlRenderer := renderer.NewHTMLRenderer()
		cfg := config.NewConfig()

		log.Printf("Stack frames captured: %d", len(report.GetStack()))
		for i, frame := range report.GetStack() {
			log.Printf("Frame %d: %s in %s:%d", i, frame.GetFunction(), frame.GetFile(), frame.GetLine())
		}

		htmlRenderer.RenderErrorPage(report, w, r, cfg, "/Users/akouta/Projects/govel", "", "")
	}
}

// captureEnhancedStackTrace captures stack trace with runtime.Caller and fallback to debug.Stack
func captureEnhancedStackTrace(report *models.ErrorReport) {
	// Method 1: Try runtime.Caller for detailed frame info
	captured := false
	for i := 1; i < 50; i++ { // Start from 1 to skip this function, increase depth
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		funcName := "unknown"
		if fn != nil {
			funcName = fn.Name()
		}

		// Skip runtime internal functions but keep some for context
		if !strings.Contains(funcName, "runtime.go") || i < 10 {
			frame := models.NewStackFrame()
			frame.SetFunction(funcName)
			frame.SetFile(file)
			frame.SetLine(line)

			// Try to read code snippet around the line
			addCodeSnippet(frame, file, line)

			report.AddStackFrame(*frame)
			captured = true
		}
	}

	// Method 2: Fallback to debug.Stack() if runtime.Caller didn't work well
	if !captured || len(report.GetStack()) < 3 {
		log.Printf("Using debug.Stack() as fallback")
		parseDebugStack(report)
	}
}

// addCodeSnippet tries to read source code around the error line
func addCodeSnippet(frame *models.StackFrame, file string, line int) {
	if file == "" || line <= 0 {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		return // File not accessible, skip code snippet
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	currentLine := 1
	startLine := line - 5
	endLine := line + 5
	if startLine < 1 {
		startLine = 1
	}

	for scanner.Scan() {
		if currentLine >= startLine && currentLine <= endLine {
			frame.AddCodeLine(fmt.Sprintf("%d", currentLine), scanner.Text())
		}
		if currentLine > endLine {
			break
		}
		currentLine++
	}
}

// parseDebugStack parses the output of debug.Stack()
func parseDebugStack(report *models.ErrorReport) {
	stackBytes := debug.Stack()
	stackLines := strings.Split(string(stackBytes), "\n")

	for i := 1; i < len(stackLines); i += 2 {
		if i+1 >= len(stackLines) {
			break
		}

		funcLine := strings.TrimSpace(stackLines[i])
		fileLine := strings.TrimSpace(stackLines[i+1])

		if funcLine == "" || fileLine == "" {
			continue
		}

		// Parse file and line number
		parts := strings.Split(fileLine, ":")
		if len(parts) < 2 {
			continue
		}

		file := parts[0]
		lineNumStr := parts[1]

		// Remove +0x... suffix if present
		if spaceIdx := strings.Index(lineNumStr, " "); spaceIdx >= 0 {
			lineNumStr = lineNumStr[:spaceIdx]
		}

		lineNum, err := strconv.Atoi(lineNumStr)
		if err != nil {
			continue
		}

		frame := models.NewStackFrame()
		frame.SetFunction(funcLine)
		frame.SetFile(file)
		frame.SetLine(lineNum)

		// Try to add code snippet
		addCodeSnippet(frame, file, lineNum)

		report.AddStackFrame(*frame)
	}
}

// Recovery function with solutions demonstration
func recoverWithGoVelSolutions(w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		log.Printf("Panic recovered (with solutions): %v", rec)

		// Create error report
		report := models.NewErrorReport()
		report.SetType("runtime.panic")
		report.SetMessage(fmt.Sprintf("%v", rec))

		// Enhanced stack trace capture
		captureEnhancedStackTrace(report)

		// Add practical runnable solutions
		solution1 := models.NewRunnableSolution(
			"Generate App Key 🔑",
			"Your application is missing a secure app key. Generate a cryptographically secure key for encryption and sessions.",
			"This will generate a new secure application key and show you how to set it up",
			"Generate Key",
			"http://localhost:3000/_ignition/execute-solution")
		report.AddSolution(*solution1)

		solution2 := models.NewRunnableSolution(
			"Create Configuration File 📋",
			"Create a proper configuration file for your application with sensible defaults and environment variable support.",
			"This will create a config.go file with proper structure",
			"Create Config",
			"http://localhost:3000/_ignition/execute-solution")
		report.AddSolution(*solution2)

		solution3 := models.NewGoDocumentationSolution(
			"Learn Go Error Patterns 📚",
			"Understanding Go's error handling patterns will help you write more robust applications.")
		solution3.AddLink("Go by Example: Errors", "https://gobyexample.com/errors")
		solution3.AddLink("Effective Go: Errors", "https://golang.org/doc/effective_go#errors")
		report.AddSolution(*solution3)

		solution4 := models.NewRunnableSolution(
			"Generate Recovery Middleware 🛡️",
			"Generate a panic recovery middleware that gracefully handles errors and provides consistent error responses.",
			"This will create a recovery middleware for your HTTP handlers",
			"Generate Middleware",
			"http://localhost:3000/_ignition/execute-solution")
		report.AddSolution(*solution4)

		// Create renderer and render error page with solutions
		htmlRenderer := renderer.NewHTMLRenderer()
		cfg := config.NewConfig()

		log.Printf("Stack frames captured (with solutions): %d", len(report.GetStack()))
		for i, frame := range report.GetStack() {
			log.Printf("Frame %d: %s in %s:%d", i, frame.GetFunction(), frame.GetFile(), frame.GetLine())
		}

		htmlRenderer.RenderErrorPage(report, w, r, cfg, "/Users/akouta/Projects/govel", "", "")
	}
}

// Ignition-compatible handlers
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"can_execute_commands":true}`))
}

func handleExecuteSolution(w http.ResponseWriter, r *http.Request) {
	// This endpoint routes solution execution requests to the appropriate handler
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body to get the solution type or action
	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		// If no JSON body, check for query parameters or route based on endpoint pattern
		// For now, we'll assume this is a Laravel-compatible interface
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Route to the appropriate handler based on the solution
	// For demo purposes, default to app key generation
	// In a real implementation, you'd route based on the solution data
	handleGenerateAppKey(w, r)
}

// Action handlers for runnable solutions
func handleGenerateAppKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Generate a secure 32-byte key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Printf("Error generating app key: %v", err)
		http.Error(w, `{"success": false, "message": "Failed to generate key"}`, http.StatusInternalServerError)
		return
	}

	// Encode the key in base64
	appKey := base64.StdEncoding.EncodeToString(key)

	// Create environment variable format
	envLine := fmt.Sprintf("APP_KEY=base64:%s", appKey)

	// Try to write to .env file
	envFilePath := ".env"
	envContent := ""
	if data, err := os.ReadFile(envFilePath); err == nil {
		envContent = string(data)
	}

	// Update or add APP_KEY
	lines := strings.Split(envContent, "\n")
	appKeyFound := false
	for i, line := range lines {
		if strings.HasPrefix(line, "APP_KEY=") {
			lines[i] = envLine
			appKeyFound = true
			break
		}
	}

	if !appKeyFound {
		lines = append(lines, envLine)
	}

	newEnvContent := strings.Join(lines, "\n")
	envWritten := false

	if err := os.WriteFile(envFilePath, []byte(newEnvContent), 0644); err == nil {
		envWritten = true
		log.Printf("App key generated and saved to .env file")
	} else {
		log.Printf("Could not write to .env file: %v", err)
	}

	// Return JSON response
	response := map[string]interface{}{
		"success":     true,
		"message":     "Application key generated successfully!",
		"app_key":     appKey,
		"env_line":    envLine,
		"env_written": envWritten,
		"instructions": []string{
			"Your new application key has been generated.",
			"Add this to your .env file: " + envLine,
			"Use this key for encrypting sessions and other sensitive data.",
			"Keep this key secure and never commit it to version control.",
		},
	}

	if envWritten {
		response["instructions"] = []string{
			"✅ Application key generated and saved to .env file!",
			"The key is now ready to use for encryption and sessions.",
			"Restart your application to use the new key.",
			"Keep your .env file secure and never commit it to version control.",
		}
	}

	json.NewEncoder(w).Encode(response)
}

func handleCreateConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	configContent := `package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	AppKey      string
	AppName     string
	AppEnv      string
	AppDebug    bool
	Port        int
	DatabaseURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		AppKey:      getEnv("APP_KEY", ""),
		AppName:     getEnv("APP_NAME", "GoVel App"),
		AppEnv:      getEnv("APP_ENV", "development"),
		AppDebug:    getEnvBool("APP_DEBUG", true),
		Port:        getEnvInt("PORT", 3000),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

// IsProduction returns true if the app is running in production
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// IsDevelopment returns true if the app is running in development
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}
`

	// Try to create the config directory and file
	if err := os.MkdirAll("config", 0755); err != nil {
		log.Printf("Error creating config directory: %v", err)
		http.Error(w, `{"success": false, "message": "Failed to create config directory"}`, http.StatusInternalServerError)
		return
	}

	configFile := "config/config.go"
	fileCreated := false

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err == nil {
		fileCreated = true
		log.Printf("Config file created: %s", configFile)
	} else {
		log.Printf("Could not write config file: %v", err)
	}

	// Create example .env file
	envExample := `# Application Configuration
APP_NAME="GoVel App"
APP_ENV=development
APP_DEBUG=true
APP_KEY=

# Server Configuration
PORT=3000

# Database Configuration
DATABASE_URL=postgres://user:password@localhost/dbname?sslmode=disable
`

	envExampleCreated := false
	if err := os.WriteFile(".env.example", []byte(envExample), 0644); err == nil {
		envExampleCreated = true
		log.Printf("Example .env file created")
	}

	response := map[string]interface{}{
		"success":       true,
		"message":       "Configuration files created successfully!",
		"files_created": []string{},
		"instructions":  []string{},
	}

	if fileCreated {
		response["files_created"] = append(response["files_created"].([]string), "config/config.go")
	}
	if envExampleCreated {
		response["files_created"] = append(response["files_created"].([]string), ".env.example")
	}

	response["instructions"] = []string{
		"✅ Configuration structure created!",
		"📁 Check the config/config.go file for your configuration structure",
		"📄 Copy .env.example to .env and configure your environment variables",
		"🔧 Import and use: cfg := config.LoadConfig()",
		"🔑 Don't forget to generate an app key if you haven't already!",
	}

	json.NewEncoder(w).Encode(response)
}

func handleGenerateRecovery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	recoveryContent := `package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// ErrorResponse represents a JSON error response
type ErrorResponse struct {
	Error   string ` + "`json:\"error\"`" + `
	Message string ` + "`json:\"message\"`" + `
	Status  int    ` + "`json:\"status\"`" + `
}

// RecoveryMiddleware provides panic recovery for HTTP handlers
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Panic recovered: %v\n%s", rec, debug.Stack())
				
				// Set appropriate headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				
				// Create error response
				errorResponse := ErrorResponse{
					Error:   "internal_server_error",
					Message: "An internal server error occurred",
					Status:  http.StatusInternalServerError,
				}
				
				// In development, include panic details
				if isDevelopment() {
					errorResponse.Message = fmt.Sprintf("Panic: %v", rec)
				}
				
				json.NewEncoder(w).Encode(errorResponse)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// isDevelopment checks if the app is in development mode
// You should implement this based on your configuration
func isDevelopment() bool {
	// Replace with your actual development check
	return true
}

// Usage example:
// mux := http.NewServeMux()
// mux.HandleFunc("/", yourHandler)
// 
// // Wrap with recovery middleware
// handler := RecoveryMiddleware(mux)
// http.ListenAndServe(":8080", handler)
`

	// Try to create the middleware directory and file
	if err := os.MkdirAll("middleware", 0755); err != nil {
		log.Printf("Error creating middleware directory: %v", err)
		http.Error(w, `{"success": false, "message": "Failed to create middleware directory"}`, http.StatusInternalServerError)
		return
	}

	middlewareFile := "middleware/recovery.go"
	fileCreated := false

	if err := os.WriteFile(middlewareFile, []byte(recoveryContent), 0644); err == nil {
		fileCreated = true
		log.Printf("Recovery middleware created: %s", middlewareFile)
	} else {
		log.Printf("Could not write middleware file: %v", err)
	}

	response := map[string]interface{}{
		"success":      true,
		"message":      "Recovery middleware generated successfully!",
		"file_created": middlewareFile,
		"instructions": []string{
			"✅ Recovery middleware generated!",
			"📁 Check middleware/recovery.go for the implementation",
			"🔧 Wrap your HTTP handler: RecoveryMiddleware(yourHandler)",
			"🛡️ All panics will now be caught and return JSON error responses",
			"📋 See usage example in the generated file",
		},
	}

	if !fileCreated {
		response["success"] = false
		response["message"] = "Failed to create middleware file"
	}

	json.NewEncoder(w).Encode(response)
}
