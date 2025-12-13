package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/Candoo/thornton-pickard-api/internal/models"
	"github.com/Candoo/thornton-pickard-api/internal/utils"
)

// GetCameras retrieves all cameras with pagination
// @Summary List all cameras
// @Description Get a paginated list of all cameras
// @Tags cameras
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Param manufacturer query string false "Filter by manufacturer"
// @Param year_from query int false "Filter by year from"
// @Param year_to query int false "Filter by year to"
// @Param format query string false "Filter by format"
// @Param sort query string false "Sort by field (name, year_introduced)" default(name)
// @Param order query string false "Sort order (asc, desc)" default(asc)
// @Success 200 {object} utils.Pagination
// @Router /cameras [get]
func GetCameras(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cameras []models.Camera
		var total int64

		// Build query with filters
		query := db.Model(&models.Camera{})

		// Search
		if search := c.Query("search"); search != "" {
			query = query.Where(
				"name ILIKE ? OR manufacturer ILIKE ? OR description ILIKE ?",
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

		// Count total before pagination
		query.Count(&total)

		// Sorting
		sortField := c.DefaultQuery("sort", "name")
		sortOrder := c.DefaultQuery("order", "asc")
		if sortOrder != "asc" && sortOrder != "desc" {
			sortOrder = "asc"
		}
		query = query.Order(sortField + " " + sortOrder)

		// Apply pagination
		if err := query.Scopes(utils.Paginate(c)).Find(&cameras).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := utils.CreatePaginationResponse(c, cameras, total)
		c.JSON(http.StatusOK, response)
	}
}

// GetCamera retrieves a single camera by ID
// @Summary Get a camera
// @Description Get a camera by ID
// @Tags cameras
// @Produce json
// @Param id path int true "Camera ID"
// @Success 200 {object} models.Camera
// @Router /cameras/{id} [get]
func GetCamera(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var camera models.Camera

		if err := db.First(&camera, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Camera not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, camera)
	}
}

// CreateCamera creates a new camera
// @Summary Create a camera
// @Description Create a new camera entry
// @Tags cameras
// @Accept json
// @Produce json
// @Param camera body models.Camera true "Camera object"
// @Success 201 {object} models.Camera
// @Router /cameras [post]
func CreateCamera(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var camera models.Camera

		if err := c.ShouldBindJSON(&camera); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&camera).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, camera)
	}
}

// UpdateCamera updates an existing camera
// @Summary Update a camera
// @Description Update a camera by ID
// @Tags cameras
// @Accept json
// @Produce json
// @Param id path int true "Camera ID"
// @Param camera body models.Camera true "Camera object"
// @Success 200 {object} models.Camera
// @Router /cameras/{id} [put]
func UpdateCamera(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var camera models.Camera

		if err := db.First(&camera, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Camera not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := c.ShouldBindJSON(&camera); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&camera).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, camera)
	}
}

// DeleteCamera deletes a camera
// @Summary Delete a camera
// @Description Delete a camera by ID
// @Tags cameras
// @Param id path int true "Camera ID"
// @Success 204
// @Router /cameras/{id} [delete]
func DeleteCamera(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		if err := db.Delete(&models.Camera{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// SearchCameras searches for cameras (deprecated, use GetCameras with query params)
// @Summary Search cameras
// @Description Search cameras by name, manufacturer, or year
// @Tags cameras
// @Produce json
// @Param q query string false "Search query"
// @Param year query int false "Year introduced"
// @Success 200 {array} models.Camera
// @Router /cameras/search [get]
func SearchCameras(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		yearStr := c.Query("year")

		var cameras []models.Camera
		dbQuery := db.Model(&models.Camera{})

		if query != "" {
			dbQuery = dbQuery.Where(
				"name ILIKE ? OR manufacturer ILIKE ?",
				"%"+query+"%",
				"%"+query+"%",
			)
		}

		if yearStr != "" {
			if year, err := strconv.Atoi(yearStr); err == nil {
				dbQuery = dbQuery.Where("year_introduced = ?", year)
			}
		}

		if err := dbQuery.Find(&cameras).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, cameras)
	}
}