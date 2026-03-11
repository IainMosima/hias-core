package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var dangerousExtensions = map[string]bool{
	".exe": true, ".sh": true, ".bat": true, ".cmd": true,
	".php": true, ".jsp": true, ".com": true, ".msi": true,
	".ps1": true, ".vbs": true, ".js": true, ".wsf": true,
}

func ValidateFileName(name string) error {
	if name == "" {
		return fmt.Errorf("file name is required")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("file name must not contain path separators")
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("file name must not contain path traversal sequences")
	}
	if strings.ContainsRune(name, 0) {
		return fmt.Errorf("file name must not contain null bytes")
	}
	ext := strings.ToLower(filepath.Ext(name))
	if dangerousExtensions[ext] {
		return fmt.Errorf("file extension %s is not allowed", ext)
	}
	return nil
}

func ValidateFileSize(size int64, maxSize int) error {
	if size <= 0 {
		return fmt.Errorf("file size must be greater than 0")
	}
	if maxSize > 0 && size > int64(maxSize) {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", size, maxSize)
	}
	return nil
}

func ValidateMimeType(mimeType string, allowedTypes []string) error {
	if mimeType == "" {
		return fmt.Errorf("mime type is required")
	}
	if len(allowedTypes) == 0 {
		return nil
	}
	for _, allowed := range allowedTypes {
		if strings.EqualFold(mimeType, allowed) {
			return nil
		}
	}
	return fmt.Errorf("mime type %s is not allowed", mimeType)
}

func GenerateS3Key(entityType, entityID, docType, fileName string) string {
	sanitized := sanitizeFileName(fileName)
	timestamp := time.Now().Unix()
	unique := uuid.New().String()[:8]
	return fmt.Sprintf("%s/%s/%s/%d-%s-%s", entityType, entityID, docType, timestamp, unique, sanitized)
}

func sanitizeFileName(name string) string {
	name = filepath.Base(name)
	name = strings.Map(func(r rune) rune {
		if r == ' ' || r == '(' || r == ')' {
			return '_'
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			return r
		}
		return -1
	}, name)
	if len(name) > 100 {
		ext := filepath.Ext(name)
		name = name[:100-len(ext)] + ext
	}
	return name
}
