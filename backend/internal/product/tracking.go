package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// TrackingRequest is the request body for tracking
type TrackingRequest struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	TargetPrice float64 `json:"target_price" binding:"required"`
}

// UpdateTracking create or update the tracking of a product for a user
func UpdateTracking(c *gin.Context) {
	var req TrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	usuarioID := userIDContext.(uint)

	var tracking models.Tracking
	err := database.DB.Where("user_id = ? AND product_id = ?", usuarioID, req.ProductID).First(&tracking).Error

	if err != nil {
		tracking = models.Tracking{
			UserID:             usuarioID,
			ProductID:          req.ProductID,
			NotifyPriceChanges: true,
			NotifyTargetPrice:  req.TargetPrice,
			TrackingStartDate:  time.Now(),
		}
		if err := database.DB.Create(&tracking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la alerta de precio"})
			return
		}
	} else {
		tracking.NotifyTargetPrice = req.TargetPrice
		tracking.NotifyPriceChanges = true

		if err := database.DB.Save(&tracking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la alerta de precio"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alerta de precio actualizada con éxito"})
}
