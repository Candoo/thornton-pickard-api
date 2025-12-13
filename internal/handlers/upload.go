package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Candoo/thornton-pickard-api/internal/services"
)

var storageService = services.NewStorageService()

// UploadImage handles image uploads
// @Summary Upload an image
// @Description Upload an image file (JPEG, PNG, GIF)
// @Tags uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} map[string]string
// @Router /upload [post]
func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
			return
		}

		// Validate file type
		ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
		allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
		
		valid := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: JPG, PNG, GIF, WebP"})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Max size: 5MB"})
			return
		}

		// Save file
		url, err := storageService.SaveFile(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"url":      url,
			"filename": file.Filename,
			"size":     file.Size,
		})
	}
}

// UploadMultipleImages handles multiple image uploads
// @Summary Upload multiple images
// @Description Upload multiple image files at once
// @Tags uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Image files"
// @Success 200 {object} map[string][]string
// @Router /upload/multiple [post]
func UploadMultipleImages() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
			return
		}

		if len(files) > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 files allowed"})
			return
		}

		urls := make([]string, 0, len(files))
		
		for _, file := range files {
			// Validate file type
			ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
			allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
			
			valid := false
			for _, allowedExt := range allowedExts {
				if ext == allowedExt {
					valid = true
					break
				}
			}

			if !valid {
				continue // Skip invalid files
			}

			// Validate file size
			if file.Size > 5*1024*1024 {
				continue // Skip large files
			}

			// Save file
			url, err := storageService.SaveFile(file)
			if err == nil {
				urls = append(urls, url)
			}
		}

		if len(urls) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No valid files uploaded"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"urls":  urls,
			"count": len(urls),
		})
	}
}