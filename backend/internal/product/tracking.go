//nolint:gosec
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// TrackingRequest is the request body for tracking
type TrackingRequest struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	TargetPrice float64 `json:"target_price"`
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
		_ = updateUserSavedProducts(usuarioID, req.ProductID, true)
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

// CheckTracking checks if a user is following a specific product
func CheckTracking(c *gin.Context) {
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	usuarioID := userIDContext.(uint)
	productID := c.Param("id")

	var tracking models.Tracking
	err := database.DB.Where("user_id = ? AND product_id = ?", usuarioID, productID).First(&tracking).Error

	isFollowing := err == nil

	c.JSON(http.StatusOK, gin.H{
		"is_following": isFollowing,
		"target_price": tracking.NotifyTargetPrice,
	})
}

// UnfollowProduct deletes a product from the user's tracking list
func UnfollowProduct(c *gin.Context) {
	userIDContext, _ := c.Get("userID")
	usuarioID := userIDContext.(uint)
	productIDStr := c.Param("id")

	productID64, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	productID := uint(productID64)

	if err := database.DB.Where("user_id = ? AND product_id = ?", usuarioID, productID).Delete(&models.Tracking{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al dejar de seguir"})
		return
	}

	_ = updateUserSavedProducts(usuarioID, productID, false)

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado de tus seguimientos"})
}

func updateUserSavedProducts(userID uint, productID uint, add bool) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if add {
		// Añadir si no existe
		for _, id := range user.SavedProducts {
			if id == int64(productID) {
				return nil
			}
		}
		user.SavedProducts = append(user.SavedProducts, int64(productID))
	} else {
		// Eliminar
		newArray := []int64{}
		for _, id := range user.SavedProducts {
			if id != int64(productID) {
				newArray = append(newArray, id)
			}
		}
		user.SavedProducts = newArray
	}
	return database.DB.Save(&user).Error
}
