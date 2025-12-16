package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/Candoo/thornton-pickard-api/internal/models"
	"github.com/Candoo/thornton-pickard-api/internal/utils"
)

// Define a struct to hold the database dependency
type CameraHandler struct {
	DB *gorm.DB
}

// NewCameraHandler creates a new handler instance
func NewCameraHandler(db *gorm.DB) *CameraHandler {
	return &CameraHandler{DB: db}
}

func (h *CameraHandler) getCameraQuery(c *gin.Context) *gorm.DB {
	query := h.DB.Model(&models.Camera{})

	// Search
	if search := c.Query("search"); search != "" {
		query = query.Where(
			"name LIKE ? OR manufacturer LIKE ? OR description LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	// Filter by manufacturer
	if manufacturer := c.Query("manufacturer"); manufacturer != "" {
		query = query.Where("manufacturer = ?", manufacturer)
	}

	// Filter by year range
	if yearFrom := c.Query("year_from"); yearFrom != "" {
		query = query.Where("year_introduced >= ?", yearFrom)
	}
	if yearTo := c.Query("year_to"); yearTo != "" {
		query = query.Where("year_introduced <= ?", yearTo)
	}

	// Filter by format
	if format := c.Query("format"); format != "" {
		query = query.Where("format = ?", format)
	}

	// Sorting
	sortField := c.DefaultQuery("sort", "name")
	sortOrder := c.DefaultQuery("order", "asc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	validSortFields := map[string]bool{"name": true, "year_introduced": true, "rarity": true}
	if !validSortFields[sortField] {
		sortField = "name" // Default if invalid field is provided
	}

	query = query.Order(sortField + " " + sortOrder)
	
	return query
}

// GetCameras retrieves all cameras with pagination
// @Summary List all cameras
// @Description Get a paginated list of all cameras with filtering, sorting, and search
// @Tags cameras
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query (name, manufacturer, or description)"
// @Param manufacturer query string false "Filter by manufacturer"
// @Param year_from query int false "Filter by year from"
// @Param year_to query int false "Filter by year to"
// @Param format query string false "Filter by format"
// @Param sort query string false "Sort by field (name, year_introduced, rarity)" default(name)
// @Param order query string false "Sort order (asc, desc)" default(asc)
// @Success 200 {object} utils.Pagination
// @Router /cameras [get]
func (h *CameraHandler) GetCameras(c *gin.Context) {
	var cameras []models.Camera
	var total int64

	query := h.getCameraQuery(c)

	// Count total before pagination
	query.Count(&total)

	// Apply pagination and retrieve data
	if err := query.Scopes(utils.Paginate(c)).Find(&cameras).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cameras"})
		return
	}

	cameraResponses := make([]models.CameraResponse, len(cameras))
	for i, camera := range cameras {
		cameraResponses[i] = camera.ToCameraResponse()
	}

	response := utils.CreatePaginationResponse(c, cameraResponses, total)
	c.JSON(http.StatusOK, response)
}

// GetCamera retrieves a single camera by ID
// @Summary Get a camera
// @Description Get a camera by ID
// @Tags cameras
// @Produce json
// @Param id path int true "Camera ID"
// @Success 200 {object} models.CameraResponse
// @Router /cameras/{id} [get]
func (h *CameraHandler) GetCamera(c *gin.Context) {
	id := c.Param("id")
	var camera models.Camera

	if err := h.DB.First(&camera, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Camera not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	cameraResponse := camera.ToCameraResponse()
	c.JSON(http.StatusOK, cameraResponse)
}

// CreateCamera creates a new camera
// @Summary Create a camera
// @Description Create a new camera entry
// @Tags cameras
// @Accept json
// @Produce json
// @Param camera body models.Camera true "Camera object"
// @Success 201 {object} models.CameraResponse
// @Router /cameras [post]
func (h *CameraHandler) CreateCamera(c *gin.Context) {
	var camera models.Camera

	if err := c.ShouldBindJSON(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := h.DB.Create(&camera).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create camera: " + err.Error()})
		return
	}

	cameraResponse := camera.ToCameraResponse()
	c.JSON(http.StatusCreated, cameraResponse)
}

// UpdateCamera updates an existing camera
// @Summary Update a camera
// @Description Update a camera by ID
// @Tags cameras
// @Accept json
// @Produce json
// @Param id path int true "Camera ID"
// @Param camera body models.Camera true "Camera object"
// @Success 200 {object} models.CameraResponse
// @Router /cameras/{id} [put]
func (h *CameraHandler) UpdateCamera(c *gin.Context) {
	id := c.Param("id")
	var camera models.Camera

	if err := h.DB.First(&camera, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Camera not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := h.DB.Save(&camera).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update camera: " + err.Error()})
		return
	}

	cameraResponse := camera.ToCameraResponse()
	c.JSON(http.StatusOK, cameraResponse)
}

// DeleteCamera deletes a camera
// @Summary Delete a camera
// @Description Delete a camera by ID
// @Tags cameras
// @Param id path int true "Camera ID"
// @Success 204
// @Router /cameras/{id} [delete]
func (h *CameraHandler) DeleteCamera(c *gin.Context) {
	id := c.Param("id")
	
	// Perform soft delete using GORM
	if err := h.DB.Delete(&models.Camera{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete camera: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}