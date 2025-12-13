package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/Candoo/thornton-pickard-api/internal/handlers"
	"github.com/Candoo/thornton-pickard-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
	_ "modernc.org/sqlite"
)

func setupTestDB() *gorm.DB {
	// Try different SQLite connection strings
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		panic("failed to connect to test database: " + err.Error())
	}
	
	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Camera{},
		&models.User{},
		&models.Manufacturer{},
		&models.Ephemera{},
	)
	if err != nil {
		fmt.Printf("MIGRATION ERROR: %v\n", err)
		panic("failed to migrate test database: " + err.Error())
	}
	
	return db
}

func TestGetCameras(t *testing.T) {
	db := setupTestDB()
	
	// Create test camera
	camera := models.Camera{
		Name:           "Test Camera",
		Manufacturer:   "Test Manufacturer",
		YearIntroduced: 1900,
	}
	db.Create(&camera)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/cameras", handlers.GetCameras(db))

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cameras", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["data"])
}

func TestGetCamera(t *testing.T) {
	db := setupTestDB()
	
	// Create test camera
	camera := models.Camera{
		Name:           "Test Camera",
		Manufacturer:   "Test Manufacturer",
		YearIntroduced: 1900,
	}
	db.Create(&camera)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/cameras/:id", handlers.GetCamera(db))

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cameras/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Camera
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Camera", response.Name)
}

func TestCreateCamera(t *testing.T) {
	db := setupTestDB()

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/cameras", handlers.CreateCamera(db))

	// Create request body
	camera := models.Camera{
		Name:           "New Camera",
		Manufacturer:   "New Manufacturer",
		YearIntroduced: 1910,
	}
	jsonData, _ := json.Marshal(camera)

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cameras", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.Camera
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "New Camera", response.Name)
}

func TestSearchCameras(t *testing.T) {
	db := setupTestDB()
	
	// Create test cameras
	cameras := []models.Camera{
		{Name: "Ruby Reflex", Manufacturer: "Thornton-Pickard", YearIntroduced: 1909},
		{Name: "Imperial", Manufacturer: "Thornton-Pickard", YearIntroduced: 1895},
	}
	db.Create(&cameras)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/cameras", handlers.GetCameras(db))

	// Test search by name
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cameras?search=ruby", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(data))
}

func TestPagination(t *testing.T) {
	db := setupTestDB()
	
	// Create 15 test cameras
	for i := 1; i <= 15; i++ {
		camera := models.Camera{
			Name:           "Camera " + string(rune(48+i)), // ASCII numbers
			Manufacturer:   "Test",
			YearIntroduced: 1900 + i,
		}
		db.Create(&camera)
	}

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/cameras", handlers.GetCameras(db))

	// Test first page
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cameras?page=1&page_size=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(10), response["page_size"])
	assert.Equal(t, float64(15), response["total"])
	assert.Equal(t, float64(2), response["total_pages"])
}