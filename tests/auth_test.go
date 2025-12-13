package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/Candoo/thornton-pickard-api/internal/handlers"
	"github.com/Candoo/thornton-pickard-api/internal/models"
)

func TestRegister(t *testing.T) {
	db := setupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/register", handlers.Register(db))

	// Create request
	reqBody := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "test@example.com", response.User.Email)
}

func TestRegisterDuplicateEmail(t *testing.T) {
	db := setupTestDB()

	// Create existing user
	user := models.User{Email: "existing@example.com"}
	user.HashPassword("password123")
	db.Create(&user)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/register", handlers.Register(db))

	// Try to register with same email
	reqBody := models.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestLogin(t *testing.T) {
	db := setupTestDB()

	// Create user first
	user := models.User{Email: "test@example.com"}
	user.HashPassword("password123")
	db.Create(&user)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/login", handlers.Login(db))

	// Login request
	reqBody := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response.Token)
}

func TestLoginInvalidPassword(t *testing.T) {
	db := setupTestDB()

	// Create user
	user := models.User{Email: "test@example.com"}
	user.HashPassword("password123")
	db.Create(&user)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/login", handlers.Login(db))

	// Login with wrong password
	reqBody := models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonData, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}