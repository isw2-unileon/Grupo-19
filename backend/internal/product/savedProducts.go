package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// GetSavedProducts devuelve la lista detallada de productos guardados por el usuario
func GetSavedProducts(c *gin.Context) {
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userID := userIDContext.(uint)
	var user models.User

	// Find User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Return of an empty array if there is no products saved
	if len(user.SavedProducts) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Product{}})
		return
	}

	var savedProducts []models.Product

	// Find the saved products
	if err := database.DB.Where("product_id IN ?", []int64(user.SavedProducts)).Find(&savedProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar los productos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": savedProducts})
}
