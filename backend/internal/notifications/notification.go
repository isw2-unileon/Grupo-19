package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// GetUserNotifications returns the alerts in the database for the authenticated user
func GetUserNotifications(c *gin.Context) {
	userIDContext, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userID := userIDContext.(uint)

	var list []models.Notification
	err := database.DB.Where("user_id = ?", userID).
		Order("is_read ASC").
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener notificaciones"})
		return
	}

	c.JSON(http.StatusOK, list)
}

// MarkNotificationAsRead marks a notifications as read
func MarkNotificationAsRead(c *gin.Context) {
	// Extract the notification ID from the URL
	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificación inválido"})
		return
	}

	// Extract the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Update the "is_read" field to true in the database.
	result := database.DB.Model(&models.Notification{}).
		Where("notification_id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la notificación"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notificación no encontrada o no autorizada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notificación marcada como leída con éxito"})
}

// DeleteNotification deletes a user notification
func DeleteNotification(c *gin.Context) {
	// Extract the notification ID from the URL
	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificación inválido"})
		return
	}

	// 2. Extract the userID from the session context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Delete the record
	result := database.DB.Where("notification_id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la notificación"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notificación no encontrada o no autorizada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notificación eliminada con éxito"})
}
