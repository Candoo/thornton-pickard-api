package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/Candoo/thornton-pickard-api/internal/models"
)

func GetEphemera(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ephemera []models.Ephemera
		
		if err := db.Find(&ephemera).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ephemera)
	}
}

func GetEphemeraItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var item models.Ephemera

		if err := db.First(&item, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func CreateEphemeraItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item models.Ephemera

		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

func GetManufacturers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var manufacturers []models.Manufacturer
		
		if err := db.Find(&manufacturers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, manufacturers)
	}
}