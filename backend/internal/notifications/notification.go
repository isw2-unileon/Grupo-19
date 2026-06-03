package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// Returns the alerts in the database for the authenticated user
func GetUserNotifications(c *gin.Context) {
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	usuarioID := userIDContext.(uint)

	var list []models.Notification
	err := database.DB.Where("user_id = ?", usuarioID).Order("created_at desc").Find(&list).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar las notificaciones"})
		return
	}

	c.JSON(http.StatusOK, list)
}
