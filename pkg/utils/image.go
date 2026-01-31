package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ImageUploadDir adalah direktori untuk menyimpan uploaded images
const ImageUploadDir = "uploads/images"

// SaveBase64Image menyimpan base64 image ke file system dan return path-nya
func SaveBase64Image(base64String string) (string, error) {
	// Remove data URI scheme if present (e.g., "data:image/png;base64,...")
	if strings.Contains(base64String, ",") {
		parts := strings.Split(base64String, ",")
		if len(parts) == 2 {
			base64String = parts[1]
		}
	}

	// Decode base64 string
	imageData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Detect image type from first bytes
	ext := detectImageType(imageData)
	if ext == "" {
		return "", fmt.Errorf("unsupported image type")
	}

	// Create uploads directory if not exists
	if err := os.MkdirAll(ImageUploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(ImageUploadDir, filename)

	// Write to file
	if err := os.WriteFile(filepath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to write image file: %w", err)
	}

	// Return relative path
	return filepath, nil
}

// DeleteImage menghapus image file dari file system
func DeleteImage(imagePath string) error {
	if imagePath == "" {
		return nil
	}
	
	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	return os.Remove(imagePath)
}

// detectImageType mendeteksi tipe image dari magic bytes
func detectImageType(data []byte) string {
	if len(data) < 4 {
		return ""
	}

	// PNG: 89 50 4E 47
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return ".png"
	}

	// JPEG: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg"
	}

	// GIF: 47 49 46
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		return ".gif"
	}

	// WebP: 52 49 46 46 (RIFF) ... 57 45 42 50 (WEBP)
	if len(data) >= 12 && data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 {
		if data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
			return ".webp"
		}
	}

	return ""
}

// IsBase64Image mengecek apakah string adalah base64 image
func IsBase64Image(str string) bool {
	// Check if it starts with data URI scheme
	if strings.HasPrefix(str, "data:image/") {
		return true
	}

	// Check if it's a valid base64 string (simple check)
	if len(str) > 100 && !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
		// Try to decode to verify
		_, err := base64.StdEncoding.DecodeString(str)
		return err == nil
	}

	return false
}
