package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// IsDirectory checks if a path is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile checks if a path is a regular file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// EnsureDirectory creates a directory if it doesn't exist
func EnsureDirectory(path string) error {
	if !IsDirectory(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// FindFiles recursively finds files matching a pattern
func FindFiles(rootPath, pattern string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			matched, matchErr := filepath.Match(pattern, filepath.Base(path))
			if matchErr != nil {
				return matchErr
			}
			if matched {
				files = append(files, path)
			}
		}
		
		return nil
	})
	
	return files, err
}

// ContainsGovelModule checks if a go.mod file contains govel modules
func ContainsGovelModule(content string) bool {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "govel/") || strings.Contains(line, "@govel/") {
			return true
		}
	}
	return false
}

// GetRelativePath returns the relative path from base to target
func GetRelativePath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

// JoinPath joins path elements safely
func JoinPath(elements ...string) string {
	return filepath.Join(elements...)
}

// GetFileSize returns the size of a file
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Ensure destination directory exists
	if err := EnsureDirectory(filepath.Dir(dst)); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// RemoveFile removes a file if it exists
func RemoveFile(path string) error {
	if FileExists(path) {
		return os.Remove(path)
	}
	return nil
}

// ReadFileLines reads a file and returns its lines
func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// WriteLines writes lines to a file
func WriteLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}