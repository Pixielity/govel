package facades

import (
	storageInterfaces "govel/packages/types/src/interfaces/storage"
	facade "govel/packages/support/src"
)

// Storage provides a clean, static-like interface to the application's file storage service.
//
// This facade implements the facade pattern, providing global access to the storage
// service configured in the dependency injection container. It offers a Laravel-style
// API for file storage operations with automatic service resolution, multiple storage drivers,
// cloud integration, and comprehensive file management capabilities.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved storage service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent file operations across goroutines
//   - Supports multiple storage drivers (local, S3, Google Cloud, Azure, FTP, SFTP)
//   - Built-in file streaming, metadata management, and URL generation
//
// Behavior:
//   - First call: Resolves storage service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if storage service cannot be resolved (fail-fast behavior)
//   - Automatically handles file operations, metadata, and cloud integration
//
// Returns:
//   - StorageInterface: The application's storage service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "storage" service is not registered in the container
//   - If the resolved service doesn't implement StorageInterface
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
//   - Multiple goroutines can call Storage() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - File operations are thread-safe with proper synchronization
//
// Usage Examples:
//
//	// Basic file operations
//	// Store a file from bytes
//	data := []byte("Hello, world!")
//	err := facades.Storage().Put("documents/readme.txt", data)
//	if err != nil {
//	    log.Printf("Failed to store file: %v", err)
//	}
//
//	// Read file contents
//	content, err := facades.Storage().Get("documents/readme.txt")
//	if err != nil {
//	    log.Printf("Failed to read file: %v", err)
//	} else {
//	    fmt.Printf("File content: %s\n", string(content))
//	}
//
//	// Check if file exists
//	if facades.Storage().Exists("documents/readme.txt") {
//	    fmt.Println("File exists")
//	}
//
//	// Delete a file
//	err = facades.Storage().Delete("documents/readme.txt")
//	if err != nil {
//	    log.Printf("Failed to delete file: %v", err)
//	}
//
//	// File streaming operations
//	// Stream file to response writer
//	func ServeFile(w http.ResponseWriter, r *http.Request) {
//	    filePath := r.URL.Path[1:] // Remove leading slash
//
//	    if !facades.Storage().Exists(filePath) {
//	        http.NotFound(w, r)
//	        return
//	    }
//
//	    // Get file metadata
//	    info, err := facades.Storage().GetFileInfo(filePath)
//	    if err != nil {
//	        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//	        return
//	    }
//
//	    // Set appropriate headers
//	    w.Header().Set("Content-Type", info.MimeType)
//	    w.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
//
//	    // Stream file content
//	    reader, err := facades.Storage().ReadStream(filePath)
//	    if err != nil {
//	        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//	        return
//	    }
//	    defer reader.Close()
//
//	    io.Copy(w, reader)
//	}
//
//	// File upload handling
//	func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
//	    // Parse multipart form
//	    err := r.ParseMultipartForm(10 << 20) // 10MB max
//	    if err != nil {
//	        http.Error(w, "Failed to parse form", http.StatusBadRequest)
//	        return
//	    }
//
//	    file, header, err := r.FormFile("upload")
//	    if err != nil {
//	        http.Error(w, "No file uploaded", http.StatusBadRequest)
//	        return
//	    }
//	    defer file.Close()
//
//	    // Generate unique filename
//	    filename := fmt.Sprintf("uploads/%s_%s",
//	        time.Now().Format("20060102_150405"),
//	        header.Filename)
//
//	    // Store uploaded file
//	    err = facades.Storage().PutStream(filename, file)
//	    if err != nil {
//	        http.Error(w, "Failed to store file", http.StatusInternalServerError)
//	        return
//	    }
//
//	    // Generate public URL
//	    url, err := facades.Storage().URL(filename)
//	    if err != nil {
//	        log.Printf("Failed to generate URL: %v", err)
//	        url = "/files/" + filename
//	    }
//
//	    w.Header().Set("Content-Type", "application/json")
//	    json.NewEncoder(w).Encode(map[string]string{
//	        "filename": filename,
//	        "url":      url,
//	        "status":   "uploaded",
//	    })
//	}
//
//	// Directory operations
//	// List files in directory
//	files, err := facades.Storage().Files("uploads")
//	if err != nil {
//	    log.Printf("Failed to list files: %v", err)
//	} else {
//	    for _, file := range files {
//	        fmt.Printf("File: %s\n", file)
//	    }
//	}
//
//	// List all files recursively
//	allFiles, err := facades.Storage().AllFiles("documents")
//	if err != nil {
//	    log.Printf("Failed to list all files: %v", err)
//	} else {
//	    for _, file := range allFiles {
//	        fmt.Printf("File: %s\n", file)
//	    }
//	}
//
//	// List directories
//	dirs, err := facades.Storage().Directories("uploads")
//	if err != nil {
//	    log.Printf("Failed to list directories: %v", err)
//	} else {
//	    for _, dir := range dirs {
//	        fmt.Printf("Directory: %s\n", dir)
//	    }
//	}
//
//	// Create directory
//	err = facades.Storage().MakeDirectory("uploads/images")
//	if err != nil {
//	    log.Printf("Failed to create directory: %v", err)
//	}
//
//	// Delete directory
//	err = facades.Storage().DeleteDirectory("uploads/temp")
//	if err != nil {
//	    log.Printf("Failed to delete directory: %v", err)
//	}
//
// Advanced Storage Patterns:
//
//	// Image processing and storage
//	type ImageProcessor struct {
//	    AllowedTypes []string
//	    MaxSize      int64
//	    Quality      int
//	}
//
//	func (p *ImageProcessor) ProcessAndStore(file io.Reader, filename string) (string, error) {
//	    // Read file into memory for processing
//	    data, err := io.ReadAll(file)
//	    if err != nil {
//	        return "", err
//	    }
//
//	    // Validate file type
//	    mimeType := http.DetectContentType(data)
//	    if !p.isAllowedType(mimeType) {
//	        return "", fmt.Errorf("unsupported file type: %s", mimeType)
//	    }
//
//	    // Check file size
//	    if int64(len(data)) > p.MaxSize {
//	        return "", fmt.Errorf("file too large: %d bytes", len(data))
//	    }
//
//	    // Process image (resize, compress, etc.)
//	    processedData, err := p.processImage(data)
//	    if err != nil {
//	        return "", err
//	    }
//
//	    // Generate storage path
//	    storagePath := fmt.Sprintf("images/%s/%s",
//	        time.Now().Format("2006/01/02"),
//	        filename)
//
//	    // Store processed image
//	    err = facades.Storage().Put(storagePath, processedData)
//	    if err != nil {
//	        return "", err
//	    }
//
//	    return storagePath, nil
//	}
//
//	func (p *ImageProcessor) isAllowedType(mimeType string) bool {
//	    for _, allowed := range p.AllowedTypes {
//	        if mimeType == allowed {
//	            return true
//	        }
//	    }
//	    return false
//	}
//
//	// File backup and versioning
//	type FileVersionManager struct {
//	    MaxVersions int
//	}
//
//	func (v *FileVersionManager) SaveVersion(filePath string, data []byte) error {
//	    // Create version directory if it doesn't exist
//	    versionDir := fmt.Sprintf("%s.versions", filePath)
//	    err := facades.Storage().MakeDirectory(versionDir)
//	    if err != nil && !os.IsExist(err) {
//	        return err
//	    }
//
//	    // Generate version filename with timestamp
//	    timestamp := time.Now().Format("20060102_150405")
//	    versionPath := fmt.Sprintf("%s/%s", versionDir, timestamp)
//
//	    // Store version
//	    err = facades.Storage().Put(versionPath, data)
//	    if err != nil {
//	        return err
//	    }
//
//	    // Clean up old versions
//	    return v.cleanupOldVersions(versionDir)
//	}
//
//	func (v *FileVersionManager) GetVersion(filePath, version string) ([]byte, error) {
//	    versionPath := fmt.Sprintf("%s.versions/%s", filePath, version)
//	    return facades.Storage().Get(versionPath)
//	}
//
//	func (v *FileVersionManager) ListVersions(filePath string) ([]string, error) {
//	    versionDir := fmt.Sprintf("%s.versions", filePath)
//	    return facades.Storage().Files(versionDir)
//	}
//
//	// Cloud storage synchronization
//	type CloudSync struct {
//	    LocalDisk  string
//	    CloudDisk  string
//	    SyncPrefix string
//	}
//
//	func (c *CloudSync) SyncToCloud(filePath string) error {
//	    // Read from local storage
//	    localStorage := facades.Storage().Disk(c.LocalDisk)
//	    data, err := localStorage.Get(filePath)
//	    if err != nil {
//	        return err
//	    }
//
//	    // Upload to cloud storage
//	    cloudStorage := facades.Storage().Disk(c.CloudDisk)
//	    cloudPath := fmt.Sprintf("%s/%s", c.SyncPrefix, filePath)
//
//	    return cloudStorage.Put(cloudPath, data)
//	}
//
//	func (c *CloudSync) SyncFromCloud(filePath string) error {
//	    // Download from cloud storage
//	    cloudStorage := facades.Storage().Disk(c.CloudDisk)
//	    cloudPath := fmt.Sprintf("%s/%s", c.SyncPrefix, filePath)
//
//	    data, err := cloudStorage.Get(cloudPath)
//	    if err != nil {
//	        return err
//	    }
//
//	    // Store in local storage
//	    localStorage := facades.Storage().Disk(c.LocalDisk)
//	    return localStorage.Put(filePath, data)
//	}
//
//	func (c *CloudSync) SyncDirectory(dirPath string) error {
//	    localStorage := facades.Storage().Disk(c.LocalDisk)
//	    files, err := localStorage.AllFiles(dirPath)
//	    if err != nil {
//	        return err
//	    }
//
//	    for _, file := range files {
//	        err := c.SyncToCloud(file)
//	        if err != nil {
//	            log.Printf("Failed to sync file %s: %v", file, err)
//	            continue
//	        }
//	    }
//
//	    return nil
//	}
//
// Storage Driver Management:
//
//	// Using different storage drivers
//	// Default driver
//	err := facades.Storage().Put("file.txt", []byte("content"))
//
//	// Specific driver
//	s3Storage := facades.Storage().Disk("s3")
//	err = s3Storage.Put("backup/file.txt", []byte("content"))
//
//	localStorage := facades.Storage().Disk("local")
//	err = localStorage.Put("temp/file.txt", []byte("content"))
//
//	// Cloud storage operations
//	gcsStorage := facades.Storage().Disk("gcs")
//	err = gcsStorage.Put("documents/file.pdf", pdfData)
//
//	// Generate signed URLs for private files
//	signedURL, err := s3Storage.TemporaryURL("private/document.pdf", time.Hour)
//	if err != nil {
//	    log.Printf("Failed to generate signed URL: %v", err)
//	} else {
//	    fmt.Printf("Signed URL: %s\n", signedURL)
//	}
//
//	// File metadata operations
//	info, err := facades.Storage().GetFileInfo("documents/report.pdf")
//	if err != nil {
//	    log.Printf("Failed to get file info: %v", err)
//	} else {
//	    fmt.Printf("File: %s\n", info.Name)
//	    fmt.Printf("Size: %d bytes\n", info.Size)
//	    fmt.Printf("Modified: %v\n", info.ModTime)
//	    fmt.Printf("MIME Type: %s\n", info.MimeType)
//	}
//
//	// Set custom metadata
//	metadata := map[string]string{
//	    "author":      "John Doe",
//	    "category":    "reports",
//	    "version":     "1.0",
//	    "description": "Monthly sales report",
//	}
//	err = facades.Storage().SetMetadata("documents/report.pdf", metadata)
//	if err != nil {
//	    log.Printf("Failed to set metadata: %v", err)
//	}
//
// File Operations with Progress:
//
//	// Large file upload with progress tracking
//	type ProgressTracker struct {
//	    Total    int64
//	    Current  int64
//	    Callback func(current, total int64)
//	}
//
//	func (p *ProgressTracker) Read(buf []byte) (int, error) {
//	    n, err := p.Reader.Read(buf)
//	    p.Current += int64(n)
//
//	    if p.Callback != nil {
//	        p.Callback(p.Current, p.Total)
//	    }
//
//	    return n, err
//	}
//
//	func UploadLargeFile(filePath string, reader io.Reader, size int64) error {
//	    // Wrap reader with progress tracking
//	    progress := &ProgressTracker{
//	        Reader: reader,
//	        Total:  size,
//	        Callback: func(current, total int64) {
//	            percentage := float64(current) / float64(total) * 100
//	            fmt.Printf("\rUpload progress: %.1f%% (%d/%d bytes)",
//	                percentage, current, total)
//	        },
//	    }
//
//	    return facades.Storage().PutStream(filePath, progress)
//	}
//
// Storage Security:
//
//	// File access control
//	type FileAccessControl struct {
//	    PublicPaths  []string
//	    PrivatePaths []string
//	    AdminPaths   []string
//	}
//
//	func (f *FileAccessControl) CanAccess(filePath string, userRole string) bool {
//	    // Check public files
//	    for _, publicPath := range f.PublicPaths {
//	        if strings.HasPrefix(filePath, publicPath) {
//	            return true
//	        }
//	    }
//
//	    // Check admin files
//	    if userRole == "admin" {
//	        for _, adminPath := range f.AdminPaths {
//	            if strings.HasPrefix(filePath, adminPath) {
//	                return true
//	            }
//	        }
//	    }
//
//	    // Check private files for authenticated users
//	    if userRole != "guest" {
//	        for _, privatePath := range f.PrivatePaths {
//	            if strings.HasPrefix(filePath, privatePath) {
//	                return true
//	            }
//	        }
//	    }
//
//	    return false
//	}
//
//	// File encryption/decryption
//	func StoreEncryptedFile(filePath string, data []byte, key []byte) error {
//	    // Encrypt data before storing
//	    encryptedData, err := facades.Crypt().Encrypt(data, key)
//	    if err != nil {
//	        return err
//	    }
//
//	    return facades.Storage().Put(filePath+".encrypted", encryptedData)
//	}
//
//	func GetEncryptedFile(filePath string, key []byte) ([]byte, error) {
//	    // Retrieve encrypted data
//	    encryptedData, err := facades.Storage().Get(filePath + ".encrypted")
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Decrypt and return
//	    return facades.Crypt().Decrypt(encryptedData, key)
//	}
//
// Testing Support:
//
//	// Test storage operations
//	func TestFileOperations(t *testing.T) {
//	    // Create test storage with in-memory driver
//	    testStorage := &TestStorage{
//	        files: make(map[string][]byte),
//	    }
//
//	    // Swap storage service
//	    restore := support.SwapService("storage", testStorage)
//	    defer restore()
//
//	    // Test file operations
//	    testData := []byte("test content")
//	    err := facades.Storage().Put("test/file.txt", testData)
//	    require.NoError(t, err)
//
//	    // Verify file exists
//	    assert.True(t, facades.Storage().Exists("test/file.txt"))
//
//	    // Read and verify content
//	    content, err := facades.Storage().Get("test/file.txt")
//	    require.NoError(t, err)
//	    assert.Equal(t, testData, content)
//
//	    // Test file deletion
//	    err = facades.Storage().Delete("test/file.txt")
//	    require.NoError(t, err)
//	    assert.False(t, facades.Storage().Exists("test/file.txt"))
//	}
//
// Best Practices:
//   - Use appropriate storage drivers for your use case (local for dev, cloud for production)
//   - Implement proper file validation and security checks
//   - Use streaming operations for large files to avoid memory issues
//   - Generate unique filenames to avoid conflicts
//   - Implement file versioning for important documents
//   - Use signed URLs for secure access to private files
//   - Consider file compression for storage optimization
//   - Implement proper error handling and logging
//   - Use metadata for file organization and searching
//   - Regular cleanup of temporary and old files
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume storage service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	storage, err := facade.TryResolve[StorageInterface]("storage")
//	if err != nil {
//	    // Handle storage service unavailability gracefully
//	    log.Printf("Storage service unavailable: %v", err)
//	    return // Skip storage operations
//	}
//	err = storage.Put("file.txt", []byte("content"))
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestStorageBehavior(t *testing.T) {
//	    // Create a test storage that records all operations
//	    testStorage := &TestStorage{
//	        files:      make(map[string][]byte),
//	        operations: []string{},
//	    }
//
//	    // Swap the real storage with test storage
//	    restore := support.SwapService("storage", testStorage)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Storage() returns testStorage
//	    err := facades.Storage().Put("test.txt", []byte("content"))
//	    require.NoError(t, err)
//
//	    // Verify storage behavior
//	    assert.Contains(t, testStorage.operations, "PUT:test.txt")
//	    assert.Equal(t, []byte("content"), testStorage.files["test.txt"])
//	}
//
// Container Configuration:
// Ensure the storage service is properly configured in your container:
//
//	// Example storage registration
//	container.Singleton("storage", func() interface{} {
//	    config := storage.Config{
//	        // Default storage driver
//	        DefaultDisk: "local",
//
//	        // Storage drivers configuration
//	        Disks: map[string]storage.DiskConfig{
//	            "local": {
//	                Driver: "local",
//	                Root:   "./storage/app",
//	                URL:    "/storage",
//	                Permissions: storage.Permissions{
//	                    File:      0644,
//	                    Directory: 0755,
//	                },
//	            },
//	            "public": {
//	                Driver: "local",
//	                Root:   "./storage/app/public",
//	                URL:    "/storage",
//	                Visibility: "public",
//	            },
//	            "s3": {
//	                Driver:   "s3",
//	                Key:      os.Getenv("AWS_ACCESS_KEY_ID"),
//	                Secret:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
//	                Region:   os.Getenv("AWS_DEFAULT_REGION"),
//	                Bucket:   os.Getenv("AWS_BUCKET"),
//	                URL:      os.Getenv("AWS_URL"),
//	                Endpoint: os.Getenv("AWS_ENDPOINT"),
//	            },
//	            "gcs": {
//	                Driver:            "gcs",
//	                ProjectID:         os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
//	                KeyFile:           os.Getenv("GOOGLE_CLOUD_KEY_FILE"),
//	                Bucket:            os.Getenv("GOOGLE_CLOUD_STORAGE_BUCKET"),
//	                PathPrefix:        os.Getenv("GOOGLE_CLOUD_STORAGE_PATH_PREFIX"),
//	                ApiURI:            os.Getenv("GOOGLE_CLOUD_STORAGE_API_URI"),
//	            },
//	            "azure": {
//	                Driver:      "azure",
//	                Name:        os.Getenv("AZURE_STORAGE_NAME"),
//	                Key:         os.Getenv("AZURE_STORAGE_KEY"),
//	                Container:   os.Getenv("AZURE_STORAGE_CONTAINER"),
//	                URL:         os.Getenv("AZURE_STORAGE_URL"),
//	            },
//	            "ftp": {
//	                Driver:   "ftp",
//	                Host:     os.Getenv("FTP_HOST"),
//	                Username: os.Getenv("FTP_USERNAME"),
//	                Password: os.Getenv("FTP_PASSWORD"),
//	                Port:     21,
//	                Root:     "/",
//	                Passive:  true,
//	                SSL:      false,
//	                Timeout:  30,
//	            },
//	        },
//	    }
//
//	    storageService, err := storage.NewStorageService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create storage service: %v", err)
//	    }
//
//	    return storageService
//	})
func Storage() storageInterfaces.StorageInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "storage" service from the dependency injection container
	// - Performs type assertion to StorageInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[storageInterfaces.StorageInterface](storageInterfaces.STORAGE_TOKEN)
}

// StorageWithError provides error-safe access to the storage service.
//
// This function offers the same functionality as Storage() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle storage service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Storage() but with error handling.
//
// Returns:
//   - StorageInterface: The resolved storage instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement StorageInterface
//
// Usage Examples:
//
//	// Basic error-safe file operations
//	storage, err := facades.StorageWithError()
//	if err != nil {
//	    log.Printf("Storage service unavailable: %v", err)
//	    return // Skip file operations
//	}
//	err = storage.Put("backup.txt", []byte("data"))
//
//	// Conditional file operations
//	if storage, err := facades.StorageWithError(); err == nil {
//	    // Store optional cache file
//	    storage.Put("cache/temp.json", jsonData)
//	}
func StorageWithError() (storageInterfaces.StorageInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "storage" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[storageInterfaces.StorageInterface](storageInterfaces.STORAGE_TOKEN)
}
