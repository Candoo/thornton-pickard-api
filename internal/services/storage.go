package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type StorageService struct {
	uploadDir string
}

func NewStorageService() *StorageService {
	uploadDir := "./uploads"
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		uploadDir = dir
	}

	// Create upload directory if it doesn't exist
	os.MkdirAll(uploadDir, 0755)

	return &StorageService{
		uploadDir: uploadDir,
	}
}

func (s *StorageService) SaveFile(file *multipart.FileHeader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(s.uploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return relative URL path
	return "/uploads/" + filename, nil
}

func (s *StorageService) DeleteFile(url string) error {
	// Extract filename from URL
	filename := filepath.Base(url)
	filepath := filepath.Join(s.uploadDir, filename)
	return os.Remove(filepath)
}